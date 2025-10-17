package dto

import "github.com/lucasfarolfi/hire.me/internal/entity"

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

func NewMostAcessedUrlsDTO(shortUrls []entity.ShortenedURL) []MostAcessedUrlDTO {
	dto := make([]MostAcessedUrlDTO, 0, len(shortUrls))
	for _, su := range shortUrls {
		dto = append(dto, MostAcessedUrlDTO{
			URL:         su.Url,
			AccessTimes: int(su.AccessTimes),
		})
	}
	return dto
}
