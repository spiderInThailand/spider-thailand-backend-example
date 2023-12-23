package route

import (
	"net/http"
	"spider-go/api/handler"
	"spider-go/api/middleware"
	"spider-go/config"
	"spider-go/database"
	"spider-go/logger"
	"spider-go/repository"
	"spider-go/usecase"
	jwt_service "spider-go/utils/jwt"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(log *logger.Logger, conf *config.Root) *gin.Engine {

	// ==========================================================
	// create service
	// ==========================================================

	jwtService := jwt_service.NewJWTService(config.C().JWT.Secret, config.C().JWT.ExpireTime, conf.JWT.Issure)

	// ==========================================================
	// create repository
	// ==========================================================

	accountRepo := repository.NewAccountRepository(database.DB)
	spiderStatisticsRepo := repository.NewSpiderStatisticsRepository(database.DB)
	spiderRepo := repository.NewSpiderRepository(database.DB)
	thaiGeographiesRepo := repository.NewThaiGeographiesRepository(database.DB)

	// ==========================================================
	// create usecase
	// ==========================================================

	authoritailUsecase := usecase.NewAuthoritiesUsecase(accountRepo, jwtService)
	spiderStatisticsUsecase := usecase.NewSpiderStatisticsUsecase(spiderStatisticsRepo)
	registerSpiderUsercase := usecase.NewRegisterSpiderUsecase(spiderRepo, spiderStatisticsRepo)
	uploadImageusecase := usecase.NewUploadImageUsecase(spiderRepo)
	spiderInfoUsecase := usecase.NewSpiderInfoUsecase(spiderRepo, accountRepo, conf.File.FileImagePath)
	deleteSpiderInfoUsecase := usecase.NewDeleteSpiderInfoUsecase(spiderRepo)
	updateSpiderInfoUsecase := usecase.NewUpdateSpiderInfoUsecase(spiderRepo)
	removeSpiderImageUsecase := usecase.NewRemoveSpiderImageUsecase(spiderRepo, conf.File.FileImagePath)
	thaiGeographiesUsecase := usecase.NewThaiGeographiesUsecase(thaiGeographiesRepo, spiderRepo)
	getFamilyListUsecase := usecase.NewGetFamilyListUsecase(spiderStatisticsRepo, spiderRepo)

	// ==========================================================
	// create handler
	// ==========================================================

	createAccoutHandler := handler.NewCreateAccountHandler(authoritailUsecase)
	loginHandler := handler.NewLoginHandler(authoritailUsecase)
	registerHandler := handler.NewRegisterHandler(registerSpiderUsercase)
	spiderStatisticsHandler := handler.NewGetSpiderStatisricsHandler(spiderStatisticsUsecase, getFamilyListUsecase)
	spiderSettingHandler := handler.NewSpiderSettingHandler(uploadImageusecase, deleteSpiderInfoUsecase, updateSpiderInfoUsecase, removeSpiderImageUsecase)
	spiderInfoHandler := handler.NewSpiderInfoHandler(spiderInfoUsecase, thaiGeographiesUsecase)
	getGeographiesHandler := handler.NewGetGeographinesHandler(thaiGeographiesUsecase)

	// ==========================================================
	// create gin web service
	// ==========================================================

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())

	r.Use(limits.RequestSizeLimiter(conf.API.MaxRequestSize))
	r.Use(
		gin.Recovery(),
		// add middleware in the future
	)

	// ==========================================================
	// api for test service
	// ==========================================================

	r.GET("/test-service", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "server is running.",
		})
	})

	// ==========================================================
	// group 1: no login required
	// ==========================================================

	g1 := r.Group("")
	{
		// post
		g1.POST("", createAccoutHandler.CreateAccout)
		g1.POST("", loginHandler.Login)
		g1.POST("", spiderInfoHandler.GetOneSpiderInfoHandler)
		g1.POST("", spiderInfoHandler.GetSpiderImagesHandler)
		g1.POST("", getGeographiesHandler.GetProvinceHandler)
		g1.POST("", getGeographiesHandler.GetDistrictByProvinceNameHandler)
		g1.POST("", spiderInfoHandler.GetSpiderInfoByGeographiesHandler)
		g1.POST("", spiderInfoHandler.GetGeographiesBySpiderTypeHandler)
		g1.POST("", spiderInfoHandler.GetSpiderInfoByLocalityHandler)
		g1.POST("", spiderInfoHandler.GetSpiderListBySpiderTypeHandler)
		g1.POST("", spiderStatisticsHandler.GetSpiderStatisticsListHandler)
		g1.POST("", spiderStatisticsHandler.GetFamilyListhandler)

	}
	// **********************************************************

	// ==========================================================
	// group 2: login required
	// ==========================================================

	g2 := r.Group("")
	g2.Use(middleware.Authenticate(jwtService))
	{
		g2.POST("", loginHandler.VerifyLogin)
		g2.POST("", registerHandler.RegisterHandler)
		g2.POST("", spiderSettingHandler.UploadImageSpiderHandler)
		g2.POST("", spiderInfoHandler.GetOneSpiderInfoHandler)
		g2.POST("", spiderInfoHandler.GetSpiderInfoListManagerHandler)
		g2.POST("", spiderSettingHandler.DeleteSpiderHandler)
		g2.POST("", spiderSettingHandler.EditSpiderInfoHandler)
		g2.POST("", spiderSettingHandler.RemoveSpiderImageHandler)
	}
	// **********************************************************

	return r
}
