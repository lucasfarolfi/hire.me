package entity

type ShortenedURL struct {
	Alias       string `json:"alias" gorm:"column:alias;unique"`
	Url         string `json:"url" gorm:"column:url"`
	AccessTimes int32  `json:"access_times,omitempty" gorm:"column:access_times"`
}
