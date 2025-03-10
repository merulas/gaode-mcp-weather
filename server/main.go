package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tung/mcp/internal/handler"
	"github.com/tung/mcp/internal/logic"
	"github.com/tung/mcp/internal/service"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件，将使用系统环境变量")
	}

	// 获取API密钥
	apiKey := os.Getenv("AMAP_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 未设置AMAP_API_KEY环境变量")
	}

	// 创建服务
	weatherService := service.NewAmapWeatherService(apiKey)
	weatherLogic := logic.NewWeatherLogic(weatherService)
	weatherHandler := handler.NewWeatherHandler(weatherLogic)

	// 创建MCP处理器
	mcpHandler := handler.NewMCPHandler(weatherLogic)

	// 创建Gin路由
	router := gin.Default()

	// 注册路由
	weatherHandler.RegisterRoutes(router)
	mcpHandler.RegisterRoutes(router)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MCP天气服务启动在 :%s", port)
	log.Printf("标准API路径: http://localhost:%s/weather", port)
	log.Printf("Claude MCP API路径: http://localhost:%s/mcp", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
