package storage

import "math/big"

const REDIS_FIB_KEY = "fib:seq"

type Storage interface {
	// Set – Добавление в бд
	Set([]*big.Int) error
	// GetRange – выборка диапазона данных(например: f(1,90) => 1,2,3,4,5...)
	GetRange(int64, int64) ([]*big.Int, int64, bool, error)
	// Get – выборка произвольных данных(например: f(3,5,2,4,1) => 3,5,2,4,1...)
	Get(...int64) ([]string, []error)
}
