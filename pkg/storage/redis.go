package storage

import "github.com/go-redis/redis/v8"

type Redis redis.Client

func New(rdb *redis.Client) Storage {
	return (*Redis)(rdb)
}

func (r *Redis) Set(seq []int) error {
	if _, err := (*redis.Client)(r).Pipelined((*redis.Client)(r).Context(), func(p redis.Pipeliner) error {
		for _, v := range seq {
			p.RPush((*redis.Client)(r).Context(), REDIS_FIB_KEY, v)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Get – []int: res; int: count; bool: full?, error _+_
func (r *Redis) Get(from int, to int) ([]int, int64, bool, error) {
	n, err := (*redis.Client)(r).LLen((*redis.Client)(r).Context(), REDIS_FIB_KEY).Result()
	if err != nil {
		return nil, n, false, err
	}

	var res []int
	if n >= int64(from) {
		res = make([]int, 0, to-from+1)
		if err = (*redis.Client)(r).LRange((*redis.Client)(r).Context(),
			REDIS_FIB_KEY,
			int64(from-1),
			int64(to-1)).ScanSlice(&res); err != nil {
			return nil, n, false, err
		}
		if n >= int64(to) {
			return res, n, true, nil
		}
	}
	return res, n, false, nil
}
