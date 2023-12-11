package service

import (
	"encoding/json"
	"strconv"
	"time"
	"tinyurl/pkg/database"
	"tinyurl/pkg/errors"
	redispkg "tinyurl/pkg/redis"

	"github.com/rs/zerolog/log"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type URLGenerateServicer interface {
	CreateURLInfo(urlInfo *UrlInfo) (*UrlInfo, error)
	GetURLInfo(id int) (*UrlInfo, error)
}

type UrlInfo struct {
	gorm.Model
	URL string
}

func (i *UrlInfo) GetJson() string {
	b, _ := json.Marshal(i)
	return string(b)
}

type URLGenerateService struct {
	db    *database.Database
	cache *redispkg.RedisCache
}

func NewURLGenerateService(db *database.Database, cache *redispkg.RedisCache) *URLGenerateService {
	return &URLGenerateService{
		db:    db,
		cache: cache,
	}
}

func (s *URLGenerateService) CreateURLInfo(urlInfo *UrlInfo) (*UrlInfo, error) {
	// 要有一個決定要序列產生器
	// 目前交給DB做決定 直接auto increase
	if val, err := s.cache.Get(urlInfo.URL).Result(); err == nil {
		err := json.Unmarshal([]byte(val), urlInfo)
		if err != nil {
			log.Warn().Msgf("CreateURLInfo : json Unmarshal failed, err: %v", err)
		}
	} else if err != redis.Nil {
		log.Warn().Msgf("CreateURLInfo : get data from cache failed, err: %v", err)
	}

	if err := s.db.Where(urlInfo).First(urlInfo).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Warn().Msgf("CreateURLInfo: Check url in DB failed, err: %v", err)
		}
	} else {
		if err := s.cache.Set(urlInfo.URL, urlInfo.GetJson(), 1*time.Hour).Err(); err != nil {
			log.Warn().Msgf("CreateURLInfo : set data to cache failed, err: %v", err)
		}
		return urlInfo, nil
	}

	if err := s.db.Create(urlInfo).Error; err != nil {
		return nil, err
	}

	return urlInfo, nil
}

func (s URLGenerateService) GetURLInfo(id int) (*UrlInfo, error) {
	urlInfo := &UrlInfo{}
	if val, err := s.cache.Get(strconv.Itoa(id)).Result(); err == nil {
		err := json.Unmarshal([]byte(val), urlInfo)
		if err != nil {
			log.Warn().Msgf("GetURLInfo : json Unmarshal failed, err: %v", err)
		} else {
			return urlInfo, nil
		}
	} else if err != redis.Nil {
		log.Warn().Msgf("GetURLInfo : get data from cache failed, err: %v", err)
	}
	queryInfo := &UrlInfo{}
	queryInfo.ID = uint(id)
	if err := s.db.Where(queryInfo).First(urlInfo).Error; err == gorm.ErrRecordNotFound {
		return nil, errors.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	if err := s.cache.Set(strconv.Itoa(id), urlInfo.GetJson(), 1*time.Hour).Err(); err != nil {
		log.Warn().Msgf("GetURLInfo : set data to cache failed, err: %v", err)
	}

	return urlInfo, nil
}
