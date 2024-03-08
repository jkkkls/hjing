package common

type ResData struct {
	Code  uint32                 `json:"code"`
	Msg   string                 `json:"msg"`
	Count int64                  `json:"total"`
	Data  []interface{}          `json:"data"`
	Map   map[string]interface{} `json:"map,omitempty"`

	Summary interface{}   `json:"summary,omitempty"`
	Data2   []interface{} `json:"data2"`
	Data3   []interface{} `json:"data3"`

	Succ         bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`

	Str     string `json:"str"`
	AppId   string `json:"appId"`
	AppName string `json:"appName"`
}
