package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type c struct {
	cParams `json:"cache"`
	c       *cache.Cache
}

type cParams struct {
	DefaultExpiration int `json:"defaultExpiration"`
	CleanupInterval   int `json:"cleanupInterval"`
}

func (c *c) C() *cache.Cache {
	if c.c == nil {
		tmpc := cache.New(time.Duration(c.cParams.DefaultExpiration)*time.Minute,
			time.Duration(c.cParams.CleanupInterval)*time.Minute)

		c.c = tmpc
	}

	return c.c
}
