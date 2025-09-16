package routes

import (
	"Test_Fleetify/controllers"
	"Test_Fleetify/repositories"
	"Test_Fleetify/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	employeeRepo := repositories.NewEmployeeRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)

	// Initialize services
	employeeService := services.NewEmployeeService(employeeRepo, db)
	departmentService := services.NewDepartmentService(departmentRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo, employeeRepo, db)

	// Initialize controllers
	employeeController := controllers.NewEmployeeController(employeeService)
	departmentController := controllers.NewDepartmentController(departmentService)
	attendanceController := controllers.NewAttendanceController(attendanceService)

	// Employee routes
	employeeRoutes := router.Group("/employees")
	{
		employeeRoutes.POST("", employeeController.CreateEmployee)
		employeeRoutes.GET("", employeeController.GetAllEmployees)
		employeeRoutes.GET("/:id", employeeController.GetEmployeeByID)
		employeeRoutes.PUT("/:id", employeeController.UpdateEmployee)
		employeeRoutes.DELETE("/:id", employeeController.DeleteEmployee)
	}

	// Department routes
	departmentRoutes := router.Group("/departments")
	{
		departmentRoutes.POST("", departmentController.CreateDepartment)
		departmentRoutes.GET("", departmentController.GetAllDepartments)
		departmentRoutes.GET("/:id", departmentController.GetDepartmentByID)
		departmentRoutes.PUT("/:id", departmentController.UpdateDepartment)
		departmentRoutes.DELETE("/:id", departmentController.DeleteDepartment)
	}

	// Attendance routes
	attendanceRoutes := router.Group("/attendance")
	{
		attendanceRoutes.POST("/clock-in", attendanceController.ClockIn)
		attendanceRoutes.PUT("/clock-out", attendanceController.ClockOut)
		attendanceRoutes.GET("/logs", attendanceController.GetAttendanceLogs)
	}
}
