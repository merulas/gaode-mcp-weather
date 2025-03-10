package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tung/mcp/internal/bean"
)

// AmapWeatherService 高德地图天气服务接口
type AmapWeatherService interface {
	GetHourlyWeather(location string) (*bean.WeatherResponse, error)
}

// amapWeatherService 高德地图天气服务实现
type amapWeatherService struct {
	apiKey        string
	baseURL       string
	locationCache *cache.Cache
	cacheDir      string
	cacheFile     string
}

// NewAmapWeatherService 创建新的高德地图天气服务
func NewAmapWeatherService(apiKey string) AmapWeatherService {
	// 创建缓存目录
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "amap_weather")
	cacheFile := filepath.Join(cacheDir, "location_cache.json")

	// 确保缓存目录存在
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// 初始化缓存
	c := cache.New(24*time.Hour, 1*time.Hour)

	// 尝试从文件加载缓存
	if _, err := os.Stat(cacheFile); err == nil {
		data, err := os.ReadFile(cacheFile)
		if err == nil {
			var cacheData map[string]string
			if err := json.Unmarshal(data, &cacheData); err == nil {
				for k, v := range cacheData {
					c.Set(k, v, cache.NoExpiration)
				}
			}
		}
	}

	return &amapWeatherService{
		apiKey:        apiKey,
		baseURL:       "https://restapi.amap.com/v3/weather/weatherInfo",
		locationCache: c,
		cacheDir:      cacheDir,
		cacheFile:     cacheFile,
	}
}

// GetHourlyWeather 获取每小时天气预报
func (s *amapWeatherService) GetHourlyWeather(location string) (*bean.WeatherResponse, error) {
	// 尝试从缓存获取城市编码
	cityCode, found := s.getCachedCityCode(location)
	if !found {
		// 如果缓存中没有，则使用输入的位置作为城市编码
		// 高德地图API支持城市名称、区域编码等多种方式
		cityCode = location
		// 缓存城市编码
		s.cacheCityCode(location, cityCode)
	}

	// 获取实况天气
	liveWeather, err := s.getLiveWeather(cityCode)
	if err != nil {
		return nil, fmt.Errorf("获取实况天气失败: %w", err)
	}

	// 获取天气预报
	forecastWeather, err := s.getForecastWeather(cityCode)
	if err != nil {
		return nil, fmt.Errorf("获取天气预报失败: %w", err)
	}

	// 构建响应
	response := s.buildWeatherResponse(liveWeather, forecastWeather)
	return response, nil
}

// getLiveWeather 获取实况天气
func (s *amapWeatherService) getLiveWeather(cityCode string) (*bean.AmapWeatherResponse, error) {
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("city", cityCode)
	params.Add("extensions", "base")
	params.Add("output", "JSON")

	resp, err := http.Get(s.baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResp bean.AmapWeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, err
	}

	if weatherResp.Status != "1" {
		return nil, fmt.Errorf("API返回错误: %s", weatherResp.Info)
	}

	return &weatherResp, nil
}

// getForecastWeather 获取天气预报
func (s *amapWeatherService) getForecastWeather(cityCode string) (*bean.AmapWeatherResponse, error) {
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("city", cityCode)
	params.Add("extensions", "all")
	params.Add("output", "JSON")

	resp, err := http.Get(s.baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResp bean.AmapWeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, err
	}

	if weatherResp.Status != "1" {
		return nil, fmt.Errorf("API返回错误: %s", weatherResp.Info)
	}

	return &weatherResp, nil
}

// buildWeatherResponse 构建天气响应
func (s *amapWeatherService) buildWeatherResponse(liveWeather, forecastWeather *bean.AmapWeatherResponse) *bean.WeatherResponse {
	if len(liveWeather.Lives) == 0 || len(forecastWeather.Forecasts) == 0 || len(forecastWeather.Forecasts[0].Casts) == 0 {
		return &bean.WeatherResponse{}
	}

	live := liveWeather.Lives[0]
	forecast := forecastWeather.Forecasts[0]

	// 构建当前天气状况
	temp, _ := strconv.ParseFloat(live.Temperature, 64)
	humidity, _ := strconv.Atoi(live.Humidity)

	currentConditions := bean.CurrentConditions{
		Temperature: bean.Temperature{
			Value: temp,
			Unit:  "C",
		},
		WeatherText:      live.Weather,
		RelativeHumidity: humidity,
		Precipitation:    strings.Contains(live.Weather, "雨") || strings.Contains(live.Weather, "雪"),
		ObservationTime:  live.ReportTime,
	}

	// 构建每小时天气预报
	// 注意：高德地图API只提供按天的预报，不提供每小时预报
	// 这里我们将按天的预报转换为模拟的每小时预报
	hourlyForecasts := make([]bean.HourlyForecast, 0)

	// 获取当前小时
	now := time.Now()
	currentHour := now.Hour()

	// 为今天和明天的每个小时创建预报
	for i := 0; i < 24; i++ {
		hour := (currentHour + i) % 24
		day := i / 24 // 0表示今天，1表示明天

		if day >= len(forecast.Casts) {
			break
		}

		cast := forecast.Casts[day]

		// 根据小时确定使用白天还是晚上的数据
		var weatherText, tempStr string
		var precipitationType, precipitationIntensity string
		var precipitationProbability int

		if hour >= 6 && hour < 18 {
			// 白天
			weatherText = cast.DayWeather
			tempStr = cast.DayTemp
		} else {
			// 晚上
			weatherText = cast.NightWeather
			tempStr = cast.NightTemp
		}

		// 根据天气文本判断降水类型和强度
		precipitationType = "None"
		precipitationIntensity = "None"
		precipitationProbability = 0

		if strings.Contains(weatherText, "雨") {
			precipitationType = "Rain"
			precipitationProbability = 80

			if strings.Contains(weatherText, "小雨") {
				precipitationIntensity = "Light"
				precipitationProbability = 60
			} else if strings.Contains(weatherText, "中雨") {
				precipitationIntensity = "Moderate"
				precipitationProbability = 80
			} else if strings.Contains(weatherText, "大雨") || strings.Contains(weatherText, "暴雨") {
				precipitationIntensity = "Heavy"
				precipitationProbability = 90
			} else {
				precipitationIntensity = "Light"
				precipitationProbability = 60
			}
		} else if strings.Contains(weatherText, "雪") {
			precipitationType = "Snow"
			precipitationProbability = 80

			if strings.Contains(weatherText, "小雪") {
				precipitationIntensity = "Light"
				precipitationProbability = 60
			} else if strings.Contains(weatherText, "中雪") {
				precipitationIntensity = "Moderate"
				precipitationProbability = 80
			} else if strings.Contains(weatherText, "大雪") || strings.Contains(weatherText, "暴雪") {
				precipitationIntensity = "Heavy"
				precipitationProbability = 90
			} else {
				precipitationIntensity = "Light"
				precipitationProbability = 60
			}
		}

		// 转换温度
		temp, _ := strconv.ParseFloat(tempStr, 64)

		// 创建相对时间字符串
		relativeTime := fmt.Sprintf("+%d hour", i+1)

		hourlyForecast := bean.HourlyForecast{
			RelativeTime: relativeTime,
			Temperature: bean.Temperature{
				Value: temp,
				Unit:  "C",
			},
			WeatherText:              weatherText,
			PrecipitationProbability: precipitationProbability,
			PrecipitationType:        precipitationType,
			PrecipitationIntensity:   precipitationIntensity,
		}

		hourlyForecasts = append(hourlyForecasts, hourlyForecast)

		// 只保留12小时的预报
		if len(hourlyForecasts) >= 12 {
			break
		}
	}

	return &bean.WeatherResponse{
		Location:          live.City,
		LocationKey:       live.Adcode,
		Country:           "中国",
		CurrentConditions: currentConditions,
		HourlyForecast:    hourlyForecasts,
	}
}

// getCachedCityCode 从缓存获取城市编码
func (s *amapWeatherService) getCachedCityCode(location string) (string, bool) {
	if cityCode, found := s.locationCache.Get(location); found {
		return cityCode.(string), true
	}
	return "", false
}

// cacheCityCode 缓存城市编码
func (s *amapWeatherService) cacheCityCode(location, cityCode string) {
	s.locationCache.Set(location, cityCode, cache.DefaultExpiration)

	// 将缓存保存到文件
	cacheData := make(map[string]string)
	items := s.locationCache.Items()
	for k, v := range items {
		cacheData[k] = v.Object.(string)
	}

	data, err := json.Marshal(cacheData)
	if err == nil {
		os.WriteFile(s.cacheFile, data, 0644)
	}
}
