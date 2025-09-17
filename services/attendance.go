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
        return dto.AttendanceResponse{}, fmt.Errorf("Karyawan tidak ditemukan")
    }

    now := time.Now()
    existingAttendance, err := s.attendanceRepo.FindByEmployeeIDAndDate(request.EmployeeID, now)
    if err == nil && existingAttendance.ID != 0 {
        return dto.AttendanceResponse{}, fmt.Errorf("Sudah Absen hari ini")
    }

    attendanceID, err := utils.GenerateAttendanceClockID(s.db)
    if err != nil {
        return dto.AttendanceResponse{}, fmt.Errorf("failed to generate attendance ID: %v", err)
    }

    // DEBUG: Log waktu untuk troubleshooting
    fmt.Printf("Waktu sekarang: %v\n", now)
    fmt.Printf("Max clock in time dari DB: %v\n", employee.Department.MaxClockInTime)

    // Parse max clock in time dari department
    maxClockInTimeStr := employee.Department.MaxClockInTime
    var maxClockInTime time.Time

    // Coba parse dengan format HH:MM dulu (yang paling umum)
    maxClockInTime, err = time.Parse("15:04", maxClockInTimeStr)
    if err != nil {
        // Jika gagal, coba parse dengan format HH:MM:SS
        maxClockInTime, err = time.Parse("15:04:05", maxClockInTimeStr)
        if err != nil {
            return dto.AttendanceResponse{}, fmt.Errorf("format waktu department tidak valid: %v", maxClockInTimeStr)
        }
    }

    // Gunakan timezone yang sama dengan waktu sekarang
    loc := now.Location()
    
    // Buat waktu comparison dengan tanggal yang sama
    maxClockInToday := time.Date(
        now.Year(), now.Month(), now.Day(),
        maxClockInTime.Hour(), maxClockInTime.Minute(), maxClockInTime.Second(), 0, loc,
    )

    fmt.Printf("Max clock in today: %v\n", maxClockInToday)
    fmt.Printf("Perbandingan: now (%v) > maxClockInToday (%v) = %v\n", 
        now, maxClockInToday, now.After(maxClockInToday))

    // Check jika terlambat
    isLate := now.After(maxClockInToday)
    message := "Absen Masuk (Ontime)"
    attendanceType := 1 // Ontime

    if isLate {
        message = "Absen Masuk (Terlambat)"
        attendanceType = 2 // Late
        
        // Hitung keterlambatan dalam menit
        lateDuration := now.Sub(maxClockInToday)
        lateMinutes := int(lateDuration.Minutes())
        message = fmt.Sprintf("Absen Masuk (Terlambat %d menit)", lateMinutes)
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

    // Create attendance history
    history := models.AttendanceHistory{
        EmployeeID:     request.EmployeeID,
        AttendanceID:   createdAttendance.AttendanceID,
        DateAttendance: now,
        AttendanceType: int8(attendanceType),
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
        Message:      message,
        EmployeeName: employee.Name,
        Department:   employee.Department.DepartmentName,
        IsLate:       isLate,
    }, nil
}

func (s *attendanceService) ClockOut(request dto.ClockOutRequest) (dto.AttendanceResponse, error) {
	employee, err := s.employeeRepo.FindByEmployeeID(request.EmployeeID)
	if err != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("Karyawan tidak ditemukan")
	}

	now := time.Now()
	attendance, err := s.attendanceRepo.FindByEmployeeIDAndDate(request.EmployeeID, now)
	if err != nil || attendance.ID == 0 {
		return dto.AttendanceResponse{}, fmt.Errorf("Belum melakukan absen masuk hari ini")
	}

	// Check if already clocked out
	if attendance.ClockOut != nil {
		return dto.AttendanceResponse{}, fmt.Errorf("Sudah absen keluar hari ini")
	}

	// Update attendance record
	clockOutTime := now
	attendance.ClockOut = &clockOutTime

	updatedAttendance, err := s.attendanceRepo.UpdateAttendance(attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	// Hitung apakah clock in terlambat (untuk response)
	// Parse max clock in time dari department
	maxClockInTimeStr := employee.Department.MaxClockInTime
	var maxClockInTime time.Time
	var isLate bool = false

	// Coba parse dengan format HH:MM dulu
	maxClockInTime, err = time.Parse("15:04", maxClockInTimeStr)
	if err != nil {
		// Jika gagal, coba parse dengan format HH:MM:SS
		maxClockInTime, err = time.Parse("15:04:05", maxClockInTimeStr)
	}

	if err == nil {
		// Gunakan timezone yang sama dengan waktu clock in
		loc := updatedAttendance.ClockIn.Location()
		
		// Buat waktu comparison dengan tanggal yang sama dengan clock in
		maxClockInToday := time.Date(
			updatedAttendance.ClockIn.Year(), 
			updatedAttendance.ClockIn.Month(), 
			updatedAttendance.ClockIn.Day(),
			maxClockInTime.Hour(), 
			maxClockInTime.Minute(), 
			maxClockInTime.Second(), 0, loc,
		)

		// Check jika clock in terlambat
		isLate = updatedAttendance.ClockIn.After(maxClockInToday)
	}

	message := "Absen Keluar"
	attendanceType := 3 // Clock Out

	// Create attendance history
	history := models.AttendanceHistory{
		EmployeeID:     request.EmployeeID,
		AttendanceID:   updatedAttendance.AttendanceID,
		DateAttendance: now,
		AttendanceType: int8(attendanceType),
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
		Message:      message,
		EmployeeName: employee.Name,
		Department:   employee.Department.DepartmentName,
		IsLate:       isLate, 
	}, nil
}
func (s *attendanceService) GetAttendanceLogs(filter dto.AttendanceFilter) ([]dto.AttendanceResponse, int64, error) {
    attendances, total, err := s.attendanceRepo.FindAllAttendance(filter)
    if err != nil {
        return nil, 0, err
    }

    responses := make([]dto.AttendanceResponse, 0, len(attendances))

    for _, a := range attendances {
        emp := a.Employee
        dept := emp.Department

        responses = append(responses, dto.AttendanceResponse{
            ID:           a.ID,
            EmployeeID:   a.EmployeeID,
            EmployeeName: emp.Name,
            AttendanceID: a.AttendanceID,
            ClockIn:      utils.FormatTS(a.ClockIn),
            ClockOut:     utils.FormatTSPtr(a.ClockOut),
            CreatedAt:    utils.FormatTS(a.CreatedAt),
            Message:      "Attendance Record",
            Department:   dept.DepartmentName,
        })
    }

    return responses, total, nil
}