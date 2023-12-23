package usecase

import (
	"context"
	"fmt"
	api_model "spider-go/api/model"
	"spider-go/config"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

var mockSpiderInfo = api_model.SpiderInfo{
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

func TestRegisterSpiderUsecase_Register(t *testing.T) {

	config.LoadConfig("./../config", "config")

	type args struct {
		ctx      context.Context
		req      api_model.SpiderInfo
		username string
	}
	tests := []struct {
		name       string
		args       args
		buildStubs func(mock_domain.MockSpiderRepository, mock_domain.MockStatisticsRepository)
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success with insert new spider statistic case",
			args: args{
				ctx:      context.TODO(),
				req:      mockSpiderInfo,
				username: "testSuccess",
			},
			buildStubs: successRegisterWithInsertNewStatistic,
			wantErr:    false,
		},
		{
			name: "success with update spider statistic case",
			args: args{
				ctx:      context.TODO(),
				req:      mockSpiderInfo,
				username: "testSuccess",
			},
			buildStubs: successRegisterWithUpdateStatistic,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSpiderRepo := mock_domain.NewMockSpiderRepository(ctrl)
			mockStatisticsRepo := mock_domain.NewMockStatisticsRepository(ctrl)

			tt.buildStubs(*mockSpiderRepo, *mockStatisticsRepo)

			usecase := NewRegisterSpiderUsecase(mockSpiderRepo, mockStatisticsRepo)

			_, err := usecase.Register(tt.args.ctx, tt.args.req, tt.args.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register fail: wantErr %v, but got error %v, error is %v", tt.wantErr, (err != nil), err)
			}
		})
	}
}

func successRegisterWithInsertNewStatistic(
	SpiderRepo mock_domain.MockSpiderRepository,
	statisticsRepo mock_domain.MockStatisticsRepository,
) {

	tn := time.Now()

	// find spider statistic with family in step of validating spider statistic
	statisticsRepo.EXPECT().FindSpiderStatisticsByFamily(
		gomock.Any(),
		gomock.Eq("Agelenidae"),
	).Return(&model.SpiderStatistics{}, fmt.Errorf("mongo not found"))

	statisticsRepo.EXPECT().UpsertSpiderStatistics(
		gomock.Any(),
		gomock.Eq("Agelenidae"),
		EqSpiderStatistic(
			model.SpiderStatistics{
				FamilyName: "Agelenidae",
				Genus: []model.GenusGroup{
					{
						GenusName: "Draconarius",
						Species: []model.SpeciesGroup{
							{
								SpeciesName: "abbreviatus",
							},
						},
					},
				},
				CreatedAt: tn,
				UpdatedAt: tn,
			},
		),
	).Return(nil)

	SpiderRepo.EXPECT().InsertNewSpider(
		gomock.Any(),
		EqSpiderInfo(
			model.SpiderInfo{
				Family:       mockSpiderInfo.Family,
				Genus:        mockSpiderInfo.Genus,
				Species:      mockSpiderInfo.Species,
				Author:       mockSpiderInfo.Author,
				PublishYear:  mockSpiderInfo.PublishYear,
				Country:      mockSpiderInfo.Country,
				CountryOther: mockSpiderInfo.OtherCountries,
				Altitude:     mockSpiderInfo.Altitude,
				Method:       mockSpiderInfo.Method,
				Habital:      mockSpiderInfo.Habital,
				Microhabital: mockSpiderInfo.Microhabital,
				Designate:    mockSpiderInfo.Designate,
				Status:       model.SPIDER_INFO_STATUS_ACTIVE,
				Address: []model.Address{
					{
						Province: mockSpiderInfo.Address[0].Province,
						District: mockSpiderInfo.Address[0].District,
						Locality: mockSpiderInfo.Address[0].Locality,
						Position: []model.Position{
							{
								Name:      mockSpiderInfo.Address[0].Position[0].Name,
								Latitude:  mockSpiderInfo.Address[0].Position[0].Latitude,
								Longitude: mockSpiderInfo.Address[0].Position[0].Longitude,
							},
						},
					},
				},
				Paper:     mockSpiderInfo.Paper,
				CreatedBy: "testSuccess",
			},
		),
	).Return(nil)
}

func successRegisterWithUpdateStatistic(
	SpiderRepo mock_domain.MockSpiderRepository,
	statisticsRepo mock_domain.MockStatisticsRepository,
) {

	tn := time.Now()

	// find spider statistic with family in step of validating spider statistic
	prepareSpiderStatistics := model.SpiderStatistics{
		FamilyName: "Agelenidae",
		Genus: []model.GenusGroup{
			{
				GenusName: "Draconarius",
			},
		},
		CreatedAt: tn,
		UpdatedAt: tn,
	}

	statisticsRepo.EXPECT().FindSpiderStatisticsByFamily(
		gomock.Any(),
		gomock.Eq("Agelenidae"),
	).Return(&prepareSpiderStatistics, nil)

	// upsert spider statistic with family in step of validating spider statistic
	statisticsRepo.EXPECT().UpsertSpiderStatistics(
		gomock.Any(),
		gomock.Eq("Agelenidae"),
		EqSpiderStatistic(
			model.SpiderStatistics{
				FamilyName: "Agelenidae",
				Genus: []model.GenusGroup{
					{
						GenusName: "Draconarius",
						Species: []model.SpeciesGroup{
							{
								SpeciesName: "abbreviatus",
							},
						},
					},
				},
				CreatedAt: tn,
				UpdatedAt: tn,
			},
		),
	).Return(nil)

	SpiderRepo.EXPECT().InsertNewSpider(
		gomock.Any(),
		EqSpiderInfo(
			model.SpiderInfo{
				Family:       mockSpiderInfo.Family,
				Genus:        mockSpiderInfo.Genus,
				Species:      mockSpiderInfo.Species,
				Author:       mockSpiderInfo.Author,
				PublishYear:  mockSpiderInfo.PublishYear,
				Country:      mockSpiderInfo.Country,
				CountryOther: mockSpiderInfo.OtherCountries,
				Altitude:     mockSpiderInfo.Altitude,
				Method:       mockSpiderInfo.Method,
				Habital:      mockSpiderInfo.Habital,
				Microhabital: mockSpiderInfo.Microhabital,
				Designate:    mockSpiderInfo.Designate,
				Status:       model.SPIDER_INFO_STATUS_ACTIVE,
				Address: []model.Address{
					{
						Province: mockSpiderInfo.Address[0].Province,
						District: mockSpiderInfo.Address[0].District,
						Locality: mockSpiderInfo.Address[0].Locality,
						Position: []model.Position{
							{
								Name:      mockSpiderInfo.Address[0].Position[0].Name,
								Latitude:  mockSpiderInfo.Address[0].Position[0].Latitude,
								Longitude: mockSpiderInfo.Address[0].Position[0].Longitude,
							},
						},
					},
				},
				Paper:     mockSpiderInfo.Paper,
				CreatedBy: "testSuccess",
			},
		),
	).Return(nil)
}
