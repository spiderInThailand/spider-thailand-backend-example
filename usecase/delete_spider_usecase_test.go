package usecase

import (
	"context"
	"fmt"
	"spider-go/config"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"testing"

	"github.com/golang/mock/gomock"
)

type commonBuildStub struct {
	spiderRepo *mock_domain.MockSpiderRepository
}

func TestDeleteSpiderInfoUsecase(t *testing.T) {

	tempDir := t.TempDir()
	config.C().File.FileImagePath = tempDir

	type arge struct {
		spiderUUID string
	}

	tc := []struct {
		name      string
		arge      arge
		buildStub func(*commonBuildStub)
		wantErr   bool
	}{
		{
			name: "successfull",
			arge: arge{
				spiderUUID: "SPIDER_565391ff-9197-47ce-b86e-311d7b901f53",
			},
			buildStub: successfull,
			wantErr:   false,
		},
		{
			name: "find_spider_info_not_found",
			arge: arge{
				spiderUUID: "SPIDER_565391ff-9197-47ce-b86e-311d7b901f53",
			},
			buildStub: find_spider_info_not_found,
			wantErr:   true,
		},
		{
			name: "delete_spider_info_failed",
			arge: arge{
				spiderUUID: "SPIDER_565391ff-9197-47ce-b86e-311d7b901f53",
			},
			buildStub: delete_spider_info_failed,
			wantErr:   true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpiderRepo := mock_domain.NewMockSpiderRepository(ctrl)

			commonStub := commonBuildStub{
				spiderRepo: mockSpiderRepo,
			}

			tt.buildStub(&commonStub)

			usecase := NewDeleteSpiderInfoUsecase(mockSpiderRepo)

			err := usecase.DeleteSpiderInfoUsecase(context.TODO(), tt.arge.spiderUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[TestDeleteSpiderInfoUsecase] failed, wantErr: %v, but got Err: %v, error: %+v", tt.wantErr, (err != nil), err)
			}
		})
	}

}

func successfull(stub *commonBuildStub) {
	spiderInfo := model.SpiderInfo{
		SpiderUUID: "SPIDER_565391ff-9197-47ce-b86e-311d7b901f53",
		ImageFile:  []string{},
	}

	stub.spiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_565391ff-9197-47ce-b86e-311d7b901f53"),
	).Return(&spiderInfo, nil)

	stub.spiderRepo.EXPECT().DeleteSpiderInfoWithSpiderUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_565391ff-9197-47ce-b86e-311d7b901f53"),
	).Return(nil)
}

func find_spider_info_not_found(stub *commonBuildStub) {
	stub.spiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_565391ff-9197-47ce-b86e-311d7b901f53"),
	).Return(nil, fmt.Errorf("MONGO NOT FOUND"))
}

func delete_spider_info_failed(stub *commonBuildStub) {
	spiderInfo := model.SpiderInfo{
		SpiderUUID: "SPIDER_565391ff-9197-47ce-b86e-311d7b901f53",
		ImageFile:  []string{},
	}

	stub.spiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_565391ff-9197-47ce-b86e-311d7b901f53"),
	).Return(&spiderInfo, nil)

	stub.spiderRepo.EXPECT().DeleteSpiderInfoWithSpiderUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_565391ff-9197-47ce-b86e-311d7b901f53"),
	).Return(fmt.Errorf("delete spider result is zero"))
}
