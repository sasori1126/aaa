package middlewares

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/entities"
	"encoding/json"
	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
	"net/http"
	"time"
)

func Timeout() gin.HandlerFunc {
	timeoutMessage := &entities.ApiError{
		Message: "request timeout, try again",
		Code:    http.StatusRequestTimeout,
		Errors: entities.E{
			"request": "Request timeout",
		},
	}

	msg, _ := json.Marshal(timeoutMessage)

	return timeout.Timeout(
		timeout.WithTimeout(time.Second*configs.Timeout),
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),
		timeout.WithDefaultMsg(string(msg)),
		timeout.WithCallBack(func(r *http.Request) {
			configs.Logger.Info("Request failed at: " + r.URL.String() + " method: " + r.Method)
		}),
	)
}
