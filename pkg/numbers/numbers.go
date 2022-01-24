package numbers

import "math/big"

type Numbers interface {
	// Generate – генерация чисел
	Generate([]*big.Int, int64) ([]*big.Int, error)
}
