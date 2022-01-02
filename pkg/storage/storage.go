package storage

import (
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	REDIS_FIB_KEY = "fib:seq"
	REDIS_FIB_LOG = "fib:error:log"
)

type iStorage interface {
	Set([]int) error
	Check(int, int) ([]int, int64, bool, error)
}

type Fibonacci struct {
	rdb *redis.Client
}

func (f *Fibonacci) Set(seq []int) error {
	if _, err := f.rdb.Pipelined(f.rdb.Context(), func(p redis.Pipeliner) error {
		for _, v := range seq {
			p.RPush(f.rdb.Context(), REDIS_FIB_KEY, v)
		}
		return nil
	}); err != nil {
		return err
	}
	//log.Println(11111)
	return nil
}

// Check – []int: res; int: count; bool: full?, error _+_
func (f *Fibonacci) Check(from int, to int) ([]int, int64, bool, error) {
	n, err := f.rdb.LLen(f.rdb.Context(), REDIS_FIB_KEY).Result()
	if err != nil {
		f.storageErrors(err)
		return nil, n, false, err
	}

	var res []int
	if n >= int64(from) {
		res = make([]int, 0, to-from+1)
		if err = f.rdb.LRange(f.rdb.Context(),
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

func NewFibonacci(rdb *redis.Client) *Fibonacci {
	return &Fibonacci{rdb: rdb}
}

func (f *Fibonacci) storageErrors(err error) {
	f.rdb.LPush(f.rdb.Context(), REDIS_FIB_LOG, "service: "+err.Error()+time.Now().String())
}
