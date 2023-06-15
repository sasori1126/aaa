package main

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/controllers"
	"axis/ecommerce-backend/pkg/logger"
	"axis/ecommerce-backend/pkg/routes"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "axis/ecommerce-backend/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

func main() {
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	//shutdown, err := otel.InitOtel()
	//if err != nil {
	//	bugsnag.Notify(err)
	//	log.Fatalf("failed to start tracing %v", err)
	//}
	//defer func() {
	//	err := shutdown(context.Background())
	//	if err != nil {
	//		bugsnag.Notify(err)
	//		log.Fatalf("failed to start tracing %v", err)
	//	}
	//}()

	sl := logger.InitLogger()
	defer func(sLogger *zap.SugaredLogger) {
		err := sLogger.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(sl)

	configs.Logger = sl
	err := configs.InitConfig()
	if err != nil {
		sl.Panic(err)
	}

	sl.Infow("Starting application....")
	serve, err := controllers.InitServer()
	if err != nil {
		log.Fatalln("could not init serve", err)
	}

	router := routes.StartApp(serve)
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Failed to start server.")
		}
	}()

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		fmt.Println(sig)
		done <- true
	}()
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("Forced to shutdown")
	}
}
