package route

import (
	"csm-api/config"
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SiteRoute(safeDB *sqlx.DB, timeSheetDB *sqlx.DB, r *store.Repository, apiConfig *config.ApiConfig) chi.Router {
	router := chi.NewRouter()

	siteHandler := &handler.HandlerSite{
		Service: &service.ServiceSite{
			SafeDB:            safeDB,
			SafeTDB:           safeDB,
			Store:             r,
			ProjectStore:      r,
			ProjectDailyStore: r,
			SitePosStore:      r,
			SiteDateStore:     r,
			UserService: &service.ServiceUser{
				SafeDB:      safeDB,
				TimeSheetDB: timeSheetDB,
				Store:       r,
			},
			ProjectService: &service.ServiceProject{
				SafeDB:    safeDB,
				Store:     r,
				UserStore: r,
			},
			WeatherApiService: &service.ServiceWeather{
				ApiKey:       apiConfig,
				SafeDB:       safeDB,
				SafeTDB:      safeDB,
				Store:        r,
				SitePosStore: r,
			},
			AddressSearchAPIService: &service.ServiceAddressSearch{
				ApiKey: apiConfig,
			},
		},
		CodeService: &service.ServiceCode{
			SafeDB: safeDB,
			Store:  r,
		},
	}

	router.Get("/", siteHandler.List)                // 현장관리 조회
	router.Get("/nm", siteHandler.SiteNameList)      // 현장명 조회
	router.Get("/stats", siteHandler.StatsList)      // 현장상태조회
	router.Post("/", siteHandler.Add)                // 현장 생성
	router.Put("/", siteHandler.Modify)              // 수정
	router.Put("/non-use", siteHandler.ModifyNonUse) // 현장 사용안함

	return router
}
