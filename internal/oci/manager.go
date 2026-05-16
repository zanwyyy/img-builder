package oci

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

type Manager struct {
	// Bạn có thể thêm cấu hình registry mặc định ở đây nếu muốn
	BaseImage string
}

func NewManager() *Manager {
	return &Manager{
		// Sử dụng distroless làm base mặc định vì siêu nhẹ và an toàn
		BaseImage: "gcr.io/distroless/static-debian11",
	}
}

// CreateAndPush đóng gói binary thành image và push lên registry
func (m *Manager) CreateAndPush(ctx context.Context, binaryPath, imageName, tag string) error {
	// 1. Phân tích reference của image (ví dụ: index.docker.io/myuser/myapp:latest)
	fullImageName := fmt.Sprintf("%s:%s", imageName, tag)
	ref, err := name.ParseReference(fullImageName)
	if err != nil {
		return fmt.Errorf("invalid image name: %v", err)
	}

	// 2. Kéo Base Image từ remote về (dạng object trong memory)
	baseRef, _ := name.ParseReference(m.BaseImage)
	baseImg, err := remote.Image(baseRef, remote.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to fetch base image: %v", err)
	}

	// 3. Tạo một Layer mới từ file binary đã compile
	// Chúng ta sẽ đặt file binary vào đường dẫn "/app" bên trong Container
	layer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return os.Open(binaryPath)
	})
	if err != nil {
		return fmt.Errorf("failed to create layer from binary: %v", err)
	}

	// 4. "Chồng" layer binary lên base image
	// Mutate.AppendLayers sẽ tạo ra một Image object mới với layer binary ở trên cùng
	img, err := mutate.AppendLayers(baseImg, layer)
	if err != nil {
		return fmt.Errorf("failed to append layer: %v", err)
	}

	// 5. Cấu hình Entrypoint (Lệnh chạy khi khởi động container)
	cfg, err := img.ConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get config file: %v", err)
	}

	// Copy cấu hình để chỉnh sửa
	newCfg := cfg.Config.DeepCopy()
	newCfg.Entrypoint = []string{"/app"} // Chỉ định chạy file binary tại /app

	// Cập nhật lại config cho Image
	img, err = mutate.Config(img, *newCfg)
	if err != nil {
		return fmt.Errorf("failed to mutate config: %v", err)
	}

	// 6. Push image lên Registry với cơ chế Auth chuẩn
	// authn.DefaultKeychain sẽ tự động tìm token/password trong ~/.docker/config.json
	err = remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain), remote.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to push image to registry: %v", err)
	}

	return nil
}
