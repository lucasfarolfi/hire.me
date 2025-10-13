package entity

type Shortener struct {
	Alias       string "json:alias"
	Url         string "json:url"
	AccessTimes int32  "json:access_times,omitempty"
}
