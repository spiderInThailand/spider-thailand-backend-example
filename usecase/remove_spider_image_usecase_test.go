package usecase

import (
	"context"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"spider-go/repository"
	"testing"

	"github.com/golang/mock/gomock"
)

type commonBuildStubRemoveSpiderImage struct {
	spiderRepo *mock_domain.MockSpiderRepository
}

func TestRemoveSpiderImageUsecase_RemoveSpiderImageBySpiderImageNameList(t *testing.T) {

	tempDir := t.TempDir()

	fileImagePathTemp := tempDir

	type args struct {
		spiderUUID        string
		spiderImageListRM []string
	}
	tests := []struct {
		name      string
		args      args
		buildStub func(*commonBuildStubRemoveSpiderImage)
		wantErr   bool
	}{
		{
			name: "remove_spider_image_success",
			args: args{
				spiderUUID: "SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21",
				spiderImageListRM: []string{
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_0a999e48-262e-43d7-8074-b0745ea40dac.jpeg",
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_ab1a74d6-ae14-4265-9f48-613685599d2b.jpeg",
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_e5d57f78-c616-40fd-a7a7-a5ccbf04e00d.jpeg",
				},
			},
			buildStub: remove_spider_image_success,
			wantErr:   false,
		},
		{
			name: "spider_uuid_note_found_error",
			args: args{
				spiderUUID: "SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21",
				spiderImageListRM: []string{
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_0a999e48-262e-43d7-8074-b0745ea40dac.jpeg",
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_ab1a74d6-ae14-4265-9f48-613685599d2b.jpeg",
					"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_e5d57f78-c616-40fd-a7a7-a5ccbf04e00d.jpeg",
				},
			},
			buildStub: spider_uuid_note_found_error,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpiderRepo := mock_domain.NewMockSpiderRepository(ctrl)

			commonStubs := commonBuildStubRemoveSpiderImage{
				spiderRepo: mockSpiderRepo,
			}

			tt.buildStub(&commonStubs)

			u := NewRemoveSpiderImageUsecase(commonStubs.spiderRepo, fileImagePathTemp)
			if err := u.RemoveSpiderImageBySpiderImageNameList(context.TODO(), tt.args.spiderUUID, tt.args.spiderImageListRM); (err != nil) != tt.wantErr {
				t.Errorf("RemoveSpiderImageUsecase.RemoveSpiderImageBySpiderImageNameList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func remove_spider_image_success(stub *commonBuildStubRemoveSpiderImage) {

	spiderInfo := model.SpiderInfo{
		ImageFile: []string{
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_0a999e48-262e-43d7-8074-b0745ea40dac.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_ab1a74d6-ae14-4265-9f48-613685599d2b.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_e5d57f78-c616-40fd-a7a7-a5ccbf04e00d.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_4068fe22-3b2c-4eb0-9e9a-ad2ea7e2ce3e.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_629ebfa4-0ae3-433b-aa77-7ffb295fae8c.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_0b06f645-552e-4950-a3b8-864efa94d692.jpeg",
		},
	}

	stub.spiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21"),
	).Return(&spiderInfo, nil)

	stub.spiderRepo.EXPECT().UpdateImageFileToSpiderInfo(
		gomock.Any(),
		gomock.Eq([]string{
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_4068fe22-3b2c-4eb0-9e9a-ad2ea7e2ce3e.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_629ebfa4-0ae3-433b-aa77-7ffb295fae8c.jpeg",
			"Image_SPIDER_826d1d15-e0a2-472e-b73d-fbc491a9da1b_0b06f645-552e-4950-a3b8-864efa94d692.jpeg",
		}),
		gomock.Eq("SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21"),
	).Return(nil)
}

func spider_uuid_note_found_error(stub *commonBuildStubRemoveSpiderImage) {

	stub.spiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_94fb3db9-cda2-4410-ab82-72424a5a1e21"),
	).Return(nil, repository.ErrorMongoNotFound)

}
