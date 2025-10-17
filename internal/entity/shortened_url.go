package entity

type ShortenedURL struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Alias       string `gorm:"column:alias;unique"`
	Url         string `gorm:"column:url"`
	AccessTimes int32  `gorm:"column:access_times"`
}

func NewShortenedURL(alias, url string) *ShortenedURL {
	return &ShortenedURL{Alias: alias, Url: url, AccessTimes: 0}
}
