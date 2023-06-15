package configs

import (
	"errors"
	"time"
)

var ErrOrderFound = errors.New("record found")

const (
	HashIdMinLength int    = 42
	Locale          string = "en"
	// BadRequest 400 Bad Request – This means that client-side input fails validation.
	BadRequest int = 400
	// Unauthorized 401 Unauthorized – This means the user isn’t not authorized to access a resource. It usually returns when the user isn’t authenticated.
	Unauthorized int = 401
	// Forbidden 403 Forbidden – This means the user is authenticated, but it’s not allowed to access a resource.
	Forbidden int = 403
	// NotFound 404 Not Found – This indicates that a resource is not found.
	NotFound int = 404
	// ServerError 500 Internal server error – This is a generic server error. It probably shouldn’t be thrown explicitly.
	ServerError int = 500
	// ServiceUnavailable 503 Service Unavailable – This indicates that something unexpected happened on server side (It can be anything like server overload, some parts of the system failed, etc.).
	ServiceUnavailable      int           = 503
	Ok                      int           = 200
	Created                 int           = 201
	NoContent               int           = 204
	ManyRequest             int           = 429
	Conflict                int           = 409
	Timeout                 time.Duration = 5
	User                    string        = "user"
	AxisEmail               string        = "wayne@axisforestry.com"
	AxisPartRyanEmail       string        = "ryan@biglist.ca"
	AxisPartSaleEmail       string        = "parts@axisforestry.com"
	PaymentMethodCard                     = "card"
	PaymentMethodPhoneOrder               = "phoneOrder"
	PaymentMethodOnAccount                = "onAccount"
	DeliveryMethodPickUp                  = "pickUp"
	DeliveryMethodDeliver                 = "deliver"
	OrderPending                          = "pending"
	OrderPendingOnPhone                   = "pending phone payment"
	OrderPendingOnAccount                 = "pending on account payment"
	OrderPaid                             = "paid"
	OrderCancelled                        = "cancelled"
	OrderFailed                           = "failed"
	OrderShipped                          = "shipped"
	OrderDelivered                        = "delivered"
	PaymentSuccessful                     = "success"
	PaymentFailed                         = "failed"
)

var ErrOrderAlreadyExist = ErrOrderFound
