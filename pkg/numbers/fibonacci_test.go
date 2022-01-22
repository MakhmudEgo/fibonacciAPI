package numbers

import (
	"log"
	"math/big"
	"reflect"
	"testing"
)

func TestFibonacci_Check(t *testing.T) {
	tests := []struct {
		name string
		arg  int64
		want bool
	}{
		{"3", 3, true},
		{"0", 0, true},
		{"4", 4, false},
		{"12", 12, false},
		{"20", 20, false},
		{"21", 21, true},
		{"46368", 46368, true},
		{"46367", 46367, false},
		{"99194853094755497", 99194853094755497, true},
		{"160500643816367088", 160500643816367088, true},
		{"160500643816360000", 160500643816360000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := Number(big.NewInt(tt.arg)); got != tt.want {
				log.Println(got, tt.want)
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_Generate(t *testing.T) {
	type fields struct {
		prev int64
		next int64
	}
	tests := []struct {
		name    string
		fields  fields
		arg     int
		want    []*big.Int
		wantErr bool
	}{
		{"prev: -1, next: -1", fields{-1, -1}, 3, []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1)}, false},
		{"prev: -1, next: -1", fields{-1, -1}, 3, []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1)}, false},
		{"prev: 1, next: 4", fields{-1, 4}, 3, nil, true},
		{"prev: -31, next: 4", fields{-1, 4}, 3, nil, true},
		{"prev: 5, next: 8", fields{5, 8}, 10, []*big.Int{big.NewInt(13), big.NewInt(21), big.NewInt(34), big.NewInt(55), big.NewInt(89), big.NewInt(144), big.NewInt(233), big.NewInt(377), big.NewInt(610), big.NewInt(987)}, false},
		{"prev: 0, next: 0", fields{0, 0}, 10, nil, true},
		{"prev: 1, next: 1", fields{1, 1}, 1, []*big.Int{big.NewInt(2)}, false},
		{"prev: 1, next: 1, arg: 94", fields{1, 1}, 94, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FibonacciWithArgs(big.NewInt(tt.fields.prev), big.NewInt(tt.fields.next))
			res := make([]*big.Int, 0, len(tt.want))
			got, err := f.Generate(res, tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestFibonacci_SetDst(t *testing.T) {}

//func TestFibonacci_checkDst(t *testing.T) {}

func TestFibonacci_isValidArgs(t *testing.T) {
	type fields struct {
		prev int64
		next int64
	}
	tests := []struct {
		name   string
		fields fields
		arg    int
		want   bool
	}{
		{"prev: 0, next: 0", fields{0, 0}, 1, false},
		{"prev: 1, next: 1", fields{1, 1}, 1, true},
		{"prev: -1, next: -1", fields{-1, -1}, 1, true},
		{"prev: -1, next: -1", fields{-1, -1}, 0, false},
		{"prev: 3, next: 5", fields{3, 5}, 1, true},
		{"prev: -3, next: -5", fields{-3, -5}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fibonacci{
				prev: big.NewInt(tt.fields.prev),
				next: big.NewInt(tt.fields.next),
			}
			if got := f.isValidArgs(tt.arg); got != tt.want {
				t.Errorf("isValidArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
