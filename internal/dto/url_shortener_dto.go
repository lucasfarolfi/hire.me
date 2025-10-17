package dto

type CreatedShortenedURLDTO struct {
	Alias      string         `json:"alias"`
	URL        string         `json:"url"`
	Statistics *StatisticsDTO `json:"statistics"`
}

type StatisticsDTO struct {
	TimeTaken string `json:"time_taken"`
}

func NewCreatedShortenedURLDTO(alias, url, timeTaken string) *CreatedShortenedURLDTO {
	return &CreatedShortenedURLDTO{alias, url, &StatisticsDTO{timeTaken}}
}

type ShortenedUrlRetrieveDTO struct {
	URL string `json:"url"`
}

type MostAcessedUrlDTO struct {
	URL         string `json:"url"`
	AccessTimes int    `json:"access_times"`
}
