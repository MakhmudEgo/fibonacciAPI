package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

const (
	REDIS_FIB_KEY = "fib:seq"
	REDIS_FIB_LOG = "fib:error:log"
)

type iService interface {
	Execute(int, int) ([]int, error)
	parse() error
}

type Fibonacci struct {
	rdb *redis.Client
}

func NewFibonacci(rdb *redis.Client) *Fibonacci {
	return &Fibonacci{rdb: rdb}
}

func (f *Fibonacci) Execute(from, to int) ([]int, error) {
	if err := f.parse(from, to); err != nil {
		f.serviceErrors(err)
		return nil, err
	}

	n, err := f.rdb.LLen(f.rdb.Context(), REDIS_FIB_KEY).Result()
	if err != nil {
		f.serviceErrors(err)
		return nil, err
	}

	var prev, next *redis.StringCmd
	if n > 1 {
		_, err := f.rdb.Pipelined(f.rdb.Context(), func(p redis.Pipeliner) error {
			prev = p.LIndex(f.rdb.Context(), REDIS_FIB_KEY, n-2)
			next = p.LIndex(f.rdb.Context(), REDIS_FIB_KEY, n-1)
			return nil
		})
		if err != nil {
			f.serviceErrors(err)
			return nil, err
		}
	} else {
		println("kekkkkkkkkk")
		return nil, nil
	}

	log.Println("prev:", prev.String())
	log.Println("next:", next.String())

	//get prev next -- llen-1 -- lindex -1
	// if llen == 1

	return nil, nil
}

func (f *Fibonacci) parse(from, to int) error {
	var err error
	if from < 1 {
		err = fmt.Errorf("not a natural number \"from\" – %d", from)
		f.serviceErrors(err)
	} else if from > to {
		err = fmt.Errorf("\"to(%d)\" is less than \"from(%d)\"", to, from)
		f.serviceErrors(err)
	}
	return err
}

func (f *Fibonacci) serviceErrors(err error) {
	f.rdb.LPush(f.rdb.Context(), REDIS_FIB_LOG, "service: "+err.Error()+time.Now().String())
}
