package controllers

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	departmentService services.DepartmentService
}

func NewDepartmentController(departmentService services.DepartmentService) *DepartmentController {
	return &DepartmentController{departmentService: departmentService}
}

func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
	var request dto.DepartmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	response, err := c.departmentService.CreateDepartment(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create department", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.BuildResponse(true, "Department created successfully", response))
}

func (c *DepartmentController) GetAllDepartments(ctx *gin.Context) {
	departments, err := c.departmentService.GetAllDepartments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get departments", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Departments retrieved successfully", departments))
}

func (c *DepartmentController) GetDepartmentByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	department, err := c.departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.BuildErrorResponse("Department not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Department retrieved successfully", department))
}

func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	var request dto.DepartmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	department, err := c.departmentService.UpdateDepartment(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update department", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Department updated successfully", department))
}

func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid ID", err.Error()))
		return
	}

	err = c.departmentService.DeleteDepartment(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to delete department", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.BuildResponse(true, "Department deleted successfully", nil))
}
