package dto

type EmployeeRequest struct {
	EmployeeID   string `json:"employee_id"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
}

type EmployeeResponse struct {
	ID           uint   `json:"id"`
	EmployeeID   string `json:"employee_id"`
	DepartmentID uint   `json:"department_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
}
type DepartmentRequest struct {
	DepartmentName  string `json:"department_name" binding:"required"`
	MaxClockInTime  string `json:"max_clock_in_time" binding:"required"`
	MaxClockOutTime string `json:"max_clock_out_time" binding:"required"`
}

type DepartmentResponse struct {
	ID              uint   `json:"id"`
	DepartmentName  string `json:"department_name"`
	MaxClockInTime  string `json:"max_clock_in_time"`
	MaxClockOutTime string `json:"max_clock_out_time"`
}
type ClockInRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
}

type ClockOutRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
}

type AttendanceResponse struct {
	ID           uint    `json:"id"`
	EmployeeID   string  `json:"employee_id"`
	AttendanceID string  `json:"attendance_id"`
	ClockIn      string  `json:"clock_in"`  // Format: "2006-01-02 15:04"
	ClockOut     *string `json:"clock_out"` // Pointer agar bisa null
	CreatedAt    string  `json:"created_at"`
	IsOnTime     bool    `json:"is_on_time"`
	IsLate       bool    `json:"is_late"`
	IsEarly      bool    `json:"is_early"`
	Message      string  `json:"message"`
}

type AttendanceHistoryResponse struct {
	ID             uint   `json:"id"`
	EmployeeID     string `json:"employee_id"`
	AttendanceID   string `json:"attendance_id"`
	DateAttendance string `json:"date_attendance"`
	AttendanceType int8   `json:"attendance_type"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
}

type AttendanceFilter struct {
	Date       string `form:"date"`
	Department uint   `form:"department"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
