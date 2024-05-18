package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/petshop-system/petshop-api-gateway/configuration/db"
	"github.com/petshop-system/petshop-api-gateway/configuration/environment"
	"github.com/petshop-system/petshop-api-gateway/server"
	database "github.com/petshop-system/petshop-api-gateway/server/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var loggerSugar *zap.SugaredLogger

func init() {

	err := envconfig.Process("setting", &environment.Setting)
	if err != nil {
		panic(err.Error())
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any
	loggerSugar = logger.Sugar()

}

func main() {

	postgresConnectionDB := db.NewPostgresDB(environment.Setting.Postgres.DBUser, environment.Setting.Postgres.DBPassword,
		environment.Setting.Postgres.DBName, environment.Setting.Postgres.DBHost, environment.Setting.Postgres.DBPort, loggerSugar)

	gatewayDB := database.NewGatewayDB(database.WithDB(postgresConnectionDB),
		database.WithLogger(loggerSugar))

	serveReverseProxyPass := server.NewServerPass(loggerSugar, gatewayDB, environment.Setting.Application.TickerReloadRouters)

	contextPath := environment.Setting.Server.Context

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "json/application")
		w.WriteHeader(http.StatusOK)
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(map[string]string{
			"status": "OK",
		})
		w.Write(body.Bytes())
	})

	http.Handle("/", &serveReverseProxyPass)

	loggerSugar.Infow("server started", "port", environment.Setting.Server.Port,
		"contextPath", contextPath)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", environment.Setting.Server.Port), nil); err != nil {
		loggerSugar.Fatalw("could not start running server", "err", err)
	}
}
