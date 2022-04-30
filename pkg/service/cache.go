package service

import (
	goCache "github.com/patrickmn/go-cache"
)

var cache *goCache.Cache

func GetCache() *goCache.Cache {
	if cache == nil {
		cache = goCache.New(goCache.NoExpiration, goCache.NoExpiration)
	}
	return cache
}
