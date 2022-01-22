package storage

import "math/big"

const REDIS_FIB_KEY = "fib:seq"

type Storage interface {
	Set([]*big.Int) error
	Get(int, int) ([]*big.Int, int64, bool, error)
}
