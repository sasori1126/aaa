package routes

import (
	"net/http"

	account2 "axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/controllers"
	"axis/ecommerce-backend/pkg/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func StartApp(c *controllers.Serve) *gin.Engine {
	userService := services.NewDefaultUserService(account2.NewUserRepoDb(c.Db))
	gin.SetMode(gin.ReleaseMode)
	apiRoute := gin.Default()
	//apiRoute.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
	//	APIKey:          "53ed68f7dd1b7ccc5902aee30e6ac982",
	//	ProjectPackages: []string{"main", "axis/ecommerce-backend"},
	//}))
	apiRoute.HandleMethodNotAllowed = true
	apiRoute.Use(otelgin.Middleware("Axisforestry"))
	apiRoute.Use(middlewares.CORS())
	apiRoute.NoRoute(middlewares.NoRoute())
	apiRoute.NoMethod(middlewares.NoMethod())
	apiRoute.MaxMultipartMemory = 8 << 20
	apiRoute.GET("/s3img", c.S3Resource)
	v1 := apiRoute.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", c.Login)
			auth.POST("/register", c.RegisterUser)
			auth.POST("/token/refresh", c.Refresh)
			auth.POST("/password/reset", c.ResetPassword)
			auth.POST("/password/forgot", c.ForgotPassword)
		}

		// ====== user account
		account := v1.Group("/account").Use(middlewares.AuthorizeJwt(services.NewDefaultUserService(account2.NewUserRepoDb(c.Db))))
		{
			account.GET("/profile", c.UserProfile)
			account.POST("/profile/update", c.UpdateUserProfile)
			account.POST("/equipments", c.AddEquipment)
			account.GET("/orders", c.UserOrders)
			account.POST("/address", c.AddAddress)
			account.PUT("/address/:addressId/update", c.UpdateAddress)
			account.DELETE("/address/:addressId/delete", c.DeleteAddress)
			account.PUT("/address/:addressId/set-default", c.SetDefaultAddress)
			account.GET("/address", c.UserAddresses)
			account.GET("/equipments", c.GetUserEquipments)
			account.GET("/logout", c.Logout)
			account.PUT("/password/reset", c.AccountResetPassword)
		}

		// ========= cart endpoints
		cart := v1.Group("/carts").Use(middlewares.AuthorizeJwt(userService))
		{
			cart.POST("/place-order", c.MakePayment)
			cart.POST("/add-cart-item", c.AddCartItem)
			cart.POST("/increase-cart-item-quantity", c.IncreaseCartItem)
			cart.POST("/decrease-cart-item-quantity", c.DecreaseCartItem)
			cart.DELETE("/delete-cart-items", c.RemoveCartItem)
			cart.GET("/my-cart", c.GetUserActiveCart)
			cart.GET("/delivery-rates/:addressId", c.GetDeliveryRates)
		}

		// ========= converge pay endpoints
		convergepay := v1.Group("/converge").Use(middlewares.AuthorizeJwt(userService))
		{
			convergepay.POST("/token", c.GetConvergePayToken)
		}

		// ========= users endpoints
		user := v1.Group("/users").Use(middlewares.AuthorizeJwt(userService))
		{
			user.GET("/", c.GetAllUsers)
		}

		// ========= public ecommerce endpoints
		eco := v1.Group("/")
		{
			eco.GET("/models", c.GetModels)
			eco.GET("/parts", c.GetParts)
			eco.GET("/parts/search", c.SearchParts)
			eco.GET("/parts/:partId", c.GetPartById)
			eco.GET("/categories", c.GetDiagramCats)
			eco.GET("/head-types", c.GetHeadTypes)
			eco.GET("/heads", c.GetHeads)
			eco.POST("/heads/kesla/orders", c.KeslaHeadOrder) // move to log in
			eco.POST("/heads/axis/orders", c.AxisHeadOrder)
			eco.POST("/heads/controller/orders", c.ControllerOrder)
			eco.GET("/heads/head-type/:headTypeId", c.GetHeadByHeadType)
			eco.GET("/categories/:categoryId", c.GetDiagramCat)
			eco.GET("/heads/:headId", c.GetHeadById)
			eco.GET("/distributors", c.GetDistributors)
			eco.GET("/diagrams", c.GetDiagrams)
			eco.GET("/diagrams/:diagramId", c.GetDiagramById)
			eco.GET("/diagrams/search", c.GetDiagramsSearch)
			eco.GET("/distributors/:distributorId", c.GetDistributor)
		}

		// ======= admin endpoints
		admin := v1.Group("/admin")
		{
			distributors := admin.Group("/distributors").Use(middlewares.AuthorizeJwt(userService))
			{
				distributors.POST("", c.CreateDistributor)
				distributors.GET("", c.GetDistributors)
				distributors.GET("/:distributorId/show", c.GetDistributor)
				distributors.PUT("/:distributorId/update", c.UpdateDistributor)
				distributors.DELETE("/:distributorId/delete", c.DeleteDistributor)
			}

			tool := admin.Group("/tools").Use(middlewares.AuthorizeJwt(userService), middlewares.UserRole([]string{"SuperAdmin", "Admin"}))
			{
				tool.GET("/price-difference", c.LoadPrices)
				tool.GET("/duplicates", c.GetDuplicates)
				tool.PATCH("/price-update", c.PriceUpdate)
				tool.POST("/upload-price-list", c.UploadFile)
				tool.POST("/merge-parts", c.MergeParts)
				tool.POST("/diagrams", c.SaveDiagram)
				tool.PUT("/diagrams/:diagramId/update", c.UpdateDiagram)
				tool.DELETE("/diagrams/:diagramId", c.DeleteDiagram)
				tool.POST("/diagrams/images", c.SaveDiagramImages)
				tool.PUT("/diagrams/images", c.UpdateDiagramImages)
				tool.GET("/diagrams/images/add-from-csv", c.AddDiagramImagesFromCsv)
				tool.DELETE("/diagrams/images/:imageId", c.DeleteDiagramImage)
				tool.GET("/diagrams/images", c.GetFigureImages)
			}

			// ========= heads endpoints
			heads := admin.Group("/heads").Use(middlewares.AuthorizeJwt(userService))
			{
				heads.POST("", c.CreateHead)
				heads.POST("/types", c.CreateHeadType)
			}

			// ========= parts endpoints
			parts := admin.Group("/parts").Use(middlewares.AuthorizeJwt(userService))
			{
				parts.POST("", c.CreatePart)
				parts.GET("", c.GetParts)
				parts.GET("/:partId", c.GetPartById)
				parts.PUT("/:partId/update", c.UpdatePart)
			}

			// ========= diagram endpoints
			diagram := admin.Group("/diagrams").Use(middlewares.AuthorizeJwt(userService))
			{
				diagram.POST("", c.CreateDiagram)
				diagram.GET("", c.GetAdminDiagrams)
				diagram.GET("/:diagramId", c.GetDiagramById)
				diagram.POST("/categories", c.CreateDiagramCat)
				diagram.GET("/categories", c.GetDiagramCats)
				diagram.DELETE("/categories/:catId", c.DeleteDiagramCat)
				diagram.PUT("/categories/:catId/update", c.UpdateDiagramCat)
				diagram.POST("/sub-categories", c.CreateDiagramSubCat)
			}

			// ========= manufacturer endpoints
			manufacturer := admin.Group("/manufacturers")
			{
				manufacturer.POST("", c.CreateManufacturer)
				manufacturer.PUT("/:manufacturerId/update", c.UpdateManufacturer)
				manufacturer.DELETE("/:manufacturerId", c.DeleteManufacturer)
				manufacturer.GET("", c.GetManufacturers)
			}

			// ========= taxes admin endpoints
			tax := admin.Group("/taxes").Use(middlewares.AuthorizeJwt(userService))
			{
				tax.POST("", c.CreateTax)
				tax.GET("/:userId", c.GetTaxesForTaxExemptions)
				tax.PUT("/:userId/:taxId/update", c.SaveTaxExemption)
				tax.DELETE("/:userId/:taxId", c.DeleteTaxExemption)
			}

			stats := admin.Group("/stats").Use(middlewares.AuthorizeJwt(userService))
			{
				stats.GET("/dashboard", c.DashboardStats)
				stats.GET("/cache/:name", c.CacheInvalidate)
			}

			// ========= models endpoints
			models := admin.Group("/models").Use(middlewares.AuthorizeJwt(userService))
			{
				models.POST("", c.CreateModel)
				models.GET("", c.GetModels)
				models.GET("/:modelId", c.GetModel)
				models.PUT("/:modelId/update", c.UpdateModel)
				models.DELETE("/:modelId", c.DeleteModel)
			}

			// ======== images endpoints
			images := admin.Group("/images").Use(middlewares.AuthorizeJwt(userService))
			{
				images.DELETE("/:imageId", c.DeleteImage)
			}

			// ========= orders endpoints
			orders := admin.Group("/orders").Use(middlewares.AuthorizeJwt(userService))
			{
				orders.POST("/add-payment", c.AddOrderPayment)
				orders.GET("", c.GetOrders)
				orders.GET("/:orderID", c.GetOrderByID)
				orders.PATCH("/:orderID/update-status", c.OrderUpdateStatus)
			}

			carts := admin.Group("/carts").Use(middlewares.AuthorizeJwt(userService))
			{
				carts.GET("", c.GetCarts)
				carts.GET("/:cartID", c.GetCartByID)
				carts.PATCH("/:cartID/update", c.OrderUpdateStatus)
			}

			// ========= controllers endpoints
			ctl := admin.Group("/controllers").Use(middlewares.AuthorizeJwt(userService), middlewares.UserRole([]string{"SuperAdmin", "Admin", "Staff"}))
			{
				ctl.POST("", c.CreateController)
				ctl.PUT("/:controllerId/update", c.UpdateController)
				ctl.DELETE("/:controllerId", c.DeleteController)
				ctl.GET("", c.GetControllers)
			}

			// ========= users endpoints
			users := admin.Group("/users").Use(middlewares.AuthorizeJwt(userService), middlewares.UserRole([]string{"SuperAdmin", "Admin"}))
			{
				users.GET("", c.GetAllUsers)
				users.GET("/:userId/show", c.GetUser)
				users.POST("/ask-to-reset-password", c.AskUserToResetPassword)
				users.PATCH("/assign-role", c.AssignUserRoleOnAccount)
				users.PATCH("/on-account", c.UpdateUserOnAccount)
				users.GET("/search", c.SearchUsers)
				users.PATCH("/update-user-currency", c.UpdateUserCurrency)
			}
		}
	}

	v1.POST("/customers/requests/temp", c.TempCustomerRequest)
	v1.GET("/customers/orders/temp", c.DownloadTempRequests)
	v1.GET("/account/verification/token/:token", c.ValidateAccountEmail)
	v1.POST("/account/verification/token", c.RequestNewEmailVerificationToken)
	v1.POST("/contact-us", c.ContactUs)
	v1.GET("/orders", c.GetOrders)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "pong",
		})
	})

	apiRoute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return apiRoute
}
