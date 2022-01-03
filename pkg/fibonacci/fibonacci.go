package fibonacci

import (
	"errors"
	"fmt"
	"math"
)

type iFibonacci interface {
	Generate(int) ([]int, error)
	Check(int) bool
	checkDst(int)
	isValidArgs(int) bool
	init(int) int
}

type Fibonacci struct {
	prev, next int
	dst        []int
}

func (f *Fibonacci) SetDst(dst []int) {
	f.dst = dst
}

func (f *Fibonacci) isValidArgs(n int) bool {
	if (f.prev == -1 && f.next == -1 ||
		f.prev == 1 && f.next == 1) && n > 0 {
		return true
	}
	if f.prev == f.next || n < 1 ||
		!f.Check(f.prev) || !f.Check(f.next) {
		return false
	}
	return true
}

func NewFibonacci(dst []int) *Fibonacci {
	return &Fibonacci{prev: -1, next: -1, dst: dst}
}

func NewFibonacciWithArgs(prev int, next int, dst []int) *Fibonacci {
	return &Fibonacci{prev: prev, next: next, dst: dst}
}

// Generate – генератор последовательности
func (f *Fibonacci) Generate(n int) ([]int, error) {
	if !f.isValidArgs(n) {
		return nil, fmt.Errorf("bad args:\nprev – %d\nnext – %d\nn – %d",
			f.prev, f.next, n)
	}
	f.checkDst(n)
	i := f.init(n)

	for ; i < n; i++ {
		f.prev, f.next = f.next, f.prev+f.next
		if f.next < 0 {
			return nil, errors.New("overflow int")
		}
		f.dst = append(f.dst, f.next)
	}
	return f.dst, nil
}

// Check – проверка на число Фибоначчи
func (f *Fibonacci) Check(n int) bool {
	// todo:bug overflow int**
	x := math.Sqrt(5*math.Pow(float64(n), 2) + 4)
	y := math.Sqrt(5*math.Pow(float64(n), 2) - 4)
	if x == float64(int(x)) ||
		y == float64(int(y)) {
		return true
	}
	return false
}

func (f *Fibonacci) checkDst(n int) {
	if f.dst == nil {
		f.dst = make([]int, 0, n)
	}
}

func (f *Fibonacci) init(n int) int {
	var i int
	if f.prev == -1 {
		if f.next == -1 {
			if n > 1 {
				f.dst = append(f.dst, 0, 1)
				i = 2
				f.prev, f.next = 0, 1
			} else {
				f.dst = append(f.dst, 0)
				i = 1
			}
		} else {
			f.dst = append(f.dst, 1)
			f.prev, f.next = 0, 1
			i = 1
		}
	}
	return i
}
