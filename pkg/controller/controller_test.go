package controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln("No .env file found")
	}
}

func TestFibonacci_ServeHTTP(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	pong, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		t.Fatal(err)
	} else {
		log.Println(pong)
	}
	type fields struct {
		rdb *redis.Client
	}

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	rcdr := httptest.NewRecorder()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"kek",
			fields{rdb},
			args{rcdr, httptest.NewRequest("GET", "/?from=1&to=13", nil)},
			`[0,1,1,2,3,5,8,13,21,34,55,89,144]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fibonacci{
				rdb: tt.fields.rdb,
			}

			f.ServeHTTP(tt.args.w, tt.args.r)
			resp, err := io.ReadAll(tt.args.w.Body)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(string(resp), tt.want+"\n") {
				t.Errorf("ServeHTTP() resp = %s, want %s", resp, tt.want)
			}
		})
	}
}
