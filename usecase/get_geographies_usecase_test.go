package usecase

import (
	"context"
	"reflect"
	mock_domain "spider-go/domain/mock"
	"spider-go/model"
	"spider-go/repository"
	"testing"

	"github.com/golang/mock/gomock"
)

type commonStubsGetGeographiesBySpiderType struct {
	mockSpiderRepo          *mock_domain.MockSpiderRepository
	mockThaiGeographiesRepo *mock_domain.MockThaiGeographiesRepository
}

var mockResultSpiderInfoList = []model.SpiderInfo{
	{
		SpiderUUID:   "SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b",
		Family:       "family_1",
		Genus:        "genus_1",
		Species:      "species_1",
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
			{
				Province: "Chiang Mai",
				District: "Muang ",
				Locality: "Doi Suthep-Pui National Park",
				Position: []model.Position{
					{
						Name:      "Doi Pui",
						Latitude:  18.81665591,
						Longitude: 98.88327951,
					},
				},
			},
		},
		Paper:     []string{"Test2023"},
		ImageFile: []string{"fwejiknfiow;ehnfiowhefn"},
	},
	{
		SpiderUUID:   "SPIDER_c6ef5023-94fc-41c8-a88d-87303c75999b",
		Family:       "family_1",
		Genus:        "genus_1",
		Species:      "species_2",
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
						Name:      "Doi Inthaonon 02",
						Latitude:  18.58889676,
						Longitude: 98.48697532,
					},
				},
			},
		},
		Paper:     []string{"Test2023"},
		ImageFile: []string{"fwejiknfiow;ehnfiowhefn"},
	},
}

func TestThaiGeographiesUsecase_GetGeographiesBySpiderType(t *testing.T) {

	type args struct {
		family  string
		genus   string
		species string
	}
	tests := []struct {
		name    string
		args    args
		stubs   func(*commonStubsGetGeographiesBySpiderType)
		want    []model.LocationResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success_get_family_1",
			args: args{
				family:  "family_1",
				genus:   "",
				species: "",
			},
			stubs: success_get_family_1,
			want: []model.LocationResult{
				{
					Province: "Chiang Mai",
					Locality: []model.LocalityResult{
						{
							Name: "Doi Inthanon National Park",
							SubLocaltion: []string{
								"Doi Inthaonon",
								"Doi Inthaonon 02",
							},
						},
						{
							Name: "Doi Suthep-Pui National Park",
							SubLocaltion: []string{
								"Doi Pui",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success_get_family_1",
			args: args{
				family:  "family_1",
				genus:   "genus_1",
				species: "species_2",
			},
			stubs: success_get_family_1_genus_1_species_2,
			want: []model.LocationResult{
				{
					Province: "Chiang Mai",
					Locality: []model.LocalityResult{
						{
							Name: "Doi Inthanon National Park",
							SubLocaltion: []string{
								"Doi Inthaonon 02",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail_mongo_not_found",
			args: args{
				family:  "family_1",
				genus:   "genus_1",
				species: "species_3",
			},
			stubs:   fail_mongo_not_found,
			want:    []model.LocationResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commonStubs := commonStubsGetGeographiesBySpiderType{
				mockSpiderRepo:          mock_domain.NewMockSpiderRepository(ctrl),
				mockThaiGeographiesRepo: mock_domain.NewMockThaiGeographiesRepository(ctrl),
			}

			tt.stubs(&commonStubs)

			usecase := NewThaiGeographiesUsecase(
				commonStubs.mockThaiGeographiesRepo,
				commonStubs.mockSpiderRepo,
			)

			got, err := usecase.GetGeographiesBySpiderType(context.TODO(), tt.args.family, tt.args.genus, tt.args.species)
			if (err != nil) != tt.wantErr {
				t.Errorf("ThaiGeographiesUsecase.GetGeographiesBySpiderType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThaiGeographiesUsecase.GetGeographiesBySpiderType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func success_get_family_1(stub *commonStubsGetGeographiesBySpiderType) {

	familyOneSpiderInfo := []model.SpiderInfo{
		mockResultSpiderInfoList[0],
		mockResultSpiderInfoList[1],
	}

	stub.mockSpiderRepo.EXPECT().FindSpiderInfoBySpiderType(
		gomock.Any(),
		gomock.Eq("family_1"),
		gomock.Eq(""),
		gomock.Eq(""),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(familyOneSpiderInfo, nil)
}

func success_get_family_1_genus_1_species_2(stub *commonStubsGetGeographiesBySpiderType) {

	familyOneSpiderInfo := []model.SpiderInfo{
		mockResultSpiderInfoList[1],
	}

	stub.mockSpiderRepo.EXPECT().FindSpiderInfoBySpiderType(
		gomock.Any(),
		gomock.Eq("family_1"),
		gomock.Eq("genus_1"),
		gomock.Eq("species_2"),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(familyOneSpiderInfo, nil)
}

func fail_mongo_not_found(stub *commonStubsGetGeographiesBySpiderType) {

	stub.mockSpiderRepo.EXPECT().FindSpiderInfoBySpiderType(
		gomock.Any(),
		gomock.Eq("family_1"),
		gomock.Eq("genus_1"),
		gomock.Eq("species_3"),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return([]model.SpiderInfo{}, repository.ErrorMongoNotFound)
}
