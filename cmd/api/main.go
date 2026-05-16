package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zanwyyy/platform/internal/builder"
	"github.com/zanwyyy/platform/internal/compiler"
	"github.com/zanwyyy/platform/internal/deli/http"
	"github.com/zanwyyy/platform/internal/git"
	"github.com/zanwyyy/platform/internal/oci"
	// ... các import khác
)

func main() {
	r := gin.Default()

	// Khởi tạo các linh kiện (Dependency Injection)
	gitClient := &git.Client{}
	compiler := &compiler.GoCompiler{}
	ociManager := &oci.Manager{}
	builderSvc := builder.NewService(gitClient, compiler, ociManager)

	handler := http.NewHandler(builderSvc)

	// Route để Platform API gọi sang
	r.POST("/api/v1/build", handler.CreateBuild)

	r.Run(":8081") // Chạy trên port 8081 để không trùng với API chính (8080)
}
