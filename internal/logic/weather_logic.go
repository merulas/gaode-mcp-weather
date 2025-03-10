package logic

import (
	"github.com/tung/mcp/internal/bean"
	"github.com/tung/mcp/internal/service"
)

// WeatherLogic 天气逻辑接口
type WeatherLogic interface {
	GetHourlyWeather(location string) (*bean.WeatherResponse, error)
}

// weatherLogic 天气逻辑实现
type weatherLogic struct {
	weatherService service.WeatherService
}

// NewWeatherLogic 创建新的天气逻辑
func NewWeatherLogic(weatherService service.WeatherService) WeatherLogic {
	return &weatherLogic{
		weatherService: weatherService,
	}
}

// GetHourlyWeather 获取每小时天气预报
func (l *weatherLogic) GetHourlyWeather(location string) (*bean.WeatherResponse, error) {
	return l.weatherService.GetHourlyWeather(location)
}
