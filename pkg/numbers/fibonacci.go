package numbers

import (
	"fmt"
	"math/big"
)

type fibonacci struct {
	prev, next *big.Int
}

func (f *fibonacci) isValidArgs(n int) bool {
	if (f.prev.Cmp(big.NewInt(-1)) == 0 && f.next.Cmp(big.NewInt(-1)) == 0 ||
		f.prev.Cmp(big.NewInt(1)) == 0 && f.next.Cmp(big.NewInt(1)) == 0) && n > 0 {
		return true
	}
	if f.prev.Cmp(f.next) == 0 || n < 1 { /*||
		!Number(f.prev) || !Number(f.next)*/
		return false
	}
	return true
}

func Fibonacci() Numbers {
	return &fibonacci{prev: big.NewInt(-1), next: big.NewInt(-1)}
}

func FibonacciWithArgs(prev *big.Int, next *big.Int) Numbers {
	return &fibonacci{prev: prev, next: next}
}

// Generate – генератор последовательности
func (f *fibonacci) Generate(dst []*big.Int, n int) ([]*big.Int, error) {
	if !f.isValidArgs(n) {
		return nil, fmt.Errorf("bad args:\nprev – %d\nnext – %d\nn – %d",
			f.prev, f.next, n)
	}
	dst, i := f.init(dst, n)

	for ; i < n; i++ {
		f.prev.Add(f.prev, f.next)
		f.prev, f.next = f.next, f.prev
		v := new(big.Int)
		v.Set(f.next)
		dst = append(dst, v)
	}
	return dst, nil
}

// Number – проверка на число Фибоначчи
// todo::bug::из-за не идеального флоата
func Number(n *big.Int) bool {
	if n.Int64() == 0 {
		return true
	}
	x := &big.Float{}
	y := &big.Float{}

	x.SetInt(n).Mul(x, x).Mul(x, big.NewFloat(5))
	y.SetInt(n).Mul(y, y).Mul(y, big.NewFloat(5))

	tmp1, _ := x.Int(nil)
	tmp2, _ := y.Int(nil)
	tmp1.Add(tmp1, big.NewInt(4))
	tmp2.Sub(tmp2, big.NewInt(4))

	x.Add(x, big.NewFloat(4.0)).Sqrt(x)
	y.Sub(y, big.NewFloat(4.0)).Sqrt(y)

	r1, _ := x.Int(nil)
	r2, _ := y.Int(nil)
	return r1.Mul(r1, r1).Cmp(tmp1) == 0 || r2.Mul(r2, r2).Cmp(tmp2) == 0
}

func (f *fibonacci) init(dst []*big.Int, n int) ([]*big.Int, int) {
	var i int
	if f.prev.Cmp(big.NewInt(-1)) == 0 && f.next.Cmp(big.NewInt(-1)) == 0 {
		if n > 1 {
			dst = append(dst, big.NewInt(0), big.NewInt(1))
			i = 2
			f.prev, f.next = big.NewInt(0), big.NewInt(1)
		} else {
			dst = append(dst, big.NewInt(0))
			i = 1
		}
	} else if f.prev.Cmp(big.NewInt(-1)) == 0 {
		dst = append(dst, big.NewInt(1))
		f.prev, f.next = big.NewInt(0), big.NewInt(1)
		i = 1
	}

	return dst, i
}
