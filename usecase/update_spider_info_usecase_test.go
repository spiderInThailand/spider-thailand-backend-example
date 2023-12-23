package usecase

import (
	"context"
	api_model "spider-go/api/model"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"testing"

	"github.com/golang/mock/gomock"
)

type commonStubsUpdateSpider struct {
	mockSpiderRepo *mock_domain.MockSpiderRepository
}

var mockDataSpiderInfo = api_model.SpiderInfo{
	SpiderUUID:   "SPIDER_8a5bbf23-8ccd-4068-ae19-145095e0847b",
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
	Designate:    "Short retrolateral apophysis",
	Address: []api_model.Address{
		{
			Province: "Chiang Mai",
			District: "Chomthong",
			Locality: "Doi Inthanon National Park",
			Position: []api_model.Position{
				{
					Name:      "Doi Inthaonon",
					Latitude:  18.58889676,
					Longitude: 98.48697532,
				},
			},
		},
	},
	Paper: []string{"Test2023"},
}

func TestUpdateSpiderInfoUsecase(t *testing.T) {

	type args struct {
		spiderInfoReq api_model.SpiderInfo
	}

	tests := []struct {
		name    string
		args    args
		stubs   func(*commonStubsUpdateSpider)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "update_spider_info_success_case",
			args: args{
				spiderInfoReq: mockDataSpiderInfo,
			},
			stubs:   update_spider_info_success_case,
			wantErr: false,
		},
		{
			name: "update_spider_info_not_found_case",
			args: args{
				spiderInfoReq: mockDataSpiderInfo,
			},
			stubs:   update_spider_info_not_found_case,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commonStubs := commonStubsUpdateSpider{
				mockSpiderRepo: mock_domain.NewMockSpiderRepository(ctrl),
			}

			tt.stubs(&commonStubs)

			usecase := NewUpdateSpiderInfoUsecase(commonStubs.mockSpiderRepo)

			err := usecase.UpdateSpiderInfoUsecase(context.TODO(), tt.args.spiderInfoReq)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("[TestUpdateSpiderInfoUsecase] fail wantErr is %v, but gotErr is %v, error: %+v", tt.wantErr, gotErr, err)
			}

		})
	}
}

func update_spider_info_success_case(mockStubs *commonStubsUpdateSpider) {

	spiderInfo := model.SpiderInfo{
		SpiderUUID:   "SPIDER_8a5bbf23-8ccd-4068-ae19-145095e0847b",
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
		Paper: []string{"Test2023"},
	}

	mockStubs.mockSpiderRepo.EXPECT().UpdateSpiderInfo(
		gomock.Any(),
		gomock.Eq("SPIDER_8a5bbf23-8ccd-4068-ae19-145095e0847b"),
		EqSpiderInfo(spiderInfo),
	).Return(true, nil)
}

func update_spider_info_not_found_case(mockStubs *commonStubsUpdateSpider) {

	spiderInfo := model.SpiderInfo{
		SpiderUUID:   "SPIDER_8a5bbf23-8ccd-4068-ae19-145095e0847b",
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
		Paper: []string{"Test2023"},
	}

	mockStubs.mockSpiderRepo.EXPECT().UpdateSpiderInfo(
		gomock.Any(),
		gomock.Eq("SPIDER_8a5bbf23-8ccd-4068-ae19-145095e0847b"),
		EqSpiderInfo(spiderInfo),
	).Return(false, nil)
}
