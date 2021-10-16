package repo

import (
	"errors"
	"testing"
	"time"

	"github.com/Ray0427/url-shortener/cache"
	"github.com/Ray0427/url-shortener/config"
	mock_cache "github.com/Ray0427/url-shortener/mock/cache"
	mock_dao "github.com/Ray0427/url-shortener/mock/dao"
	"github.com/Ray0427/url-shortener/model"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func TestUrlRepo_CreateUrl(t *testing.T) {
	type fields struct {
		Dao    *mock_dao.MockUrlDaoInterface
		cache  *mock_cache.MockCacheInterface
		config config.Config
	}
	type args struct {
		url      string
		expireAt time.Time
	}
	layout := "2006-01-02T15:04:05Z"
	sampleUrl := "https://ipinfo.io"
	sampleExpireAt, _ := time.Parse(layout, "2021-02-08T09:20:41Z")
	sampleHashId := "94vz1DdAwM"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUrlDao := mock_dao.NewMockUrlDaoInterface(ctrl)
	mockCache := mock_cache.NewMockCacheInterface(ctrl)
	mockConfig := config.Config{}
	mockConfig.HashID.Salt = "test"
	mockConfig.HashID.MinLength = 10
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "CreateSuccess",
			prepare: func(f *fields) {
				f.Dao.EXPECT().Create(sampleUrl, sampleExpireAt).Return(model.Url{
					FullUrl:  sampleUrl,
					ExpireAt: sampleExpireAt,
				}, nil)
				f.cache.EXPECT().SetUrl(sampleHashId, cache.Url{
					Url:      sampleUrl,
					ExpireAt: sampleExpireAt,
				}).Return(nil)
			},
			args: args{
				url:      sampleUrl,
				expireAt: sampleExpireAt,
			},
			want:    sampleHashId,
			wantErr: false,
		},
		{
			name: "DatabaseError",
			prepare: func(f *fields) {
				f.Dao.EXPECT().Create(sampleUrl, sampleExpireAt).Return(model.Url{}, errors.New(""))
			},
			args: args{
				url:      sampleUrl,
				expireAt: sampleExpireAt,
			},
			wantErr: true,
		},
		{
			name: "CacheSuccess",
			prepare: func(f *fields) {
				f.Dao.EXPECT().Create(sampleUrl, sampleExpireAt).Return(model.Url{
					FullUrl:  sampleUrl,
					ExpireAt: sampleExpireAt,
				}, nil)
				f.cache.EXPECT().SetUrl(sampleHashId, cache.Url{
					Url:      sampleUrl,
					ExpireAt: sampleExpireAt,
				}).Return(errors.New(""))
			},
			args: args{
				url:      sampleUrl,
				expireAt: sampleExpireAt,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				Dao:    mockUrlDao,
				cache:  mockCache,
				config: mockConfig,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			ur := &UrlRepo{
				Dao:    f.Dao,
				cache:  f.cache,
				config: f.config,
			}
			got, err := ur.CreateUrl(tt.args.url, tt.args.expireAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlRepo.CreateUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UrlRepo.CreateUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlRepo_GetUrl(t *testing.T) {
	type fields struct {
		Dao    *mock_dao.MockUrlDaoInterface
		cache  *mock_cache.MockCacheInterface
		config config.Config
	}
	type args struct {
		hashID string
	}

	url1 := "https://ipinfo.io"
	hashId1 := "3wedgpzLRq"
	id1 := 1
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUrlDao := mock_dao.NewMockUrlDaoInterface(ctrl)
	mockCache := mock_cache.NewMockCacheInterface(ctrl)
	mockConfig := config.Config{}
	mockConfig.HashID.Salt = "test"
	mockConfig.HashID.MinLength = 10
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.

		{
			name: "GetCacheSuccess",
			prepare: func(f *fields) {
				expireAt := time.Now().AddDate(0, 0, 1)
				url := cache.Url{
					Url:      url1,
					ExpireAt: expireAt,
				}
				f.cache.EXPECT().GetUrl(hashId1, &cache.Url{}).SetArg(1, url).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    url1,
			wantErr: false,
		},
		{
			name: "GetDBSuccess",
			prepare: func(f *fields) {
				expireAt := time.Now().AddDate(0, 0, 1)
				modelUrl := model.Url{
					FullUrl:  url1,
					ExpireAt: expireAt,
				}
				f.Dao.EXPECT().Get(uint(id1)).Return(modelUrl, nil)
				url := cache.Url{}
				f.cache.EXPECT().GetUrl(hashId1, &url).Return(redis.Nil)
				f.cache.EXPECT().SetUrl(hashId1, cache.Url{
					Url:      url1,
					ExpireAt: expireAt,
				}).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    url1,
			wantErr: false,
		},
		{
			name: "GetCacheEmpty",
			prepare: func(f *fields) {
				expireAt := time.Now().AddDate(0, 0, -1)
				url := cache.Url{
					Url:      url1,
					ExpireAt: expireAt,
				}
				f.cache.EXPECT().GetUrl(hashId1, &cache.Url{}).SetArg(1, url).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetCacheExpire",
			prepare: func(f *fields) {
				// expireAt := time.Now().AddDate(0, 0, 1)
				f.cache.EXPECT().GetUrl(hashId1, &cache.Url{}).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "InvalidHashId",
			prepare: func(f *fields) {
				// expireAt := time.Now().AddDate(0, 0, 1)
				// modelUrl := model.Url{
				// 	FullUrl:  url1,
				// 	ExpireAt: expireAt,
				// }
				// f.Dao.EXPECT().Get(uint(id1)).Return(modelUrl, nil)
				url := cache.Url{}
				f.cache.EXPECT().GetUrl(hashId1, &url).Return(redis.Nil)
				f.config.HashID.MinLength = 1
				// f.cache.EXPECT().SetUrl(hashId1, cache.Url{
				// 	Url:      url1,
				// 	ExpireAt: expireAt,
				// }).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "DBNotFound",
			prepare: func(f *fields) {
				modelUrl := model.Url{}
				f.Dao.EXPECT().Get(uint(id1)).Return(modelUrl, gorm.ErrRecordNotFound)
				url := cache.Url{}
				f.cache.EXPECT().GetUrl(hashId1, &url).Return(redis.Nil)
				f.cache.EXPECT().SetUrl(hashId1, nil).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "DBError",
			prepare: func(f *fields) {
				modelUrl := model.Url{}
				f.Dao.EXPECT().Get(uint(id1)).Return(modelUrl, errors.New("DB error"))
				url := cache.Url{}
				f.cache.EXPECT().GetUrl(hashId1, &url).Return(redis.Nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetDBExpire",
			prepare: func(f *fields) {
				expireAt := time.Now().AddDate(0, 0, -1)
				modelUrl := model.Url{
					FullUrl:  url1,
					ExpireAt: expireAt,
				}
				f.Dao.EXPECT().Get(uint(id1)).Return(modelUrl, nil)
				url := cache.Url{}
				f.cache.EXPECT().GetUrl(hashId1, &url).Return(redis.Nil)
				f.cache.EXPECT().SetUrl(hashId1, cache.Url{
					Url:      url1,
					ExpireAt: expireAt,
				}).Return(nil)
			},
			args: args{
				hashID: hashId1,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				Dao:    mockUrlDao,
				cache:  mockCache,
				config: mockConfig,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			ur := &UrlRepo{
				Dao:    f.Dao,
				cache:  f.cache,
				config: f.config,
			}
			got, err := ur.GetUrl(tt.args.hashID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlRepo.GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UrlRepo.GetUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
