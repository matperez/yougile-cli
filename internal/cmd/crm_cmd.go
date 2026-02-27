package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// createContactPersonBody is used to build CreateContactPersonDto with optional Fields (anonymous in client).
type createContactPersonBody struct {
	Title     string `json:"title"`
	ProjectId string `json:"projectId"`
	Fields    *struct {
		Email           *string `json:"email,omitempty"`
		Phone           *string `json:"phone,omitempty"`
		Address         *string `json:"address,omitempty"`
		Position        *string `json:"position,omitempty"`
		AdditionalPhone *string `json:"additionalPhone,omitempty"`
	} `json:"fields,omitempty"`
}

// NewCrmContactPersonsCreateCmd returns the "crm contact-persons create" command.
func NewCrmContactPersonsCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, projectID, email, phone, address, position, additionalPhone string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a contact person in a CRM project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" || projectID == "" {
				return fmt.Errorf("title and project-id are required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			build := createContactPersonBody{Title: title, ProjectId: projectID}
			if email != "" || phone != "" || address != "" || position != "" || additionalPhone != "" {
				build.Fields = &struct {
					Email           *string `json:"email,omitempty"`
					Phone           *string `json:"phone,omitempty"`
					Address         *string `json:"address,omitempty"`
					Position        *string `json:"position,omitempty"`
					AdditionalPhone *string `json:"additionalPhone,omitempty"`
				}{}
				if email != "" {
					build.Fields.Email = &email
				}
				if phone != "" {
					build.Fields.Phone = &phone
				}
				if address != "" {
					build.Fields.Address = &address
				}
				if position != "" {
					build.Fields.Position = &position
				}
				if additionalPhone != "" {
					build.Fields.AdditionalPhone = &additionalPhone
				}
			}
			raw, _ := json.Marshal(build)
			var body client.CrmContactPersonsControllerCreateJSONRequestBody
			if err := json.Unmarshal(raw, &body); err != nil {
				return fmt.Errorf("build request: %w", err)
			}
			resp, err := api.CrmContactPersonsControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create contact person: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create contact person: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Contact person created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "contact name/title")
	c.Flags().StringVar(&projectID, "project-id", "", "CRM project ID")
	c.Flags().StringVar(&email, "email", "", "email")
	c.Flags().StringVar(&phone, "phone", "", "phone")
	c.Flags().StringVar(&address, "address", "", "address")
	c.Flags().StringVar(&position, "position", "", "position")
	c.Flags().StringVar(&additionalPhone, "additional-phone", "", "additional phone")
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewCrmContactsByExternalIdCmd returns the "crm contacts by-external-id" command.
func NewCrmContactsByExternalIdCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var provider, chatID string
	c := &cobra.Command{
		Use:   "by-external-id",
		Short: "Find contact by external integration (provider and chat ID)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if provider == "" || chatID == "" {
				return fmt.Errorf("provider and chat-id are required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.CrmExternalIdControllerFindContactByExternalIdParams{
				Provider: provider,
				ChatId:   chatID,
			}
			resp, err := api.CrmExternalIdControllerFindContactByExternalIdWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("find contact by external id: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("find contact by external id: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			if resp.JSON200 == nil {
				_, err = fmt.Fprintln(out, "No contact found")
				return err
			}
			return output.PrintJSON(out, resp.JSON200)
		},
	}
	c.Flags().StringVar(&provider, "provider", "", "external integration provider")
	c.Flags().StringVar(&chatID, "chat-id", "", "chat ID in the external messenger")
	_ = c.MarkFlagRequired("provider")
	_ = c.MarkFlagRequired("chat-id")
	return c
}

// NewCrmCmd returns the "crm" parent command.
func NewCrmCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "crm",
		Short: "CRM contacts and contact persons",
	}
	contactPersons := &cobra.Command{Use: "contact-persons", Short: "Contact persons"}
	contactPersons.AddCommand(NewCrmContactPersonsCreateCmd(resolvePath, outputJSON))
	c.AddCommand(contactPersons)
	contacts := &cobra.Command{Use: "contacts", Short: "Contacts"}
	contacts.AddCommand(NewCrmContactsByExternalIdCmd(resolvePath, outputJSON))
	c.AddCommand(contacts)
	return c
}
