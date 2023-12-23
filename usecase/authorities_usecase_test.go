package usecase

// import (
// 	"context"
// 	"reflect"
// 	mock_domain "spider-go/domain/mock"
// 	"spider-go/model"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func TestAuthorities_Login(t *testing.T) {

// 	type args struct {
// 		username string
// 		password string
// 	}
// 	tests := []struct {
// 		name            string
// 		args            args
// 		buildStubs      func(mock_domain.MockAccountRepository, mock_domain.MockRedisRepository)
// 		wantAccountInfo *model.Account
// 		wantErr         error
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "success",
// 			args: args{
// 				username: "unittest",
// 				password: "unittestsuccess",
// 			},
// 			buildStubs: successLogin,
// 			wantAccountInfo: &model.Account{
// 				Username:     "unittest",
// 				HashPassword: "$2a$10$jfEYdAp0ZjH4LSyBl1NMYuH97s4sbuEgiv1kRBxNqzBsBnuenRD8i",
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "invalid password",
// 			args: args{
// 				username: "unittest",
// 				password: "nopassword",
// 			},
// 			buildStubs:      invalidPassword,
// 			wantAccountInfo: nil,
// 			wantErr:         ErrorAuthoritiesInvalidPassword,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockAccountRepository := mock_domain.NewMockAccountRepository(ctrl)
// 			mockRedisRepository := mock_domain.NewMockRedisRepository(ctrl)

// 			tt.buildStubs(*mockAccountRepository, *mockRedisRepository)

// 			usecase := NewAuthoritiesUsecase(mockAccountRepository, mockRedisRepository)
// 			gotAccountInfo, _, err := usecase.Login(context.TODO(), tt.args.username, tt.args.password)
// 			if err != tt.wantErr {
// 				t.Errorf("Authorities.Login() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotAccountInfo, tt.wantAccountInfo) {
// 				t.Errorf("Authorities.Login() = %v, want %v", gotAccountInfo, tt.wantAccountInfo)
// 			}
// 		})
// 	}
// }

// func successLogin(mockAccountRepo mock_domain.MockAccountRepository, mockRedisRepo mock_domain.MockRedisRepository) {
// 	mockAccountRepo.EXPECT().FindAccountByUsername(
// 		gomock.Any(),
// 		gomock.Eq("unittest"),
// 	).Return(
// 		&model.Account{
// 			Username:     "unittest",
// 			HashPassword: "$2a$10$jfEYdAp0ZjH4LSyBl1NMYuH97s4sbuEgiv1kRBxNqzBsBnuenRD8i",
// 		},
// 		nil,
// 	)

// 	mockRedisRepo.EXPECT().SetDataToRedisWithTTL(
// 		gomock.Any(),
// 		gomock.Any(),
// 		gomock.Eq(model.Login{Username: "unittest"}),
// 		gomock.Any(),
// 	).Return(nil)
// }

// func invalidPassword(mockAccountRepo mock_domain.MockAccountRepository, mockRedisRepo mock_domain.MockRedisRepository) {
// 	mockAccountRepo.EXPECT().FindAccountByUsername(
// 		gomock.Any(),
// 		gomock.Eq("unittest"),
// 	).Return(
// 		&model.Account{
// 			Username:     "unittest",
// 			HashPassword: "$2a$10$jfEYdAp0ZjH4LSyBl1NMYuH97s4sbuEgiv1kRBxNqzBsBnuenRD8i",
// 		},
// 		nil,
// 	)
// }
