package compiler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GoCompiler struct{}

// BuildGo thực hiện biên dịch mã nguồn Go thành file binary tĩnh (static binary)
func (c *GoCompiler) BuildGo(ctx context.Context, sourcePath string) (string, error) {
	// 1. Kiểm tra xem thư mục source có file go.mod không
	if _, err := os.Stat(filepath.Join(sourcePath, "go.mod")); os.IsNotExist(err) {
		return "", fmt.Errorf("không tìm thấy go.mod trong thư mục nguồn")
	}

	// 2. Định nghĩa tên và đường dẫn file binary đầu ra
	// Mình sẽ để tạm trong thư mục hệ thống, bạn có thể tùy chỉnh lại
	binaryName := "app-main"
	outputPath := filepath.Join(os.TempDir(), binaryName)

	// 3. Chuẩn bị lệnh build: go build -o <output> <source>
	// -ldflags="-s -w": Giảm kích thước binary bằng cách xóa bỏ debug thông tin
	cmd := exec.CommandContext(ctx, "go", "build", "-ldflags", "-s -w", "-o", outputPath, sourcePath)

	// 4. Thiết lập ENV để build static binary
	// CGO_ENABLED=0: Không dùng thư viện C của host OS (giúp chạy được trên Scratch/Distroless)
	// GOOS=linux: Đảm bảo chạy được trên Linux (môi trường container)
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS=linux",
		"GOARCH=amd64", // Bạn có thể lấy từ runtime.GOARCH nếu muốn build theo máy chủ
	)

	// 5. Thực thi lệnh và lấy log
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("lỗi khi biên dịch: %v\nOutput: %s", err, string(output))
	}

	// 6. Kiểm tra xem file đã thực sự được tạo ra chưa
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("lệnh build thành công nhưng không tìm thấy file binary tại %s", outputPath)
	}

	return outputPath, nil
}
