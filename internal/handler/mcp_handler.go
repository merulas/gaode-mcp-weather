package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tung/mcp/internal/bean"
	"github.com/tung/mcp/internal/logic"
)

// MCPHandler MCP处理器接口
type MCPHandler interface {
	HandleMCPRequest(c *gin.Context)
	RegisterRoutes(router *gin.Engine)
}

// mcpHandler MCP处理器实现
type mcpHandler struct {
	weatherLogic logic.WeatherLogic
}

// NewMCPHandler 创建新的MCP处理器
func NewMCPHandler(weatherLogic logic.WeatherLogic) MCPHandler {
	return &mcpHandler{
		weatherLogic: weatherLogic,
	}
}

// RegisterRoutes 注册路由
func (h *mcpHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/mcp", h.HandleMCPRequest)
}

// HandleMCPRequest 处理MCP请求
func (h *mcpHandler) HandleMCPRequest(c *gin.Context) {
	var req bean.MCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, bean.NewMCPErrorResponse("无效的请求格式"))
		return
	}

	// 根据请求的名称分发到不同的处理函数
	switch req.Name {
	case "weather":
		h.handleWeatherRequest(c, req)
	default:
		c.JSON(http.StatusBadRequest, bean.NewMCPErrorResponse("未知的请求名称: "+req.Name))
	}
}

// handleWeatherRequest 处理天气请求
func (h *mcpHandler) handleWeatherRequest(c *gin.Context, req bean.MCPRequest) {
	// 解析参数
	var weatherReq bean.WeatherMCPRequest
	paramsJSON, err := json.Marshal(req.Parameters)
	if err != nil {
		c.JSON(http.StatusBadRequest, bean.NewMCPErrorResponse("参数解析失败"))
		return
	}

	if err := json.Unmarshal(paramsJSON, &weatherReq); err != nil {
		c.JSON(http.StatusBadRequest, bean.NewMCPErrorResponse("参数格式错误"))
		return
	}

	if weatherReq.Location == "" {
		c.JSON(http.StatusBadRequest, bean.NewMCPErrorResponse("缺少必要参数: location"))
		return
	}

	// 调用逻辑层获取天气数据
	response, err := h.weatherLogic.GetHourlyWeather(weatherReq.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bean.NewMCPErrorResponse(err.Error()))
		return
	}

	// 返回MCP格式的响应
	c.JSON(http.StatusOK, bean.NewMCPResponse(response))
}
