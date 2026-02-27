package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewFilesUploadCmd returns the "files upload" command.
func NewFilesUploadCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "upload [file]",
		Short: "Upload a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			path := args[0]
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer func() { _ = f.Close() }()

			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			name := filepath.Base(path)
			part, err := w.CreateFormFile("file", name)
			if err != nil {
				return fmt.Errorf("create form file: %w", err)
			}
			if _, err := io.Copy(part, f); err != nil {
				return fmt.Errorf("write file part: %w", err)
			}
			if err := w.Close(); err != nil {
				return fmt.Errorf("close multipart: %w", err)
			}

			contentType := w.FormDataContentType()
			resp, err := api.FileControllerUploadFileWithBodyWithResponse(context.Background(), contentType, bytes.NewReader(buf.Bytes()))
			if err != nil {
				return fmt.Errorf("upload: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("upload: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("upload: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "URL: %s\n", resp.JSON200.FullUrl)
			return err
		},
	}
}

// NewFilesCmd returns the "files" parent command.
func NewFilesCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "files",
		Short: "File operations",
	}
	c.AddCommand(NewFilesUploadCmd(resolvePath, outputJSON))
	return c
}
