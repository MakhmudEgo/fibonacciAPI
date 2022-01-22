package service

import "math/big"

type Service interface {
	Execute(int, int) ([]*big.Int, error)
}
