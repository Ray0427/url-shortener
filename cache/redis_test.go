package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestCache_SetUrl(t *testing.T) {
	type args struct {
		hashId string
		url    interface{}
	}

	tests := []struct {
		name    string
		prepare func(mock redismock.ClientMock, args args)
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "SetSuccess",
			prepare: func(mock redismock.ClientMock, args args) {
				val, _ := json.Marshal(args.url)
				mock.ExpectSet("HASH_ID:"+args.hashId, val, time.Hour).SetVal("ok")
			},
			args: args{
				hashId: "abc",
				url: Url{
					Url:      "https://ipinfo.io",
					ExpireAt: time.Now().AddDate(0, 0, 1),
				},
			},
			wantErr: false,
		},
		{
			name: "SetEmptySuccess",
			prepare: func(mock redismock.ClientMock, args args) {
				val, _ := json.Marshal(args.url)
				mock.ExpectSet("HASH_ID:"+args.hashId, val, time.Hour).SetVal("ok")
			},
			args: args{
				hashId: "abc",
				url:    nil,
			},
			wantErr: false,
		},
		{
			name: "SetFailed",
			prepare: func(mock redismock.ClientMock, args args) {
				val, _ := json.Marshal(args.url)
				mock.ExpectSet("HASH_ID:"+args.hashId, val, time.Hour).SetErr(errors.New("set failed"))
			},
			args: args{
				hashId: "abc",
				url:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			db, mock := redismock.NewClientMock()
			if tt.prepare != nil {
				tt.prepare(mock, tt.args)
			}
			c := &Cache{
				client: db,
				ctx:    ctx,
			}
			if err := c.SetUrl(tt.args.hashId, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Cache.SetUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_GetUrl(t *testing.T) {
	type fields struct {
		client *redis.Client
		ctx    context.Context
	}
	type args struct {
		hashId string
		url    interface{}
	}
	tests := []struct {
		name    string
		prepare func(mock redismock.ClientMock, args args)
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "GetSuccess",
			args: args{
				hashId: "abc",
				url: Url{
					Url:      "https://ipinfo.io",
					ExpireAt: time.Now().AddDate(0, 0, 1),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			db, mock := redismock.NewClientMock()
			if tt.prepare != nil {
				tt.prepare(mock, tt.args)
			}
			c := &Cache{
				client: db,
				ctx:    ctx,
			}
			if err := c.GetUrl(tt.args.hashId, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Cache.GetUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
