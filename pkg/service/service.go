package service

type Service interface {
	Execute(int, int) ([]int, error)
}
