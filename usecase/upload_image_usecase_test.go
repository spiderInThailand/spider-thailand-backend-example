package usecase

import (
	"context"
	"errors"
	"spider-go/config"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type commonStubsUploadImage struct {
	mockSpiderRepo *mock_domain.MockSpiderRepository
}

var (
	normal_spiderUUID = "SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21"
	normal_image      = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAYAAACp8Z5+AAABdWlDQ1BrQ0dDb2xvclNwYWNlRGlzcGxheVAzAAAokXWQvUvDUBTFT6tS0DqIDh0cMolD1NIKdnFoKxRFMFQFq1OafgltfCQpUnETVyn4H1jBWXCwiFRwcXAQRAcR3Zw6KbhoeN6XVNoi3sfl/Ticc7lcwBtQGSv2AijplpFMxKS11Lrke4OHnlOqZrKooiwK/v276/PR9d5PiFlNu3YQ2U9cl84ul3aeAlN//V3Vn8maGv3f1EGNGRbgkYmVbYsJ3iUeMWgp4qrgvMvHgtMunzuelWSc+JZY0gpqhrhJLKc79HwHl4plrbWD2N6f1VeXxRzqUcxhEyYYilBRgQQF4X/8044/ji1yV2BQLo8CLMpESRETssTz0KFhEjJxCEHqkLhz634PrfvJbW3vFZhtcM4v2tpCAzidoZPV29p4BBgaAG7qTDVUR+qh9uZywPsJMJgChu8os2HmwiF3e38M6Hvh/GMM8B0CdpXzryPO7RqFn4Er/QcXKWq8UwZBywAAAFZlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA5KGAAcAAAASAAAARKACAAQAAAABAAAABKADAAQAAAABAAAABAAAAABBU0NJSQAAAFNjcmVlbnNob3TxR2DXAAAB0mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNi4wLjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczpleGlmPSJodHRwOi8vbnMuYWRvYmUuY29tL2V4aWYvMS4wLyI+CiAgICAgICAgIDxleGlmOlBpeGVsWURpbWVuc2lvbj40PC9leGlmOlBpeGVsWURpbWVuc2lvbj4KICAgICAgICAgPGV4aWY6UGl4ZWxYRGltZW5zaW9uPjQ8L2V4aWY6UGl4ZWxYRGltZW5zaW9uPgogICAgICAgICA8ZXhpZjpVc2VyQ29tbWVudD5TY3JlZW5zaG90PC9leGlmOlVzZXJDb21tZW50PgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4K6FzRVQAAADtJREFUCB0di8ENwDAMAi8xO3bHzpQlmnwrU7uCB3cS47qXnYYBmYliBu5EyTeQStBQTYGefXrz36b5AH97GMU3efz3AAAAAElFTkSuQmCC"
)

func TestUploadImageSpiderUsecase(t *testing.T) {

	tmpDir := t.TempDir()
	config.C().File.FileImagePath = tmpDir
	config.C().RedisOption.Login.KeyFormat = "login_%s_%s"
	config.C().RedisOption.Login.TTL = time.Duration(12)

	type arge struct {
		spiderUUID        string
		listImageEncode64 []string
	}

	tc := []struct {
		name    string
		arge    arge
		stubs   func(stubs *commonStubsUploadImage)
		wantErr bool
	}{
		{
			name: "upload_image_success_case",
			arge: arge{
				spiderUUID: normal_spiderUUID,
				listImageEncode64: []string{
					normal_image,
					normal_image,
					normal_image,
				},
			},
			stubs:   upload_image_success_case,
			wantErr: false,
		},

		{
			name: "fail_spiderUUID_not_found_case",
			arge: arge{
				spiderUUID: "SPIDER_94fb3db9-cda2-4410-xxxx-72424a5a1e21",
				listImageEncode64: []string{
					normal_image,
					normal_image,
					normal_image,
				},
			},
			stubs:   fail_spiderUUID_not_found_case,
			wantErr: true,
		},
		{
			name: "fail_decode_image_case",
			arge: arge{
				spiderUUID: normal_spiderUUID,
				listImageEncode64: []string{
					"data:image/xxxx;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAYAAACp8Z5+AAABdWlDQ1BrQ0dDb2xvclNwYWNlRGlzcGxheVAzAAAokXWQvUvDUBTFT6tS0DqIDh0cMolD1NIKdnFoKxRFMFQFq1OafgltfCQpUnETVyn4H1jBWXCwiFRwcXAQRAcR3Zw6KbhoeN6XVNoi3sfl/Ticc7lcwBtQGSv2AijplpFMxKS11Lrke4OHnlOqZrKooiwK/v276/PR9d5PiFlNu3YQ2U9cl84ul3aeAlN//V3Vn8maGv3f1EGNGRbgkYmVbYsJ3iUeMWgp4qrgvMvHgtMunzuelWSc+JZY0gpqhrhJLKc79HwHl4plrbWD2N6f1VeXxRzqUcxhEyYYilBRgQQF4X/8044/ji1yV2BQLo8CLMpESRETssTz0KFhEjJxCEHqkLhz634PrfvJbW3vFZhtcM4v2tpCAzidoZPV29p4BBgaAG7qTDVUR+qh9uZywPsJMJgChu8os2HmwiF3e38M6Hvh/GMM8B0CdpXzryPO7RqFn4Er/QcXKWq8UwZBywAAAFZlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA5KGAAcAAAASAAAARKACAAQAAAABAAAABKADAAQAAAABAAAABAAAAABBU0NJSQAAAFNjcmVlbnNob3TxR2DXAAAB0mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNi4wLjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczpleGlmPSJodHRwOi8vbnMuYWRvYmUuY29tL2V4aWYvMS4wLyI+CiAgICAgICAgIDxleGlmOlBpeGVsWURpbWVuc2lvbj40PC9leGlmOlBpeGVsWURpbWVuc2lvbj4KICAgICAgICAgPGV4aWY6UGl4ZWxYRGltZW5zaW9uPjQ8L2V4aWY6UGl4ZWxYRGltZW5zaW9uPgogICAgICAgICA8ZXhpZjpVc2VyQ29tbWVudD5TY3JlZW5zaG90PC9leGlmOlVzZXJDb21tZW50PgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4K6FzRVQAAADtJREFUCB0di8ENwDAMAi8xO3bHzpQlmnwrU7uCB3cS47qXnYYBmYliBu5EyTeQStBQTYGefXrz36b5AH97GMU3efz3AAAAAElFTkSuQmCC",
				},
			},
			stubs:   fail_decode_image_case,
			wantErr: true,
		},
		{
			name: "fail_update_file_name_to_mongo_error_case",
			arge: arge{
				spiderUUID: normal_spiderUUID,
				listImageEncode64: []string{
					normal_image,
					normal_image,
					normal_image,
				},
			},
			stubs: fail_update_file_name_to_mongo_error_case,

			wantErr: true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpiderRepo := mock_domain.NewMockSpiderRepository(ctrl)

			commonStubsUploadImage := commonStubsUploadImage{
				mockSpiderRepo: mockSpiderRepo,
			}

			tt.stubs(&commonStubsUploadImage)

			usecase := NewUploadImageUsecase(mockSpiderRepo)

			err := usecase.UploadImageSpiderUsecase(
				context.TODO(),
				tt.arge.spiderUUID,
				tt.arge.listImageEncode64,
			)

			if tt.wantErr != (err != nil) {
				t.Errorf("[upload image usecase] want error is `%v` but got error: %+v", tt.wantErr, err)
			}

		})
	}
}

func upload_image_success_case(stubs *commonStubsUploadImage) {

	spiderInfo := model.SpiderInfo{
		SpiderUUID: normal_spiderUUID,
		ImageFile:  []string{},
	}

	stubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq(normal_spiderUUID),
	).Return(&spiderInfo, nil)

	stubs.mockSpiderRepo.EXPECT().UpdateImageFileToSpiderInfo(
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(normal_spiderUUID),
	).Return(nil)
}

func fail_spiderUUID_not_found_case(stubs *commonStubsUploadImage) {

	stubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_94fb3db9-cda2-4410-xxxx-72424a5a1e21"),
	).Return(nil, errors.New("MONGO_NOT_FOUND"))

}

func fail_decode_image_case(stubs *commonStubsUploadImage) {

	spiderInfo := model.SpiderInfo{
		SpiderUUID: normal_spiderUUID,
		ImageFile:  []string{},
	}

	stubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq(normal_spiderUUID),
	).Return(&spiderInfo, nil)
}

func fail_update_file_name_to_mongo_error_case(stubs *commonStubsUploadImage) {

	spiderInfo := model.SpiderInfo{
		SpiderUUID: normal_spiderUUID,
		ImageFile:  []string{},
	}

	stubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq(normal_spiderUUID),
	).Return(&spiderInfo, nil)

	stubs.mockSpiderRepo.EXPECT().UpdateImageFileToSpiderInfo(
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(normal_spiderUUID),
	).Return(ErrorMongoTechnicalFail)
}
