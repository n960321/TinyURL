package service

import (
	"log"
	"tinyurl/pkg/base58"
)

type URLGenerateServicer interface {
	CreateShortenedURL(url string) (urlKey string)
	GetShortenedURL(urlKey string) (url string)
}

type URLGenerateService struct {
}

type Url struct {
	Id  int
	Url string
}

func NewURLGenerateService() *URLGenerateService {
	return &URLGenerateService{}
}

func (s *URLGenerateService) CreateShortenedURL(url string) string {
	log.Panicf("not yet implemented")
	// 要有一個決定要序列產生器
	// 得到後放進去DB 且檢查 cache 有沒有 如果有就砍掉

	return ""
}

func (s URLGenerateService) GetShortenedURL(urlKey string) string {

	decodeValue, err := base58.DecodeToInt(urlKey)
	if err != nil {
		log.Printf("the url_key is invaild, err:%v \n", err)
		return ""
	}

	log.Printf("decodeValue: %v\n", decodeValue)
	log.Panicf("not yet implemented")

	return ""
}
