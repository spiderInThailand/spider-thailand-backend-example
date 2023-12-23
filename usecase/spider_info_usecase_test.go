package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"os"
	"path"
	"reflect"
	"spider-go/config"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"spider-go/repository"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

type mockStubsSpiderInfoUsecase struct {
	mockAccRepo    *mock_domain.MockAccountRepository
	mockSpiderRepo *mock_domain.MockSpiderRepository
}

const (
	PNG_IMAGE_BASE64 = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAARElEQVR4nAA0AMv/BDxwkSMxNVIZ+8jBvwQJGDIaFSQnI/r17+0E9er629jy3A389/j6AwwMH/Xx+gP45djSxwEAAP//b58b81Ms3yEAAAAASUVORK5CYII="
)

var mockResultSpiderInfo = model.SpiderInfo{
	SpiderUUID:   "SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b",
	Family:       "Agelenidae",
	Genus:        "Draconarius",
	Species:      "abbreviatus",
	Author:       "Dankittipakul & Wang",
	PublishYear:  "2003",
	Country:      "Thailand",
	Altitude:     "1000, 1750 m",
	Method:       "pitfall trap, litter sample",
	Habital:      "pine forest, evergreen hill forest",
	Microhabital: "N/A",
	Status:       model.SPIDER_INFO_STATUS_ACTIVE,
	Designate:    "Short retrolateral apophysis",
	Address: []model.Address{
		{
			Province: "Chiang Mai",
			District: "Chomthong",
			Locality: "Doi Inthanon National Park",
			Position: []model.Position{
				{
					Name:      "Doi Inthaonon",
					Latitude:  18.58889676,
					Longitude: 98.48697532,
				},
			},
		},
	},
	Paper:     []string{"Test2023"},
	ImageFile: []string{"fwejiknfiow;ehnfiowhefn"},
}

// ======================================================================
// TestSpiderInfoUsecase_GetSpiderInfoUsecase
// ======================================================================
func TestSpiderInfoUsecase_GetSpiderInfoUsecase(t *testing.T) {

	type args struct {
		ctx        context.Context
		spiderUUID string
		username   string
	}
	tests := []struct {
		name       string
		args       args
		want       model.SpiderInfo
		buildStubs func(*mockStubsSpiderInfoUsecase)
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success_without_login",
			args: args{
				context.TODO(),
				"SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b",
				"",
			},
			buildStubs: success_without_login,
			want:       mockResultSpiderInfo,
			wantErr:    false,
		},
		{
			name: "success_with_login",
			args: args{
				context.TODO(),
				"SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b",
				"test",
			},
			buildStubs: success_with_login,
			want:       mockResultSpiderInfo,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			stubs := mockStubsSpiderInfoUsecase{
				mockAccRepo:    mock_domain.NewMockAccountRepository(ctrl),
				mockSpiderRepo: mock_domain.NewMockSpiderRepository(ctrl),
			}

			tt.buildStubs(&stubs)

			usecase := NewSpiderInfoUsecase(stubs.mockSpiderRepo, stubs.mockAccRepo, "")

			spiderInfoResult, err := usecase.GetSpiderInfoUsecase(context.TODO(), tt.args.spiderUUID, tt.args.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("[TestSpiderInfoUsecase_GetSpiderInfoUsecase] wantErr is `%v` but got `%v`, and error is `%+v`", tt.wantErr, (err != nil), err)
			}

			if !reflect.DeepEqual(*spiderInfoResult, tt.want) {
				t.Errorf("[TestSpiderInfoUsecase_GetSpiderInfoUsecase]\nwant: `%+v`\ngot: `%v`", tt.want, spiderInfoResult)
			}
		})
	}
}

func success_without_login(commonStubs *mockStubsSpiderInfoUsecase) {

	commonStubs.mockAccRepo.EXPECT().FindAccountByUsername(
		gomock.Any(),
		gomock.Eq(""),
	).Return(nil, repository.ErrorMongoNotFound)

	commonStubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b"),
	).Return(&mockResultSpiderInfo, nil)
}

func success_with_login(commonStubs *mockStubsSpiderInfoUsecase) {

	account := model.Account{
		Username: "test",
		Role:     "admin",
	}

	commonStubs.mockAccRepo.EXPECT().FindAccountByUsername(
		gomock.Any(),
		gomock.Eq("test"),
	).Return(&account, nil)

	commonStubs.mockSpiderRepo.EXPECT().FindSpiderByUUID(
		gomock.Any(),
		gomock.Eq("SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b"),
	).Return(&mockResultSpiderInfo, nil)
}

// **********************************************************************

// ======================================================================
// TestSpiderInfoUsecase_GetSpiderImageUsecase
// ======================================================================
func TestSpiderInfoUsecase_GetSpiderImageUsecase(t *testing.T) {
	tempDir := t.TempDir()

	type args struct {
		fileImages []string
	}

	tc := []struct {
		name         string
		args         args
		needTempFile bool
		want         []model.SpiderImageList
		wantErr      bool
	}{
		// can't not test jpeg case
		{
			name: "success_get_png_file",
			args: args{
				fileImages: []string{"SPIDER_dcb5dd72-d7c9-4b89-a039-abe670fcf3002023-03-04T19:44:22+0700-volume-0.png"},
			},
			needTempFile: true,
			want: []model.SpiderImageList{
				{
					Name:        "SPIDER_dcb5dd72-d7c9-4b89-a039-abe670fcf3002023-03-04T19:44:22+0700-volume-0.png",
					ImageBase64: PNG_IMAGE_BASE64,
				},
			},
			wantErr: false,
		},
		{
			name: "read_file_failed",
			args: args{
				fileImages: []string{"SPIDER_dcb5dd72-d7c9-4b89-a039-abe670fcf3002023-03-04T19:44:22+0700-volume-1.png"},
			},
			needTempFile: false,
			want:         []model.SpiderImageList{},
			wantErr:      true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.needTempFile {
				for _, filename := range tt.args.fileImages {
					tempImagePNG(tempDir, filename)
				}

			}

			mockSpiderRepo := mock_domain.NewMockSpiderRepository(ctrl)
			mockAccRepo := mock_domain.NewMockAccountRepository(ctrl)

			usecase := NewSpiderInfoUsecase(mockSpiderRepo, mockAccRepo, tempDir)

			spiderImageEncode, err := usecase.GetSpiderImagesUsecase(context.TODO(), tt.args.fileImages)
			if (err != nil) != tt.wantErr {
				t.Errorf("[TestSpiderInfoUsecase_GetSpiderImageUsecase] want error: %v, but got error: %v, and error: %v", tt.wantErr, (err != nil), err)
			}

			if !reflect.DeepEqual(tt.want, spiderImageEncode) && (len(tt.want) > 0 && len(spiderImageEncode) > 0) {
				t.Errorf("[TestSpiderInfoUsecase_GetSpiderImageUsecase] want %v, but got %v", tt.want, spiderImageEncode)
			}
		})
	}

}

func tempImagePNG(pathDir, name string) {

	filePath := path.Join(pathDir, name)

	commaIndex := strings.Index(string(PNG_IMAGE_BASE64), ",")
	OriginalImage := string(PNG_IMAGE_BASE64)[commaIndex+1:]

	imageDecode, _ := base64.StdEncoding.DecodeString(OriginalImage)

	renderImageDecode := bytes.NewBuffer(imageDecode)

	pngImage, _ := png.Decode(renderImageDecode)

	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	defer file.Close()

	png.Encode(file, pngImage)
}

// **********************************************************************

// ======================================================================
// TestSpiderInfoUsecase_GetSpiderImageUsecase
// ======================================================================
func TestSpiderInfoUsecase_GetSpiderInfoListManagerUsecase(t *testing.T) {
	config.LoadConfig("./../config", "config")

	type args struct {
		username string
		page     int
		size     int
	}
	testcase := []struct {
		name       string
		args       args
		buildStubs func(*mockStubsSpiderInfoUsecase)
		wantErr    bool
	}{
		{
			name: "success_get_spider_info_list_manager",
			args: args{
				username: "test",
				page:     0,
				size:     5,
			},
			buildStubs: success_get_spider_info_list_manager,
			wantErr:    false,
		},
	}

	for _, tt := range testcase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commonStubs := mockStubsSpiderInfoUsecase{
				mockAccRepo:    mock_domain.NewMockAccountRepository(ctrl),
				mockSpiderRepo: mock_domain.NewMockSpiderRepository(ctrl),
			}

			tt.buildStubs(&commonStubs)

			usecase := NewSpiderInfoUsecase(commonStubs.mockSpiderRepo, commonStubs.mockAccRepo, "")

			_, err := usecase.GetSpiderInfoListManager(context.TODO(), tt.args.username, tt.args.page, tt.args.size)

			if (err != nil) != tt.wantErr {
				t.Errorf("want error bool %v, but got error bool %v, error is %+v", tt.wantErr, (err != nil), err)
			}

		})
	}
}

func success_get_spider_info_list_manager(commonStubs *mockStubsSpiderInfoUsecase) {

	account := model.Account{
		Username: "test",
		Role:     "admin",
	}

	commonStubs.mockAccRepo.EXPECT().FindAccountByUsername(
		gomock.Any(),
		gomock.Eq("test"),
	).Return(&account, nil)

	spiderInfoList := []model.SpiderInfo{
		mockResultSpiderInfo,
	}

	commonStubs.mockSpiderRepo.EXPECT().FindAllSpiderListManager(
		gomock.Any(),
		gomock.Eq(0),
		gomock.Eq(5),
	).Return(spiderInfoList, nil)
}

// **********************************************************************

// ======================================================================
// TestSpiderInfoUsecase_GetSpiderInfoListByGeographies
// ======================================================================
func TestSpiderInfoUsecase_GetSpiderInfoListByGeographies(t *testing.T) {

	type args struct {
		province string
		district string
		position string
	}

	testcase := []struct {
		name       string
		args       args
		buildStubs func(*mockStubsSpiderInfoUsecase)
		wantErr    bool
	}{
		{
			name: "success_get_spider_info_by_province_and_district",
			args: args{
				province: "test",
				district: "test",
				position: "",
			},
			buildStubs: success_get_spider_info_by_province_and_district,
			wantErr:    false,
		},
		{
			name: "success_get_spider_info_by_position",
			args: args{
				province: "",
				district: "",
				position: "test",
			},
			buildStubs: success_get_spider_info_by_position,
			wantErr:    false,
		},
		{
			name: "fail_data_data_request",
			args: args{
				province: "",
				district: "",
				position: "",
			},
			buildStubs: func(commonStubs *mockStubsSpiderInfoUsecase) {},
			wantErr:    true,
		},
	}

	for _, tt := range testcase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commonStubs := mockStubsSpiderInfoUsecase{
				mockAccRepo:    mock_domain.NewMockAccountRepository(ctrl),
				mockSpiderRepo: mock_domain.NewMockSpiderRepository(ctrl),
			}
			tt.buildStubs(&commonStubs)

			usecase := NewSpiderInfoUsecase(commonStubs.mockSpiderRepo, commonStubs.mockAccRepo, "")

			_, err := usecase.GetSpiderInfoListByGeographies(context.TODO(), tt.args.province, tt.args.district, tt.args.position)

			if (err != nil) != tt.wantErr {
				t.Errorf("want error bool %v, but got error bool %v, error is %+v", tt.wantErr, (err != nil), err)
			}

		})
	}
}

func success_get_spider_info_by_province_and_district(commonStubs *mockStubsSpiderInfoUsecase) {
	spiderInfoList := []model.SpiderInfo{mockResultSpiderInfo}
	commonStubs.mockSpiderRepo.EXPECT().FindSpiderInfoListByGeographies(
		gomock.Any(),
		gomock.Eq("test"),
		gomock.Eq("test"),
		gomock.Eq(""),
	).Return(spiderInfoList, nil)
}

func success_get_spider_info_by_position(commonStubs *mockStubsSpiderInfoUsecase) {
	spiderInfoList := []model.SpiderInfo{mockResultSpiderInfo}

	commonStubs.mockSpiderRepo.EXPECT().FindSpiderInfoListByGeographies(
		gomock.Any(),
		gomock.Eq(""),
		gomock.Eq(""),
		gomock.Eq("test"),
	).Return(spiderInfoList, nil)
}

func TestSpiderInfoUsecase_GetSpiderInfoListByLocality(t *testing.T) {

	type args struct {
		locality string
		page     int32
		size     int32
	}
	tests := []struct {
		name    string
		args    args
		stubs   func(*mockStubsSpiderInfoUsecase)
		want    []model.SpiderInfo
		wantErr bool
	}{
		{
			name: "success_locality",
			args: args{
				locality: "Doi Inthaonon",
				page:     0,
				size:     10,
			},
			stubs: success_locality,
			want: []model.SpiderInfo{
				mockResultSpiderInfo,
			},
			wantErr: false,
		},
		{
			name: "fail_spider_info_not_found",
			args: args{
				locality: "Doi Inthanon 02",
				page:     0,
				size:     10,
			},
			stubs:   fail_spider_info_not_found,
			want:    []model.SpiderInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commonBuildStub := mockStubsSpiderInfoUsecase{
				mockAccRepo:    mock_domain.NewMockAccountRepository(ctrl),
				mockSpiderRepo: mock_domain.NewMockSpiderRepository(ctrl),
			}

			tt.stubs(&commonBuildStub)

			usecase := NewSpiderInfoUsecase(commonBuildStub.mockSpiderRepo, commonBuildStub.mockAccRepo, "")

			got, err := usecase.GetSpiderInfoListByLocality(context.TODO(), tt.args.locality, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpiderInfoUsecase.GetSpiderInfoListByLocality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpiderInfoUsecase.GetSpiderInfoListByLocality() = %v, want %v", got, tt.want)
			}
		})
	}
}

func success_locality(stubs *mockStubsSpiderInfoUsecase) {
	stubs.mockSpiderRepo.EXPECT().FindSpiderInfoByLocality(
		gomock.Any(),
		gomock.Eq("Doi Inthaonon"),
		gomock.Eq(int32(0)),
		gomock.Eq(int32(10)),
	).Return([]model.SpiderInfo{mockResultSpiderInfo}, nil)
}

func fail_spider_info_not_found(stubs *mockStubsSpiderInfoUsecase) {
	stubs.mockSpiderRepo.EXPECT().FindSpiderInfoByLocality(
		gomock.Any(),
		gomock.Eq("Doi Inthanon 02"),
		gomock.Eq(int32(0)),
		gomock.Eq(int32(10)),
	).Return([]model.SpiderInfo{}, repository.ErrorMongoNotFound)
}
