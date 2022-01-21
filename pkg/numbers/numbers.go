package numbers

type Numbers interface {
	Generate([]int, int) ([]int, error)
}
