package models

import "time"

type Employee struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   string     `gorm:"unique;size:50" json:"employee_id"`
	DepartmentID uint       `json:"department_id"`
	Name         string     `gorm:"size:255" json:"name"`
	Address      string     `gorm:"type:text" json:"address"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

type Department struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	DepartmentName  string `gorm:"size:255" json:"department_name"`
	MaxClockInTime  string `gorm:"type:varchar(8)" json:"max_clock_in_time"`
	MaxClockOutTime string `gorm:"type:varchar(8)" json:"max_clock_out_time"`
}

type Attendance struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   string     `gorm:"size:50" json:"employee_id"`
	AttendanceID string     `gorm:"unique;size:50" ,json:"attendance_id"`
	ClockIn      time.Time  `json:"clock_in"`
	ClockOut     *time.Time `json:"clock_out"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Employee     Employee   `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}

type AttendanceHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     string    `gorm:"size:50" json:"employee_id"`
	AttendanceID   string    `gorm:"size:100;index" json:"attendance_id"`
	DateAttendance time.Time `json:"date_attendance"`
	AttendanceType int8      `gorm:"type:tinyint(1)" json:"attendance_type"`
	Description    string    `gorm:"type:text" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Employee       Employee  `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}
