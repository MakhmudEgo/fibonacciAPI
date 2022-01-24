package storage

import (
	"github.com/go-redis/redis/v8"
	"log"
	"math/big"
	"os"
)

type Redis redis.Client

func New() Storage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(pong)
	}
	return (*Redis)(rdb)
}

func (r *Redis) Set(seq []*big.Int) error {
	if _, err := (*redis.Client)(r).Pipelined((*redis.Client)(r).Context(), func(p redis.Pipeliner) error {
		for _, v := range seq {
			p.RPush((*redis.Client)(r).Context(), REDIS_FIB_KEY, v.String())
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetRange – []int: res; int: count; bool: full?, error _+_
func (r *Redis) GetRange(from int64, to int64) ([]*big.Int, int64, bool, error) {
	n, err := (*redis.Client)(r).LLen((*redis.Client)(r).Context(), REDIS_FIB_KEY).Result()
	if err != nil {
		return nil, n, false, err
	}

	var res []*big.Int
	if n >= from {
		res = make([]*big.Int, 0, to-from+1)
		resString := make([]string, 0, to-from+1)
		if err = (*redis.Client)(r).LRange((*redis.Client)(r).Context(),
			REDIS_FIB_KEY,
			from-1,
			to-1).ScanSlice(&resString); err != nil {
			return nil, n, false, err
		}
		for _, s := range resString {
			v := &big.Int{}
			v.SetString(s, 10)
			res = append(res, v)
		}
		if n >= to {
			return res, n, true, nil
		}
	}
	return res, n, false, nil
}

func (r *Redis) Get(nums ...int64) ([]string, []error) {
	result := make([]string, 0, len(nums))
	errs := make([]error, 0, len(nums))

	_, _ = (*redis.Client)(r).Pipelined((*redis.Client)(r).Context(), func(p redis.Pipeliner) error {
		for _, num := range nums {
			res, err := p.LIndex((*redis.Client)(r).Context(), REDIS_FIB_KEY, num-1).Result()
			result = append(result, res)
			errs = append(errs, err)
		}
		return nil
	})
	return result, errs
}
