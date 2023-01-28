package backlog

import (
	"context"
	"fmt"
	"io"
)

// GetSharedFiles returns the meta data of shared files
func (c *Client) GetSharedFiles(projectID int, path string) ([]*SharedFile, error) {
	return c.GetSharedFilesContext(context.Background(), projectID, path)
}

// GetSharedFilesContext returns the meta data of shared files
func (c *Client) GetSharedFilesContext(ctx context.Context, projectID int, path string) ([]*SharedFile, error) {
	u := fmt.Sprintf("/api/v2/projects/%d/files/metadata%s", projectID, path)
	offset := 0
	sharedFiles := []*SharedFile{}
	size := 100

	for {
		u2, err := c.AddOptions(u, struct {
			Offset *int `url:"offset"`
			Count  *int `url:"count"`
		}{
			Offset: Int(0),
			Count:  Int(size),
		})
		if err != nil {
			return nil, err
		}

		req, err := c.NewRequest("GET", u2, nil)
		if err != nil {
			return nil, err
		}

		sf := []*SharedFile{}
		if err := c.Do(ctx, req, &sf); err != nil {
			return nil, err
		}
		sharedFiles = append(sharedFiles, sf...)
		if len(sf) != size {
			break
		}
		offset += size
	}

	return sharedFiles, nil
}

// GetFile writes the content to writer
func (c *Client) GetFile(projectID, fileID int, w io.Writer) error {
	return c.GetFileContext(context.Background(), projectID, fileID, w)
}

// GetFileContext writes the content to writer
func (c *Client) GetFileContext(ctx context.Context, projectID, fileID int, w io.Writer) error {
	u := fmt.Sprintf("/api/v2/projects/%d/files/%d", projectID, fileID)

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, w); err != nil {
		return err
	}
	return nil
}
