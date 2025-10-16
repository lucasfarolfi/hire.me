package entity

type ShortenedURL struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Alias       string `json:"alias" gorm:"column:alias;unique"`
	Url         string `json:"url" gorm:"column:url"`
	AccessTimes int32  `json:"access_times,omitempty" gorm:"column:access_times"`
}

func NewShortenedURL(alias, url string) *ShortenedURL {
	return &ShortenedURL{Alias: alias, Url: url, AccessTimes: 0}
}
