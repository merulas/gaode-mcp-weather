package bean

// WeatherRequest 天气请求参数
type WeatherRequest struct {
	Location string `json:"location" binding:"required"`
}

// Temperature 温度信息
type Temperature struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// CurrentConditions 当前天气状况
type CurrentConditions struct {
	Temperature      Temperature `json:"temperature"`
	WeatherText      string      `json:"weather_text"`
	RelativeHumidity int         `json:"relative_humidity"`
	Precipitation    bool        `json:"precipitation"`
	ObservationTime  string      `json:"observation_time"`
}

// HourlyForecast 每小时天气预报
type HourlyForecast struct {
	RelativeTime             string      `json:"relative_time"`
	Temperature              Temperature `json:"temperature"`
	WeatherText              string      `json:"weather_text"`
	PrecipitationProbability int         `json:"precipitation_probability"`
	PrecipitationType        string      `json:"precipitation_type"`
	PrecipitationIntensity   string      `json:"precipitation_intensity"`
}

// WeatherResponse 天气响应数据
type WeatherResponse struct {
	Location          string            `json:"location"`
	LocationKey       string            `json:"location_key"`
	Country           string            `json:"country"`
	CurrentConditions CurrentConditions `json:"current_conditions"`
	HourlyForecast    []HourlyForecast  `json:"hourly_forecast"`
}

// AccuWeatherLocationResponse AccuWeather位置响应
type AccuWeatherLocationResponse []struct {
	Key           string `json:"Key"`
	LocalizedName string `json:"LocalizedName"`
	Country       struct {
		LocalizedName string `json:"LocalizedName"`
	} `json:"Country"`
}

// AccuWeatherCurrentConditionsResponse AccuWeather当前天气状况响应
type AccuWeatherCurrentConditionsResponse []struct {
	LocalObservationDateTime string `json:"LocalObservationDateTime"`
	WeatherText              string `json:"WeatherText"`
	HasPrecipitation         bool   `json:"HasPrecipitation"`
	Temperature              struct {
		Metric struct {
			Value float64 `json:"Value"`
			Unit  string  `json:"Unit"`
		} `json:"Metric"`
	} `json:"Temperature"`
	RelativeHumidity int `json:"RelativeHumidity"`
}

// AccuWeatherHourlyForecastResponse AccuWeather每小时天气预报响应
type AccuWeatherHourlyForecastResponse []struct {
	DateTime                 string `json:"DateTime"`
	IconPhrase               string `json:"IconPhrase"`
	HasPrecipitation         bool   `json:"HasPrecipitation"`
	PrecipitationType        string `json:"PrecipitationType"`
	PrecipitationIntensity   string `json:"PrecipitationIntensity"`
	PrecipitationProbability int    `json:"PrecipitationProbability"`
	Temperature              struct {
		Value float64 `json:"Value"`
		Unit  string  `json:"Unit"`
	} `json:"Temperature"`
}
