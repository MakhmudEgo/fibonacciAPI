package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"reflect"
	"testing"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln("No .env file found")
	}
}

func TestFibonacci_Execute(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	pong, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		t.Error(err)
	} else {
		log.Println(pong)
	}
	rdb.FlushDB(rdb.Context())

	type fields struct {
		rdb *redis.Client
	}
	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*big.Int
		wantErr bool
	}{
		{"from: 0, to: 10", fields{rdb}, args{0, 10}, nil, true},
		{"from: 1, to: 94", fields{rdb}, args{81, 81}, []*big.Int{big.NewInt(23416728348467685)}, false},
		{"from: 4, to: 94", fields{rdb}, args{4, 9}, []*big.Int{big.NewInt(2),
			big.NewInt(3),
			big.NewInt(5),
			big.NewInt(8),
			big.NewInt(13),
			big.NewInt(21),
		}, false},
		{"from: 114, to: 94", fields{rdb}, args{114, 94}, nil, true},
		{"from: -1, to: 34", fields{rdb}, args{114, 94}, nil, true},
		{"from: 11, to: 11", fields{rdb}, args{11, 11}, []*big.Int{big.NewInt(55)}, false},
		{"from: 11, to: 12", fields{rdb}, args{11, 12}, []*big.Int{big.NewInt(55), big.NewInt(89)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fibonacci(tt.fields.rdb)
			got, err := f.Execute(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
