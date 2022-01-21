package numbers

import (
	"reflect"
	"testing"
)

func TestFibonacci_Check(t *testing.T) {
	tests := []struct {
		arg  int
		want bool
	}{
		{3, true},
		{0, true},
		{4, false},
		{12, false},
		{20, false},
		{21, true},
		{46368, true},
		{46367, false},
		// todo:bug:middle overflow int
		{160500643816367088, true},
		// todo:bug:middle overflow int
		{160500643816360000, false},
	}
	for _, tt := range tests {
		t.Run("TestFibonacci_Check", func(t *testing.T) {

			if got := Number(tt.arg); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_Generate(t *testing.T) {
	type fields struct {
		prev int
		next int
	}
	tests := []struct {
		name    string
		fields  fields
		arg     int
		want    []int
		wantErr bool
	}{
		{"prev: -1, next: -1", fields{-1, -1}, 3, []int{0, 1, 1}, false},
		{"prev: -1, next: -1", fields{-1, -1}, 3, []int{0, 1, 1}, false},
		{"prev: 1, next: 4", fields{-1, 4}, 3, nil, true},
		{"prev: -31, next: 4", fields{-1, 4}, 3, nil, true},
		{"prev: 5, next: 8", fields{5, 8}, 10, []int{13, 21, 34, 55, 89, 144, 233, 377, 610, 987}, false},
		{"prev: 0, next: 0", fields{0, 0}, 10, nil, true},
		{"prev: 1, next: 1", fields{1, 1}, 1, []int{2}, false},
		{"prev: 1, next: 1, arg: 94", fields{1, 1}, 94, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FibonacciWithArgs(tt.fields.prev, tt.fields.next)
			res := make([]int, 0, len(tt.want))
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
		prev int
		next int
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
				prev: tt.fields.prev,
				next: tt.fields.next,
			}
			if got := f.isValidArgs(tt.arg); got != tt.want {
				t.Errorf("isValidArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
