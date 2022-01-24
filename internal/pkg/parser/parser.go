package parser

import (
	"strconv"
)

func Parse(xStr, yStr string) (int64, int64, error) {
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return 0, 0, err
	}
	return int64(x), int64(y), nil
}
