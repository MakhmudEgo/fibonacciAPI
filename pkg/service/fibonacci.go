package service

import (
	"fibonacciAPI/pkg/numbers"
	"fibonacciAPI/pkg/storage"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

type fibonacciService struct {
	rdb        *redis.Client
	prev, next int // Последние 2 числа Фибоначчи
	resN       int // С какого элемента ответ
	cash       int // С какого элемента кэшировать
}

func Fibonacci(rdb *redis.Client) Service {
	return &fibonacciService{rdb: rdb, prev: -1, next: -1}
}

func (f *fibonacciService) Execute(from, to int) ([]int, error) {
	if err := f.parse(from, to); err != nil {
		return nil, err
	}

	repo := storage.New(f.rdb)
	res, n, full, err := repo.Get(from, to)
	if err != nil {
		return nil, err
	}

	if full {
		return res, nil
	} else if res != nil {
		// no change capacity res
		if err = f.noFull(res, n); err != nil {
			return nil, err
		}
	} else {
		// new capacity res
		if err = f.noElemInCashForResponse(n); err != nil {
			return nil, err
		}
		f.resN = from - int(n) - 1
	}
	fib := numbers.FibonacciWithArgs(f.prev, f.next)
	res, err = fib.Generate(res, to-int(n))
	if err != nil {
		return nil, err
	}

	go func() {
		if err = repo.Set(res[f.cash:]); err != nil {
			log.Println(err)
		}
	}()
	return res[f.resN:], nil
}

func (f *fibonacciService) parse(from, to int) error {
	var err error
	if from < 1 {
		err = fmt.Errorf("not a natural number \"from\" – %d", from)
	} else if from > to {
		err = fmt.Errorf("\"to(%d)\" is less than \"from(%d)\"", to, from)
	}
	return err
}

func (f *fibonacciService) noFull(res []int, n int64) error {
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

func (f *fibonacciService) noElemInCashForResponse(n int64) error {
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