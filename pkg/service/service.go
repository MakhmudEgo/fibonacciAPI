package service

import (
	"fibonacciAPI/pkg/fibonacci"
	"fibonacciAPI/pkg/storage"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type iService interface {
	Execute(int, int) ([]int, error)
	parse(int, int) error
	noFull([]int, int64) error
	noElemInCashForResponse(int64) error
}

type Fibonacci struct {
	rdb        *redis.Client
	fib        *fibonacci.Fibonacci
	prev, next int // Последние 2 числа Фибоначчи
	resN       int // С какого элемента ответ
	cash       int // С какого элемента кэшировать
}

func NewFibonacci(rdb *redis.Client) *Fibonacci {
	return &Fibonacci{rdb: rdb, prev: -1, next: -1}
}

func (f *Fibonacci) Execute(from, to int) ([]int, error) {
	if err := f.parse(from, to); err != nil {
		f.serviceErrors(err)
		return nil, err
	}

	stg := storage.NewFibonacci(f.rdb)
	res, n, full, err := stg.Check(from, to)
	if err != nil {
		f.serviceErrors(err)
		return nil, err
	}

	if full {
		return res, nil
	} else if res != nil { // no change capacity res
		if err = f.noFull(res, n); err != nil {
			f.serviceErrors(err)
			return nil, err
		}
	} else { // new capacity res
		if err = f.noElemInCashForResponse(n); err != nil {
			f.serviceErrors(err)
			return nil, err
		}
		f.resN = from - int(n) - 1
	}
	f.fib = fibonacci.NewFibonacciWithArgs(f.prev, f.next, res)
	res, err = f.fib.Generate(to - int(n))
	if err != nil {
		f.serviceErrors(err)
		return nil, err
	}

	go stg.Set(res[f.cash:])
	return res[f.resN:], nil
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

func (f *Fibonacci) noFull(res []int, n int64) error {
	var err error
	f.cash = len(res)
	if len(res) > 1 {
		f.prev, f.next = res[len(res)-2], res[len(res)-1]
	} else if n > 1 {
		f.prev, err = f.rdb.LIndex(f.rdb.Context(), storage.REDIS_FIB_KEY, n-2).Int()
		if err != nil {
			return err
		}
		f.next = res[len(res)-1]
	}
	return err
}

func (f *Fibonacci) noElemInCashForResponse(n int64) error {
	if n == 0 {
		return nil
	}
	var err error
	var prevStr, nextStr *redis.StringCmd
	if _, err = f.rdb.Pipelined(f.rdb.Context(), func(p redis.Pipeliner) error {
		prevStr = p.LIndex(f.rdb.Context(), storage.REDIS_FIB_KEY, n-2)
		nextStr = p.LIndex(f.rdb.Context(), storage.REDIS_FIB_KEY, n-1)
		return nil
	}); err != nil {
		return err
	}
	var prevErr, nextErr error
	f.prev, prevErr = prevStr.Int()
	f.next, nextErr = nextStr.Int()
	if prevErr != nil || nextErr != nil {
		return err
	}
	if n == 1 {
		f.prev = -1
	}
	return err
}

func (f *Fibonacci) serviceErrors(err error) {
	f.rdb.LPush(f.rdb.Context(), storage.REDIS_FIB_LOG, "service: "+err.Error()+time.Now().String())
}
