package controllers

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AttendanceController struct {
	attendanceService services.AttendanceService
}

func NewAttendanceController(attendanceService services.AttendanceService) *AttendanceController {
	return &AttendanceController{attendanceService: attendanceService}
}

func (h *AttendanceController) ClockIn(c *gin.Context) {
	var request dto.ClockInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	response, err := h.attendanceService.ClockIn(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to clock in", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.BuildResponse(true, "Berhasil absen masuk", response))
}

func (h *AttendanceController) ClockOut(c *gin.Context) {
	var request dto.ClockOutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid request", err.Error()))
		return
	}

	response, err := h.attendanceService.ClockOut(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to clock out", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.BuildResponse(true, "Berhasil Absen Keluar", response))
}

func (h *AttendanceController) GetAttendanceLogs(c *gin.Context) {
	var filter dto.AttendanceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	// Set default values jika tidak disediakan
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	attendances, total, err := h.attendanceService.GetAttendanceLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get attendance logs", err.Error()))
		return
	}

	response := gin.H{
		"data":       attendances,
		"total":      total,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"totalPages": (total + int64(filter.Limit) - 1) / int64(filter.Limit),
	}

	c.JSON(http.StatusOK, dto.BuildResponse(true, "Attendance logs retrieved successfully", response))
}
