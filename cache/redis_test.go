package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestCache_SetUrl(t *testing.T) {
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
		prepare func(mock redismock.ClientMock)
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "SetSuccess",
			prepare: func(mock redismock.ClientMock, args ) {
				mock.ExpectSet("HASH_ID:"+"abc", Url{
					Url:      "https://ipinfo.io",
					ExpireAt: time.Now().AddDate(0, 0, 1),
				}, time.Hour).SetErr(nil)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			db, mock := redismock.NewClientMock()
			// f := fields{
			// 	client: db,
			// 	ctx:    ctx,
			// }
			if tt.prepare != nil {
				tt.prepare(mock)
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
