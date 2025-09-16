package repositories

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/domain/models"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	CreateAttendance(attendance models.Attendance) (models.Attendance, error)
	FindByEmployeeIDAndDate(employeeID string, date time.Time) (models.Attendance, error)
	UpdateAttendance(attendance models.Attendance) (models.Attendance, error)
	CreateHistory(history models.AttendanceHistory) (models.AttendanceHistory, error)
	FindAllAttendance(filter dto.AttendanceFilter) ([]models.Attendance, int64, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) CreateAttendance(attendance models.Attendance) (models.Attendance, error) {
	attendance.ClockOut = nil
	err := r.db.Create(&attendance).Error
	return attendance, err
}

func (r *attendanceRepository) FindByEmployeeIDAndDate(employeeID string, date time.Time) (models.Attendance, error) {
	var attendance models.Attendance
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("employee_id = ? AND clock_in >= ? AND clock_in < ?",
		employeeID, startOfDay, endOfDay).First(&attendance).Error
	return attendance, err
}

func (r *attendanceRepository) UpdateAttendance(attendance models.Attendance) (models.Attendance, error) {
	err := r.db.Save(&attendance).Error
	return attendance, err
}

func (r *attendanceRepository) CreateHistory(history models.AttendanceHistory) (models.AttendanceHistory, error) {
	err := r.db.Create(&history).Error
	return history, err
}

func (r *attendanceRepository) FindAllAttendance(filter dto.AttendanceFilter) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.Model(&models.Attendance{}).Preload("Employee").Preload("Employee.Department")

	if filter.Date != "" {
		date, err := time.Parse("2006-01-02", filter.Date)
		if err == nil {
			startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where("clock_in >= ? AND clock_in < ?", startOfDay, endOfDay)
		}
	}

	if filter.Department > 0 {
		query = query.Joins("JOIN employees ON attendances.employee_id = employees.employee_id").
			Where("employees.department_id = ?", filter.Department)
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(filter.Limit).Find(&attendances).Error
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}
