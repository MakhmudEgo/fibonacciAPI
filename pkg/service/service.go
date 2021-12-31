package service

import (
	"fibonacciAPI/pkg/fibonacci"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	REDIS_FIB_KEY = "fib:seq"
	REDIS_FIB_LOG = "fib:error:log"
)

type iService interface {
	Execute(int, int) ([]int, error)
	parse(int, int) error
	storage(int, int) ([]int, bool, error)
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

	res, n, full, err := f.storage(from, to)
	if err != nil {
		f.serviceErrors(err)
		return nil, err
	}
	var fib *fibonacci.Fibonacci
	if full {
		return res, nil
	} else if res != nil { // no change capacity res
		if len(res) > 1 {
			fib = fibonacci.NewFibonacciWithArgs(res[len(res)-2], res[len(res)-1], res)
		} else if n > 1 {
			prev, err := f.rdb.LIndex(f.rdb.Context(), REDIS_FIB_KEY, n-2).Int()
			if err != nil {
				f.serviceErrors(err)
				return nil, err
			}
			fib = fibonacci.NewFibonacciWithArgs(prev, res[len(res)-1], res)
		}
	} else if n != 0 { // new capacity res
		var prev, next *redis.StringCmd
		if _, err = f.rdb.Pipelined(f.rdb.Context(), func(p redis.Pipeliner) error {
			prev = p.LIndex(f.rdb.Context(), REDIS_FIB_KEY, n-2)
			next = p.LIndex(f.rdb.Context(), REDIS_FIB_KEY, n-1)
			return nil
		}); err != nil {
			f.serviceErrors(err)
			return nil, err
		}
		prevInt, prevErr := prev.Int()
		nextInt, nextErr := next.Int()
		if prevErr != nil || nextErr != nil {
			f.serviceErrors(err)
			return nil, err
		}
		if n > 1 {
			fib = fibonacci.NewFibonacciWithArgs(prevInt, nextInt, nil)
		} else {
			fib = fibonacci.NewFibonacciWithArgs(-1, nextInt, nil)
		}
	} else {
		fib = fibonacci.NewFibonacciWithArgs(-1, -1, nil)
	}
	res, err = fib.Generate(to - len(res))
	if err != nil {
		f.serviceErrors(err)
		return nil, err
	}
	return res, nil
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

// storage – []int: res; int: count; bool: full?, error _+_
func (f *Fibonacci) storage(from, to int) ([]int, int64, bool, error) {
	// get exists len seq Fibonacci
	n, err := f.rdb.LLen(f.rdb.Context(), REDIS_FIB_KEY).Result()
	if err != nil {
		f.serviceErrors(err)
		return nil, n, false, err
	}

	var res []int
	// if n >= to
	if n >= int64(from) {
		res = make([]int, 0, to-from+1)
		// todo:: proff ScanSlice vs stringSlice
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

func (f *Fibonacci) serviceErrors(err error) {
	f.rdb.LPush(f.rdb.Context(), REDIS_FIB_LOG, "service: "+err.Error()+time.Now().String())
}
