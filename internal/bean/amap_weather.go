package bean

// AmapWeatherRequest 高德地图天气请求参数
type AmapWeatherRequest struct {
	City       string `json:"city"`       // 城市编码
	Key        string `json:"key"`        // 用户在高德地图官网申请的key
	Output     string `json:"output"`     // 可选值：JSON,XML
	Extensions string `json:"extensions"` // 可选值：base/all，base:返回实况天气，all:返回预报天气
}

// AmapWeatherResponse 高德地图天气响应
type AmapWeatherResponse struct {
	Status    string                `json:"status"`              // 返回状态
	Count     string                `json:"count"`               // 返回结果数目
	Info      string                `json:"info"`                // 返回的状态信息
	InfoCode  string                `json:"infocode"`            // 返回状态说明
	Lives     []AmapLiveWeather     `json:"lives,omitempty"`     // 实况天气数据
	Forecasts []AmapForecastWeather `json:"forecasts,omitempty"` // 预报天气数据
}

// AmapLiveWeather 高德地图实况天气
type AmapLiveWeather struct {
	Province      string `json:"province"`      // 省份名
	City          string `json:"city"`          // 城市名
	Adcode        string `json:"adcode"`        // 区域编码
	Weather       string `json:"weather"`       // 天气现象
	Temperature   string `json:"temperature"`   // 实时气温，单位：摄氏度
	WindDirection string `json:"winddirection"` // 风向
	WindPower     string `json:"windpower"`     // 风力级别
	Humidity      string `json:"humidity"`      // 空气湿度
	ReportTime    string `json:"reporttime"`    // 数据发布的时间
}

// AmapForecastWeather 高德地图预报天气
type AmapForecastWeather struct {
	City       string            `json:"city"`       // 城市名称
	Adcode     string            `json:"adcode"`     // 区域编码
	Province   string            `json:"province"`   // 省份名称
	Reporttime string            `json:"reporttime"` // 预报发布时间
	Casts      []AmapWeatherCast `json:"casts"`      // 预报数据
}

// AmapWeatherCast 高德地图天气预报
type AmapWeatherCast struct {
	Date         string `json:"date"`         // 日期
	Week         string `json:"week"`         // 星期几
	DayWeather   string `json:"dayweather"`   // 白天天气现象
	NightWeather string `json:"nightweather"` // 晚上天气现象
	DayTemp      string `json:"daytemp"`      // 白天温度
	NightTemp    string `json:"nighttemp"`    // 晚上温度
	DayWind      string `json:"daywind"`      // 白天风向
	NightWind    string `json:"nightwind"`    // 晚上风向
	DayPower     string `json:"daypower"`     // 白天风力
	NightPower   string `json:"nightpower"`   // 晚上风力
}
