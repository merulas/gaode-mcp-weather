package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tung/mcp/internal/bean"
)

// WeatherService 天气服务接口
type WeatherService interface {
	GetHourlyWeather(location string) (*bean.WeatherResponse, error)
}

// weatherService 天气服务实现
type weatherService struct {
	apiKey        string
	baseURL       string
	locationCache *cache.Cache
	cacheDir      string
	cacheFile     string
}

// NewWeatherService 创建新的天气服务
func NewWeatherService(apiKey string) WeatherService {
	// 创建缓存目录
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "weather")
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

	return &weatherService{
		apiKey:        apiKey,
		baseURL:       "http://dataservice.accuweather.com",
		locationCache: c,
		cacheDir:      cacheDir,
		cacheFile:     cacheFile,
	}
}

// GetHourlyWeather 获取每小时天气预报
func (s *weatherService) GetHourlyWeather(location string) (*bean.WeatherResponse, error) {
	// 尝试从缓存获取位置键
	locationKey, found := s.getCachedLocationKey(location)
	if !found {
		// 如果缓存中没有，则从API获取
		var err error
		locationKey, err = s.getLocationKey(location)
		if err != nil {
			return nil, fmt.Errorf("获取位置键失败: %w", err)
		}
		// 缓存位置键
		s.cacheLocationKey(location, locationKey)
	}

	// 获取当前天气状况
	currentConditions, err := s.getCurrentConditions(locationKey)
	if err != nil {
		return nil, fmt.Errorf("获取当前天气状况失败: %w", err)
	}

	// 获取每小时天气预报
	hourlyForecast, err := s.getHourlyForecast(locationKey)
	if err != nil {
		return nil, fmt.Errorf("获取每小时天气预报失败: %w", err)
	}

	// 获取位置信息
	locationInfo, err := s.getLocationInfo(locationKey)
	if err != nil {
		return nil, fmt.Errorf("获取位置信息失败: %w", err)
	}

	// 构建响应
	response := &bean.WeatherResponse{
		Location:          locationInfo[0].LocalizedName,
		LocationKey:       locationKey,
		Country:           locationInfo[0].Country.LocalizedName,
		CurrentConditions: s.formatCurrentConditions(currentConditions),
		HourlyForecast:    s.formatHourlyForecast(hourlyForecast),
	}

	return response, nil
}

// getLocationKey 获取位置键
func (s *weatherService) getLocationKey(location string) (string, error) {
	url := fmt.Sprintf("%s/locations/v1/cities/search?apikey=%s&q=%s", s.baseURL, s.apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var locations bean.AccuWeatherLocationResponse
	if err := json.Unmarshal(body, &locations); err != nil {
		return "", err
	}

	if len(locations) == 0 {
		return "", fmt.Errorf("未找到位置: %s", location)
	}

	return locations[0].Key, nil
}

// getLocationInfo 获取位置信息
func (s *weatherService) getLocationInfo(locationKey string) (bean.AccuWeatherLocationResponse, error) {
	url := fmt.Sprintf("%s/locations/v1/%s?apikey=%s", s.baseURL, locationKey, s.apiKey)

	resp, err := http.Get(url)
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

	var location bean.AccuWeatherLocationResponse
	location = append(location, struct {
		Key           string `json:"Key"`
		LocalizedName string `json:"LocalizedName"`
		Country       struct {
			LocalizedName string `json:"LocalizedName"`
		} `json:"Country"`
	}{})

	if err := json.Unmarshal(body, &location[0]); err != nil {
		return nil, err
	}

	return location, nil
}

// getCurrentConditions 获取当前天气状况
func (s *weatherService) getCurrentConditions(locationKey string) (bean.AccuWeatherCurrentConditionsResponse, error) {
	url := fmt.Sprintf("%s/currentconditions/v1/%s?apikey=%s", s.baseURL, locationKey, s.apiKey)

	resp, err := http.Get(url)
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

	var currentConditions bean.AccuWeatherCurrentConditionsResponse
	if err := json.Unmarshal(body, &currentConditions); err != nil {
		return nil, err
	}

	return currentConditions, nil
}

// getHourlyForecast 获取每小时天气预报
func (s *weatherService) getHourlyForecast(locationKey string) (bean.AccuWeatherHourlyForecastResponse, error) {
	url := fmt.Sprintf("%s/forecasts/v1/hourly/12hour/%s?apikey=%s&metric=true", s.baseURL, locationKey, s.apiKey)

	resp, err := http.Get(url)
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

	var hourlyForecast bean.AccuWeatherHourlyForecastResponse
	if err := json.Unmarshal(body, &hourlyForecast); err != nil {
		return nil, err
	}

	return hourlyForecast, nil
}

// formatCurrentConditions 格式化当前天气状况
func (s *weatherService) formatCurrentConditions(currentConditions bean.AccuWeatherCurrentConditionsResponse) bean.CurrentConditions {
	if len(currentConditions) == 0 {
		return bean.CurrentConditions{}
	}

	current := currentConditions[0]
	return bean.CurrentConditions{
		Temperature: bean.Temperature{
			Value: current.Temperature.Metric.Value,
			Unit:  current.Temperature.Metric.Unit,
		},
		WeatherText:      current.WeatherText,
		RelativeHumidity: current.RelativeHumidity,
		Precipitation:    current.HasPrecipitation,
		ObservationTime:  current.LocalObservationDateTime,
	}
}

// formatHourlyForecast 格式化每小时天气预报
func (s *weatherService) formatHourlyForecast(hourlyForecast bean.AccuWeatherHourlyForecastResponse) []bean.HourlyForecast {
	result := make([]bean.HourlyForecast, len(hourlyForecast))

	for i, hour := range hourlyForecast {
		result[i] = bean.HourlyForecast{
			RelativeTime: fmt.Sprintf("+%d hour%s", i+1, func() string {
				if i+1 > 1 {
					return "s"
				}
				return ""
			}()),
			Temperature: bean.Temperature{
				Value: hour.Temperature.Value,
				Unit:  hour.Temperature.Unit,
			},
			WeatherText:              hour.IconPhrase,
			PrecipitationProbability: hour.PrecipitationProbability,
			PrecipitationType:        hour.PrecipitationType,
			PrecipitationIntensity:   hour.PrecipitationIntensity,
		}
	}

	return result
}

// getCachedLocationKey 从缓存获取位置键
func (s *weatherService) getCachedLocationKey(location string) (string, bool) {
	if value, found := s.locationCache.Get(location); found {
		return value.(string), true
	}
	return "", false
}

// cacheLocationKey 缓存位置键
func (s *weatherService) cacheLocationKey(location, locationKey string) {
	// 添加到内存缓存
	s.locationCache.Set(location, locationKey, cache.NoExpiration)

	// 保存到文件
	cacheData := make(map[string]string)

	// 尝试从文件加载现有缓存
	if _, err := os.Stat(s.cacheFile); err == nil {
		data, err := os.ReadFile(s.cacheFile)
		if err == nil {
			json.Unmarshal(data, &cacheData)
		}
	}

	// 更新缓存
	cacheData[location] = locationKey

	// 保存回文件
	data, err := json.Marshal(cacheData)
	if err == nil {
		os.WriteFile(s.cacheFile, data, 0644)
	}
}
