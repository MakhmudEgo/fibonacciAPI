package service

import "math/big"

type Service interface {
	//Execute – исполнитель
	Execute(int64, int64) ([]*big.Int, error)
}
