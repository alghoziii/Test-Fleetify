package controllers

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EmployeeController struct {
	employeeService services.EmployeeService
}

func NewEmployeeController(employeeService services.EmployeeService) *EmployeeController {
	return &EmployeeController{employeeService: employeeService}
}

func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var request dto.EmployeeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	response, err := c.employeeService.CreateEmployee(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create employee", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.BuildResponse(true, "Employee created successfully", response))
}

func (c *EmployeeController) GetAllEmployees(ctx *gin.Context) {
	employees, err := c.employeeService.GetAllEmployees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get employees", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Employees retrieved successfully", employees))
}

func (c *EmployeeController) GetEmployeeByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	employee, err := c.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.BuildErrorResponse("Employee not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Employee retrieved successfully", employee))
}

func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	var request dto.EmployeeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	employee, err := c.employeeService.UpdateEmployee(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update employee", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Employee updated successfully", employee))
}

func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	err = c.employeeService.DeleteEmployee(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to delete employee", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Employee deleted successfully", nil))
}
