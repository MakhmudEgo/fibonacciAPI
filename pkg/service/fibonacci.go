package service

import (
	"fibonacciAPI/pkg/numbers"
	"fibonacciAPI/pkg/storage"
	"fmt"
	"log"
	"math/big"
)

type fibonacciService struct {
	repo       storage.Storage
	prev, next *big.Int // Последние 2 числа Фибоначчи
	resN       int64    // С какого элемента ответ
	cash       int64    // С какого элемента кэшировать
}

func Fibonacci(repo storage.Storage) Service {
	return &fibonacciService{repo: repo, prev: big.NewInt(-1), next: big.NewInt(-1)}
}

func (f *fibonacciService) Execute(from, to int64) ([]*big.Int, error) {
	if err := f.parse(from, to); err != nil {
		return nil, err
	}
	log.Println(from, to)
	res, n, full, err := f.repo.GetRange(from, to)
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
		f.resN = from - n - 1
	}
	fib := numbers.FibonacciWithArgs(f.prev, f.next)
	res, err = fib.Generate(res, to-n)
	if err != nil {
		return nil, err
	}

	go func() {
		if err = f.repo.Set(res[f.cash:]); err != nil {
			log.Println(err)
		}
	}()
	return res[f.resN:], nil
}

func (f *fibonacciService) parse(from, to int64) error {
	var err error
	if from < 1 {
		err = fmt.Errorf("not a natural number \"from\" – %d", from)
	} else if from > to {
		err = fmt.Errorf("\"to(%d)\" is less than \"from(%d)\"", to, from)
	}
	return err
}

func (f *fibonacciService) noFull(res []*big.Int, n int64) error {
	var err error
	f.cash = int64(len(res))
	if len(res) > 1 {
		f.prev.Set(res[len(res)-2])
		f.next.Set(res[len(res)-1])
	} else if n > 1 {
		prev, err := f.repo.Get(n - 1)
		if err != nil {
			return err[0]
		}
		f.prev.SetString(prev[0], 10)
		f.next.Set(res[len(res)-1])
	}
	return err
}

func (f *fibonacciService) noElemInCashForResponse(n int64) error {
	if n == 1 {
		f.prev.SetInt64(-1)
	} else if n > 1 {
		nums, errs := f.repo.Get(n-1, n)
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
		f.prev.SetString(nums[0], 10)
		f.next.SetString(nums[1], 10)
	}
	return nil
}
