package dao

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ray0427/url-shortener/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestUrlDao_Create(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	sampleUrl := "https://ipinfo.io"
	sampleExpireAt, _ := time.Parse(layout, "2021-02-08T09:20:41Z")
	// sampleUrlId := "3wedgpzLRq"
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		url      string
		expireAt time.Time
	}
	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		args    args
		want    model.Url
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "CreateSuccess",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `urls` (`full_url`,`expire_at`) VALUES (?,?)")).WithArgs(sampleUrl, sampleExpireAt).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)
				mock.ExpectCommit()
			},
			args: args{
				url:      sampleUrl,
				expireAt: sampleExpireAt,
			},
			want: model.Url{
				ID:       1,
				FullUrl:  sampleUrl,
				ExpireAt: sampleExpireAt,
			},
			wantErr: false,
		},
		{
			name: "CreateFailed",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `urls` (`full_url`,`expire_at`) VALUES (?,?)")).WithArgs(sampleUrl, sampleExpireAt).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			args: args{
				url:      sampleUrl,
				expireAt: sampleExpireAt,
			},
			want: model.Url{
				FullUrl:  sampleUrl,
				ExpireAt: sampleExpireAt,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("Failed to open mock sql db, got error: %v", err)
			}

			if db == nil {
				t.Error("mock db is null")
			}

			if mock == nil {
				t.Error("sqlmock is null")
			}
			gDB, err := gorm.Open(mysql.New(mysql.Config{
				SkipInitializeWithVersion: true,
				Conn:                      db,
			}), &gorm.Config{})
			if err != nil {
				t.Errorf("Failed to open gorm v2 db, got error: %v", err)
			}

			if gDB == nil {
				t.Error("gorm db is null")
			}
			defer db.Close()
			tt.prepare(mock)
			d := &UrlDao{
				DB: gDB,
			}
			got, err := d.Create(tt.args.url, tt.args.expireAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlDao.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UrlDao.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlDao_Get(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	sampleUrl := "https://ipinfo.io"
	sampleExpireAt, _ := time.Parse(layout, "2021-02-08T09:20:41Z")
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		fields  fields
		args    args
		want    model.Url
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "GetSuccess",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `urls` WHERE `urls`.`id` = ? ORDER BY `urls`.`id` LIMIT 1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "full_url", "expire_at"}).AddRow(1, sampleUrl, sampleExpireAt)).WillReturnError(nil)
			},
			args: args{
				id: 1,
			},
			want: model.Url{
				ID:       1,
				FullUrl:  sampleUrl,
				ExpireAt: sampleExpireAt,
			},
			wantErr: false,
		},
		{
			name: "GetFailed",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `urls` WHERE `urls`.`id` = ? ORDER BY `urls`.`id` LIMIT 1")).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
			},
			args: args{
				id: 1,
			},
			want:    model.Url{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("Failed to open mock sql db, got error: %v", err)
			}

			if db == nil {
				t.Error("mock db is null")
			}

			if mock == nil {
				t.Error("sqlmock is null")
			}
			gDB, err := gorm.Open(mysql.New(mysql.Config{
				SkipInitializeWithVersion: true,
				Conn:                      db,
			}), &gorm.Config{})
			if err != nil {
				t.Errorf("Failed to open gorm v2 db, got error: %v", err)
			}

			if gDB == nil {
				t.Error("gorm db is null")
			}
			defer db.Close()
			tt.prepare(mock)
			d := &UrlDao{
				DB: gDB,
			}
			got, err := d.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlDao.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UrlDao.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
