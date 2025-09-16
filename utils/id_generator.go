package utils

import (
	"Test_Fleetify/domain/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GenerateEmployeeID(db *gorm.DB, prefix string, length int) (string, error) {
	var lastEmployee struct {
		EmployeeID string
	}

	// Cari employee_id terakhir yang sesuai dengan prefix
	result := db.Table("employees").
		Where("employee_id LIKE ?", prefix+"%").
		Order("employee_id DESC").
		Limit(1).
		Scan(&lastEmployee)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return "", result.Error
	}

	// Jika tidak ada data, mulai dari 1
	if lastEmployee.EmployeeID == "" {
		return fmt.Sprintf("%s%0*d", prefix, length, 1), nil
	}

	// Extract number dari employee_id terakhir
	lastID := lastEmployee.EmployeeID
	if strings.HasPrefix(lastID, prefix) {
		numberStr := lastID[len(prefix):]
		number, err := strconv.Atoi(numberStr)
		if err == nil {
			nextNumber := number + 1
			return fmt.Sprintf("%s%0*d", prefix, length, nextNumber), nil
		}
	}

	// Fallback: jika format tidak sesuai, mulai dari 1
	return fmt.Sprintf("%s%0*d", prefix, length, 1), nil
}

// GenerateAttendanceID menghasilkan attendance_id dalam format ATT-EMP001-20250115
func GenerateAttendanceID(employeeID string) string {
	return fmt.Sprintf("ATT-%s-%s", employeeID, time.Now().Format("20060102"))
}

// GenerateAttendanceID menghasilkan attendance_id sederhana: ATT-001
func GenerateAttendanceClockID(db *gorm.DB) (string, error) {
	var lastAttendance models.Attendance
	result := db.Order("id DESC").First(&lastAttendance)

	nextNumber := 1
	if result.Error == nil && lastAttendance.ID > 0 {
		nextNumber = int(lastAttendance.ID) + 1
	}

	return fmt.Sprintf("ATT-%03d", nextNumber), nil
}
