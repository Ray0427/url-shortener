package cache

import (
	"context"
	"testing"

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
	var ctx = context.TODO()
	db, mock := redismock.NewClientMock()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "SetSuccess",
			fields: fields{
				client: db,
				ctx:    ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
			}
			if err := c.SetUrl(tt.args.hashId, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Cache.SetUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
