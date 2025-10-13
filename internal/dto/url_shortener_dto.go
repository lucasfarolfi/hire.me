package dto

type URLShortenerCreatedDTO struct {
	Alias      string         "json:alias"
	URL        string         "json:url"
	Statistics *StatisticsDTO "json:statistics"
}

type StatisticsDTO struct {
	TimeTaken string "json:time_taken"
}

func NewURLShortenerCreatedDTO(alias, url, timeTaken string) *URLShortenerCreatedDTO {
	return &URLShortenerCreatedDTO{alias, url, &StatisticsDTO{timeTaken}}
}
