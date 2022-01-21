package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
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
		want    []int
		wantErr bool
	}{
		{"from: 0, to: 10", fields{rdb}, args{0, 10}, nil, true},
		{"from: 1, to: 94", fields{rdb}, args{1, 94}, nil, true},
		{"from: 4, to: 94", fields{rdb}, args{4, 94}, nil, true},
		{"from: 114, to: 94", fields{rdb}, args{114, 94}, nil, true},
		{"from: -1, to: 34", fields{rdb}, args{114, 94}, nil, true},
		{"from: 11, to: 11", fields{rdb}, args{11, 11}, []int{55}, false},
		{"from: 11, to: 11", fields{rdb}, args{11, 12}, []int{55, 89}, false},
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
