package storage

const REDIS_FIB_KEY = "fib:seq"

type Storage interface {
	Set([]int) error
	Get(int, int) ([]int, int64, bool, error)
}
