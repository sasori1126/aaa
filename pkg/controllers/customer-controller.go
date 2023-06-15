package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/customers"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/notification/mail"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) TempCustomerRequest(c *gin.Context) {
	cr := &dto.TempCustomerRequest{}
	err := c.ShouldBindJSON(cr)
	if err != nil {
		ers, ok := err.(validator.ValidationErrors)
		if ok {
			reqErrors := utils.FormatValidationError(ers)
			c.JSON(http.StatusBadRequest, utils.FormatApiError(
				"failed to process user request",
				http.StatusBadRequest,
				reqErrors,
			))
		} else {
			c.JSON(http.StatusBadRequest, utils.FormatApiError(
				"failed to process user request",
				http.StatusBadRequest, entities.E{"binding": "failed to bind request"},
			))
		}
		return
	}

	cs := services.NewDefaultCustomerService(customers.NewCustomerRepoDb(s.Db))
	apiErr := cs.CreateTempCustomerReq(cr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	var items string
	for _, item := range cr.Items {
		items += item.Title + ", "
	}

	data := map[string]string{
		"NAME":           "Customer name: " + cr.CustomerName,
		"ORDER":          "Order detail: " + items,
		"CUSTOMER_EMAIL": "Email: " + cr.CustomerEmail,
		"PHONE_NUMBER":   "Phone number: " + cr.CustomerContact,
		"ADDRESS":        "Address: " + cr.CustomerAddress,
	}

	err = mail.SendEmail(dto.MailData{
		Name:  "Wayne",
		Email: "wayne@axisforestry.com",
	}, dto.MailData{
		Name:  cr.CustomerName,
		Email: cr.CustomerEmail,
	}, data, 5)
	if err != nil {
		configs.Logger.Warn(err)
	}

	c.JSON(configs.Created, utils.GetApiResponse(gin.H{},
		"customer request saved successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) DownloadTempRequests(c *gin.Context) {
	cs := services.NewDefaultCustomerService(customers.NewCustomerRepoDb(s.Db))
	orders, apiError := cs.DownloadTempRequests()
	if apiError != nil {
		c.JSON(apiError.Code, apiError)
		return
	}
	data, err := utils.DownloadOrdersCsv(orders)
	if err != nil {
		utils.FormatApiError("Error downloading file", configs.ServerError, entities.E{})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=orders.csv")
	c.Data(http.StatusOK, "text/csv", data)
}
