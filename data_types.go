package main

const (
	AuthKey              = "Authorization"
	AuthStringPattern    = "Bearer %v"
	ContentTypeKey       = "Content-Type"
	AppOrJSONContentType = "application/json"
	GetMethod            = "GET"
	MonthTimeFrame       = "month"
	ProdEnvKey           = "prod"

	GetBucketsUrl            = "https://api.runscope.com/buckets"
	GetEnvsUrlPattern        = "https://api.runscope.com/buckets/%v/environments"
	GetTestsUrlPattern       = "https://api.runscope.com/buckets/%v/tests?count=9999999"
	GetTestMetricsUrlPattern = "https://api.runscope.com/buckets/%v/tests/%v/metrics?timeframe=%v&environment_uuid=%v"

	ProxyUrlEnvVarKey          = "PROXY_URL"
	RunScopeSecretKeyEnvVarKey = "RUN_SCOPE_SECRET_KEY"
)

var (
	XLSXHeaders = [...]string{"Project Name", "Test Name",
		"Success Rate percent", "Avg Response Time sec", "Response Time sec 50th percentile",
		"Response Time sec 95th percentile", "Response Time sec 99th percentile", "Change from Last Period Response Time sec 50th percentile",
		"Change from Last Period Response Time sec 95th percentile", "Change from Last Period Response Time sec 99th percentile"}
)

type Bucket struct {
	BucketKey   string `json:"key"`
	ProjectName string `json:"name"`
	Environments
	Tests
	HasProductionData bool
}

type Buckets struct {
	Data []Bucket `json:"data"`
}

type Environment struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Environments struct {
	Data []Environment `json:"data"`
}

type Test struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	TestMetrics
}

type Tests struct {
	Data []Test `json:"data"`
}

type TestReponseTime struct {
	SuccessRatio  float64 `json:"success_ratio"`
	AvgRespTimeMs float64 `json:"avg_response_time_ms"`
}

type TestRespPeriodicMetrics struct {
	RespTime50ThPercentile float64 `json:"response_time_50th_percentile"`
	RespTime95ThPercentile float64 `json:"response_time_95th_percentile"`
	RespTime99ThPercentile float64 `json:"response_time_99th_percentile"`
	TotalTestRuns          float64 `json:"total_test_runs"`
}

type TestMetrics struct {
	ResponseTimes             []TestReponseTime       `json:"response_times"`
	ChangesFromLastTimePeriod TestRespPeriodicMetrics `json:"change_from_last_period"`
	ThisTimePeriod            TestRespPeriodicMetrics `json:"this_time_period"`
	EnvUUID                   string                  `json:"environment_uuid"`
	SuccessRate               float64
	AvgRespTimeMs             float64
}
