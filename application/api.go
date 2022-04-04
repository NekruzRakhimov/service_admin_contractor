package application

import (
	"context"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"net/http"
	"service_admin_contractor/application/cerrors"
	"service_admin_contractor/application/config"
	"service_admin_contractor/application/controller"
	"service_admin_contractor/application/cvalidator"
	"service_admin_contractor/application/middleware"
	"service_admin_contractor/application/respond"
	"service_admin_contractor/application/service"
	"service_admin_contractor/infrastructure/logging"
	"service_admin_contractor/infrastructure/persistence/postgres"
)

// NewApi конфигурирует API
func NewApi() (http.Handler, error) {
	logging.ConfigureLogger()
	cvalidator.ConfigureValidator()

	pc := postgres.DBConn()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handleNotFoundError)

	configureHealthchecks(r, pc)

	r.Use(middleware.CorrelationHandler)
	r.Use(middleware.RequestLoggerHandler)
	r.Use(middleware.RecoveryHandler(middleware.PrintRecoveryStack(true)))

	r.Use(middleware.CORS(
		middleware.AllowedOrigins(viper.GetStringSlice(config.CorsAllowedOrigins)),
		middleware.AllowedMethods(viper.GetStringSlice(config.CorsAllowedMethods)),
		middleware.AllowedHeaders(viper.GetStringSlice(config.CorsAllowedHeaders)),
		middleware.AllowCredentials()))

	err := configureRoutes(r, pc)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func configureHealthchecks(r *mux.Router, pc *pgxpool.Pool) {
	r.Handle("/health", healthcheck.Handler(
		healthcheck.WithTimeout(viper.GetDuration(config.HealthcheckTimeout)),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(
			func(ctx context.Context) error {
				return pc.Ping(ctx)
			}),
		),
	))
}

func configureRoutes(r *mux.Router, pc *pgxpool.Pool) error {
	contractorRepo := postgres.NewContractorRepository(pc)
	bpmsUserRepo := postgres.NewBpmsUserRepository(pc)
	contractorSrvc := service.NewContractorService(contractorRepo)

	//region Contractor routes
	api := r.PathPrefix("/api/v1/admin").Subrouter()
	api.Use(middleware.AuthHandler(bpmsUserRepo))

	controller.NewContractorController(contractorSrvc).HandleRoutes(api)
	//endregion

	return nil
}

func handleNotFoundError(w http.ResponseWriter, r *http.Request) {
	respond.WithError(w, r, cerrors.ErrResourceNotFound(r))
}
