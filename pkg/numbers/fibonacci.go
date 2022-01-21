package numbers

import (
	"errors"
	"fmt"
	"math"
)

type fibonacci struct {
	prev, next int
}

func (f *fibonacci) isValidArgs(n int) bool {
	if (f.prev == -1 && f.next == -1 ||
		f.prev == 1 && f.next == 1) && n > 0 {
		return true
	}
	if f.prev == f.next || n < 1 ||
		!Number(f.prev) || !Number(f.next) {
		return false
	}
	return true
}

func Fibonacci() Numbers {
	return &fibonacci{prev: -1, next: -1}
}

func FibonacciWithArgs(prev int, next int) Numbers {
	return &fibonacci{prev: prev, next: next}
}

// Generate – генератор последовательности
func (f *fibonacci) Generate(dst []int, n int) ([]int, error) {
	if !f.isValidArgs(n) {
		return nil, fmt.Errorf("bad args:\nprev – %d\nnext – %d\nn – %d",
			f.prev, f.next, n)
	}
	dst, i := f.init(dst, n)

	for ; i < n; i++ {
		f.prev, f.next = f.next, f.prev+f.next
		if f.next < 0 {
			return nil, errors.New("overflow int")
		}
		dst = append(dst, f.next)
	}
	return dst, nil
}

// Number – проверка на число Фибоначчи
func Number(n int) bool {
	// todo:bug overflow int**
	x := math.Sqrt(5*math.Pow(float64(n), 2) + 4)
	y := math.Sqrt(5*math.Pow(float64(n), 2) - 4)
	if x == float64(int(x)) ||
		y == float64(int(y)) {
		return true
	}
	return false
}

func (f *fibonacci) init(dst []int, n int) ([]int, int) {
	var i int
	if f.prev == -1 && f.next == -1 {
		if n > 1 {
			dst = append(dst, 0, 1)
			i = 2
			f.prev, f.next = 0, 1
		} else {
			dst = append(dst, 0)
			i = 1
		}
	} else if f.prev == -1 {
		dst = append(dst, 1)
		f.prev, f.next = 0, 1
		i = 1
	}

	return dst, i
}
