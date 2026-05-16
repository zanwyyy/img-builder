package git

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Client struct{}

func (c *Client) Clone(ctx context.Context, repoURL, branch, targetPath string) error {
	_, err := git.PlainCloneContext(ctx, targetPath, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1, // Shallow clone để chỉ lấy commit mới nhất, tăng tốc cực nhanh
		Progress:      os.Stdout,
	})
	return err
}
