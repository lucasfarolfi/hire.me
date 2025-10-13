package entity

type ShortenedURL struct {
	Alias       string "json:alias"
	Url         string "json:url"
	AccessTimes int32  "json:access_times,omitempty"
}
