package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_repo "github.com/Ray0427/url-shortener/mock/repo"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/mock/gomock"
)

func Test_urlController_PostUrl(t *testing.T) {
	// type fields struct {
	// 	urlRepo repo.UrlRepoInterface
	// }

	layout := "2006-01-02T15:04:05Z"
	sampleUrl := "https://ipinfo.io"
	sampleExpireAt, _ := time.Parse(layout, "2021-02-08T09:20:41Z")
	type args struct {
		// c          *gin.Context
		body       PostUrlsParam
		statusCode int
	}
	tests := []struct {
		name string
		// fields  fields
		prepare func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface)
		args    args
	}{
		// TODO: Add test cases.
		{
			name: "PostUrlSuccess",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				bodyJson, _ := json.Marshal(args.body)
				ctx.Request = httptest.NewRequest("POST", "/api/v1/urls", bytes.NewReader(bodyJson))
				ctx.Request.Header.Add("Content-Type", binding.MIMEJSON)
				mockRepo.EXPECT().CreateUrl(args.body.Url, sampleExpireAt).Return("3wedgpzLRq", nil)
			},
			args: args{
				body: PostUrlsParam{
					Url:      sampleUrl,
					ExpireAt: sampleExpireAt,
				},
				statusCode: http.StatusOK,
			},
		},
		{
			name: "InvalidExpireAt",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				bodyJson, _ := json.Marshal(args.body)
				ctx.Request = httptest.NewRequest("POST", "/api/v1/urls", bytes.NewReader(bodyJson))
				ctx.Request.Header.Add("Content-Type", binding.MIMEJSON)
			},
			args: args{
				body: PostUrlsParam{
					Url: sampleUrl,
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "InvalidUrl",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				bodyJson, _ := json.Marshal(args.body)
				ctx.Request = httptest.NewRequest("POST", "/api/v1/urls", bytes.NewReader(bodyJson))
				ctx.Request.Header.Add("Content-Type", binding.MIMEJSON)
			},
			args: args{
				body: PostUrlsParam{
					Url:      "ipinfo",
					ExpireAt: sampleExpireAt,
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "InternalServerError",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				bodyJson, _ := json.Marshal(args.body)
				ctx.Request = httptest.NewRequest("POST", "/api/v1/urls", bytes.NewReader(bodyJson))
				ctx.Request.Header.Add("Content-Type", binding.MIMEJSON)
				mockRepo.EXPECT().CreateUrl(args.body.Url, sampleExpireAt).Return("", &repo.InternalServerError{
					Message: "DB error",
				})
			},
			args: args{
				body: PostUrlsParam{
					Url:      sampleUrl,
					ExpireAt: sampleExpireAt,
				},
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockUrlRepoInterface(ctrl)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			if tt.prepare != nil {
				tt.prepare(ctx, tt.args, mockRepo)
			}
			uc := &urlController{
				urlRepo: mockRepo,
			}
			uc.PostUrl(ctx)

			if rec.Result().StatusCode != tt.args.statusCode {
				t.Errorf("StatusCode error, want %v, got %v", tt.args.statusCode, rec.Result().StatusCode)
			}
		})
	}
}

func Test_urlController_GetId(t *testing.T) {
	// type fields struct {
	// 	urlRepo repo.UrlRepoInterface
	// }
	type args struct {
		c          *gin.Context
		param      string
		statusCode int
	}
	tests := []struct {
		name    string
		prepare func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface)
		// fields fields
		args args
	}{
		// TODO: Add test cases.
		{
			name: "GetIdSuccess",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				ctx.Request = httptest.NewRequest("GET", "/"+args.param, nil)
				mockRepo.EXPECT().GetUrl(args.param).Return("https://ipinfo.io", nil)
			},
			args: args{
				param:      "3wedgpzLRq",
				statusCode: http.StatusMovedPermanently,
			},
		},
		{
			name: "BadRequestError",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				ctx.Request = httptest.NewRequest("GET", "/"+args.param, nil)
				mockRepo.EXPECT().GetUrl(args.param).Return("", &repo.BadRequestError{
					Message: "Invalid HashId",
				})
			},
			args: args{
				param:      "3wedgpzLR",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "NotFoundError",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				ctx.Request = httptest.NewRequest("GET", "/"+args.param, nil)
				mockRepo.EXPECT().GetUrl(args.param).Return("", &repo.NotFoundError{})
			},
			args: args{
				param:      "3wedgpzLRq",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "InternalServerError",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				ctx.Request = httptest.NewRequest("GET", "/"+args.param, nil)
				mockRepo.EXPECT().GetUrl(args.param).Return("", &repo.InternalServerError{
					Message: "DB error",
				})
			},
			args: args{
				param:      "3wedgpzLRq",
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "UnknownError",
			prepare: func(ctx *gin.Context, args args, mockRepo *mock_repo.MockUrlRepoInterface) {
				ctx.Request = httptest.NewRequest("GET", "/"+args.param, nil)
				mockRepo.EXPECT().GetUrl(args.param).Return("", errors.New("unknown"))
			},
			args: args{
				param:      "3wedgpzLRq",
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockUrlRepoInterface(ctrl)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			if tt.prepare != nil {
				tt.prepare(ctx, tt.args, mockRepo)
			}
			uc := &urlController{
				urlRepo: mockRepo,
			}
			ctx.Params = []gin.Param{{Key: "url_id", Value: tt.args.param}}
			uc.GetId(ctx)
			if rec.Result().StatusCode != tt.args.statusCode {
				t.Errorf("StatusCode error, want %v, got %v", tt.args.statusCode, rec.Result().StatusCode)
			}
		})
	}
}
