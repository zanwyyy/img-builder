package builder

import "context"

// BuildRequest chứa thông tin đầu vào từ API chính gửi sang
type BuildRequest struct {
	RepoURL   string `json:"repo_url"`
	Branch    string `json:"branch"`
	ImageName string `json:"image_name"` // e.g: docker.io/user/my-app
	Tag       string `json:"tag"`        // e.g: latest hoặc commit hash
}

// GitClient định nghĩa cách lấy code về
type GitClient interface {
	Clone(ctx context.Context, repoURL, branch, targetPath string) error
}

// Compiler định nghĩa cách biến code thành binary
type Compiler interface {
	BuildGo(ctx context.Context, sourcePath string) (string, error)
}

// ImageManager định nghĩa cách tạo và đẩy Image
type ImageManager interface {
	CreateAndPush(ctx context.Context, binaryPath, imageName, tag string) error
}

// Service là bộ não điều phối toàn bộ quá trình
type Service struct {
	git      GitClient
	compiler Compiler
	oci      ImageManager
}

func NewService(g GitClient, c Compiler, o ImageManager) *Service {
	return &Service{git: g, compiler: c, oci: o}
}

func (s *Service) ProcessBuild(req BuildRequest) error {
	ctx := context.Background()

	if err := s.git.Clone(ctx, req.RepoURL, req.Branch, "/tmp/repo"); err != nil {
		return err
	}

	binaryPath, err := s.compiler.BuildGo(ctx, "/tmp/repo")
	if err != nil {
		return err
	}

	if err := s.oci.CreateAndPush(ctx, binaryPath, req.ImageName, req.Tag); err != nil {
		return err
	}

	return nil
}
