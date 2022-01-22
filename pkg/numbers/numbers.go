package numbers

import "math/big"

type Numbers interface {
	Generate([]*big.Int, int) ([]*big.Int, error)
}
