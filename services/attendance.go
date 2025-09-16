package services

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/domain/models"
	"Test_Fleetify/repositories"
	"Test_Fleetify/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AttendanceService interface {
	ClockIn(request dto.ClockInRequest) (dto.AttendanceResponse, error)
	ClockOut(request dto.ClockOutRequest) (dto.AttendanceResponse, error)
	GetAttendanceLogs(filter dto.AttendanceFilter) ([]dto.AttendanceResponse, int64, error)
}
type attendanceService struct {
	attendanceRepo repositories.AttendanceRepository
	employeeRepo   repositories.EmployeeRepository
	db             *gorm.DB
}

func NewAttendanceService(
	attendanceRepo repositories.AttendanceRepository,
	employeeRepo repositories.EmployeeRepository,
	db *gorm.DB,
) AttendanceService {
	return &attendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		db:             db,
	}
}

func (s *attendanceService) ClockIn(request dto.ClockInRequest) (dto.AttendanceResponse, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(request.EmployeeID)
	if err != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("employee not found")
	}

	// Check if already clocked in today
	now := time.Now()
	existingAttendance, err := s.attendanceRepo.FindByEmployeeIDAndDate(request.EmployeeID, now)
	if err == nil && existingAttendance.ID != 0 {
		return dto.AttendanceResponse{}, fmt.Errorf("already clocked in today")
	}

	attendanceID, err := utils.GenerateAttendanceClockID(s.db)
	if err != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("failed to generate attendance ID: %v", err)
	}

	attendance := models.Attendance{
		EmployeeID:   request.EmployeeID,
		AttendanceID: attendanceID,
		ClockIn:      now,
		ClockOut:     nil,
	}

	createdAttendance, err := s.attendanceRepo.CreateAttendance(attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	isOnTime := true
	isLate := false
	message := "Clock In (On Time)"

	maxClockInTime, _ := time.Parse("15:04:05", employee.Department.MaxClockInTime)
	currentTime := time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	if currentTime.After(maxClockInTime) {
		isOnTime = false
		isLate = true
		message = "Clock In (Late)"
	}
	// Check if clock in is on time
	var attendanceType int8 = 1 // 1 = On Time
	if isLate {
		attendanceType = 2
	}

	// Create attendance history
	history := models.AttendanceHistory{
		EmployeeID:     request.EmployeeID,
		AttendanceID:   createdAttendance.AttendanceID,
		DateAttendance: now,
		AttendanceType: attendanceType,
		Description:    message,
	}

	_, err = s.attendanceRepo.CreateHistory(history)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.AttendanceResponse{
		ID:           createdAttendance.ID,
		EmployeeID:   createdAttendance.EmployeeID,
		AttendanceID: createdAttendance.AttendanceID,
		ClockIn:      utils.FormatTS(createdAttendance.ClockIn),
		ClockOut:     utils.FormatTSPtr(createdAttendance.ClockOut),
		CreatedAt:    utils.FormatTS(createdAttendance.CreatedAt),
		IsOnTime:     isOnTime,
		IsLate:       isLate,
		IsEarly:      false,
		Message:      message,
	}, nil
}

func (s *attendanceService) ClockOut(request dto.ClockOutRequest) (dto.AttendanceResponse, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(request.EmployeeID)
	if err != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("employee not found")
	}

	// Check if already clocked in today
	now := time.Now()
	attendance, err := s.attendanceRepo.FindByEmployeeIDAndDate(request.EmployeeID, now)
	if err != nil || attendance.ID == 0 {
		return dto.AttendanceResponse{}, fmt.Errorf("you haven't clocked in today")
	}

	// Check if already clocked out
	if attendance.ClockOut != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("already clocked out today")
	}

	// Update attendance record
	clockOutTime := now
	attendance.ClockOut = &clockOutTime

	updatedAttendance, err := s.attendanceRepo.UpdateAttendance(attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	// Check if clock out is on time
	isClockInOnTime := true
	maxClockInTime, _ := time.Parse("15:04:05", employee.Department.MaxClockInTime)
	clockInTime := time.Date(0, 0, 0, attendance.ClockIn.Hour(), attendance.ClockIn.Minute(), attendance.ClockIn.Second(), 0, time.UTC)

	if clockInTime.After(maxClockInTime) {
		isClockInOnTime = false
	}
	isOnTime := true
	isEarly := false
	message := "Clock Out (On Time)"

	maxClockOutTime, _ := time.Parse("15:04:05", employee.Department.MaxClockOutTime)
	currentTime := time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	if currentTime.Before(maxClockOutTime) {
		isOnTime = false
		isEarly = true
		message = "Clock Out (Early)"
	}

	var attendanceType int8 = 1 // 1 = On Time
	if isEarly {
		attendanceType = 3 // 3 = Early
	}

	// Create attendance history
	history := models.AttendanceHistory{
		EmployeeID:     request.EmployeeID,
		AttendanceID:   updatedAttendance.AttendanceID,
		DateAttendance: now,
		AttendanceType: attendanceType,
		Description:    message,
	}

	_, err = s.attendanceRepo.CreateHistory(history)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.AttendanceResponse{
		ID:           updatedAttendance.ID,
		EmployeeID:   updatedAttendance.EmployeeID,
		AttendanceID: updatedAttendance.AttendanceID,
		ClockIn:      utils.FormatTS(updatedAttendance.ClockIn),
		ClockOut:     utils.FormatTSPtr(updatedAttendance.ClockOut),
		CreatedAt:    utils.FormatTS(updatedAttendance.CreatedAt),
		IsOnTime:     isClockInOnTime && isOnTime,
		IsLate:       !isClockInOnTime,
		IsEarly:      isEarly,
		Message:      message,
	}, nil
}

func (s *attendanceService) GetAttendanceLogs(filter dto.AttendanceFilter) ([]dto.AttendanceResponse, int64, error) {
	attendances, total, err := s.attendanceRepo.FindAllAttendance(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.AttendanceResponse
	for _, attendance := range attendances {
		// Get employee data to check department times
		employee, err := s.employeeRepo.FindByEmployeeID(attendance.EmployeeID)
		if err != nil {
			continue // Skip if employee not found
		}

		// Check clock in status
		isClockInOnTime := true
		maxClockInTime, _ := time.Parse("15:04:05", employee.Department.MaxClockInTime)
		clockInTime := time.Date(0, 0, 0, attendance.ClockIn.Hour(), attendance.ClockIn.Minute(), attendance.ClockIn.Second(), 0, time.UTC)

		if clockInTime.After(maxClockInTime) {
			isClockInOnTime = false
		}

		// Check clock out status
		isClockOutOnTime := true
		isEarly := false
		if attendance.ClockOut != nil {
			maxClockOutTime, _ := time.Parse("15:04:05", employee.Department.MaxClockOutTime)
			clockOutTime := time.Date(0, 0, 0, attendance.ClockOut.Hour(), attendance.ClockOut.Minute(), attendance.ClockOut.Second(), 0, time.UTC)

			if clockOutTime.Before(maxClockOutTime) {
				isClockOutOnTime = false
				isEarly = true
			}
		}

		// Determine message
		message := "On Time"
		if !isClockInOnTime {
			message = "Late Arrival"
		} else if isEarly {
			message = "Early Departure"
		} else if !isClockOutOnTime && attendance.ClockOut != nil {
			message = "On Time Departure"
		}

		responses = append(responses, dto.AttendanceResponse{
			ID:           attendance.ID,
			EmployeeID:   attendance.EmployeeID,
			AttendanceID: attendance.AttendanceID,
			ClockIn:      utils.FormatTS(attendance.ClockIn),
			ClockOut:     utils.FormatTSPtr(attendance.ClockOut),
			CreatedAt:    utils.FormatTS(attendance.CreatedAt),
			IsOnTime:     isClockInOnTime && isClockOutOnTime,
			IsLate:       !isClockInOnTime,
			IsEarly:      isEarly,
			Message:      message,
		})
	}

	return responses, total, nil
}
