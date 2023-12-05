package service

import (
	"log"
	"tinyurl/pkg/base58"
	"tinyurl/pkg/database"

	"gorm.io/gorm"
)

type URLGenerateServicer interface {
	CreateShortURL(url string) (urlKey string)
	GetShortURL(urlKey string) (url string)
}

type URLGenerateService struct {
	DB *database.Database
}

type UrlInfo struct {
	gorm.Model
	URL string
}

func NewURLGenerateService(db *database.Database) *URLGenerateService {
	return &URLGenerateService{DB: db}
}

func (s *URLGenerateService) CreateShortURL(url string) string {
	urlInfo := UrlInfo{URL: url}

	// 要有一個決定要序列產生器
	// 目前交給DB做決定 直接auto increase
	s.DB.Create(&urlInfo)
	shortURL := base58.EncodeFromInt(int(urlInfo.ID))
	// 得到後放進去DB 且檢查 cache 有沒有 如果有就砍掉

	return shortURL
}

func (s URLGenerateService) GetShortURL(urlKey string) string {

	decodeValue, err := base58.DecodeToInt(urlKey)
	if err != nil {
		log.Printf("the url_key is invaild, err:%v \n", err)
		return ""
	}
	urlInfo := &UrlInfo{}
	log.Printf("decodeValue: %v\n", decodeValue)
	r := s.DB.Where("id = ?",decodeValue).First(urlInfo)
	if r.Error != nil {
		return ""
	}

	return urlInfo.URL
}
