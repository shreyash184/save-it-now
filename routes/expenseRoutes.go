package routes

import (
	"expense-tracker/controllers"
	"expense-tracker/middleware"
	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/expense", controllers.AddExpense)
		protected.GET("/expenses", controllers.GetExpenses)
	}
}
