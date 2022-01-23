package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"reflect"
	"testing"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
func TestRedis_Get(t *testing.T) {

	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name      string
		r         *Redis
		args      args
		want      []*big.Int
		n         int64
		full      bool
		wantErr   bool
		data      []*big.Int
		setData   bool
		resetData bool
	}{
		{
			name: "empty", r: (*Redis)(r), args: args{1, 10},
			want: nil,
			n:    0, full: false, wantErr: false,
			data:    nil,
			setData: false,
		},
		{
			name: "empty", r: (*Redis)(r), args: args{1, 10},
			want: []*big.Int{big.NewInt(0), big.NewInt(1)},
			n:    2, full: false, wantErr: false,
			data:    []*big.Int{big.NewInt(0), big.NewInt(1)},
			setData: true,
		},
		{
			name: "empty", r: (*Redis)(r), args: args{1, 2},
			want: []*big.Int{big.NewInt(0), big.NewInt(1)},
			n:    2, full: true, wantErr: false,
			data:    []*big.Int{big.NewInt(0), big.NewInt(1)},
			setData: false,
		},
		{
			name: "empty", r: (*Redis)(r), args: args{1, 30},
			want: []*big.Int{big.NewInt(0), big.NewInt(1)},
			n:    2, full: false, wantErr: false,
			data:    []*big.Int{big.NewInt(0), big.NewInt(1)},
			setData: false,
		},
		{
			name: "empty", r: (*Redis)(r), args: args{1, 30},
			want: []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(2)},
			n:    4, full: false, wantErr: false,
			data:    []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(2)},
			setData: true, resetData: true,
		},
		{
			name: "empty", r: (*Redis)(r), args: args{1, 4},
			want: []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(2)},
			n:    4, full: true, wantErr: false,
			data:    []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(2)},
			setData: false, resetData: false,
		},
	}
	r.FlushDB(context.Background())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resetData {
				r.FlushDB(context.Background())
			}
			if tt.setData {
				if err := tt.r.Set(tt.data); err != nil {
					t.Errorf(err.Error())
				}
			}
			got, got1, got2, err := tt.r.Get(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.n {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.n)
			}
			if got2 != tt.full {
				t.Errorf("Get() got2 = %v, want %v", got2, tt.full)
			}
		})
	}
}
