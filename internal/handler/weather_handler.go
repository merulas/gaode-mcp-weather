package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tung/mcp/internal/bean"
	"github.com/tung/mcp/internal/logic"
)

// WeatherHandler 天气处理器接口
type WeatherHandler interface {
	GetHourlyWeather(c *gin.Context)
	RegisterRoutes(router *gin.Engine)
}

// weatherHandler 天气处理器实现
type weatherHandler struct {
	weatherLogic logic.WeatherLogic
}

// NewWeatherHandler 创建新的天气处理器
func NewWeatherHandler(weatherLogic logic.WeatherLogic) WeatherHandler {
	return &weatherHandler{
		weatherLogic: weatherLogic,
	}
}

// RegisterRoutes 注册路由
func (h *weatherHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/weather", h.GetHourlyWeather)
}

// GetHourlyWeather 获取每小时天气预报
func (h *weatherHandler) GetHourlyWeather(c *gin.Context) {
	var req bean.WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	response, err := h.weatherLogic.GetHourlyWeather(req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
