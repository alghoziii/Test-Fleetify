package repositories

import (
	"Test_Fleetify/domain/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee models.Employee) (models.Employee, error)
	FindAllEmployee() ([]models.Employee, error)
	FindByID(id uint) (models.Employee, error)
	FindByEmployeeID(employeeID string) (models.Employee, error)
	UpdateEmployee(employee models.Employee) (models.Employee, error)
	FindByNameAndDepartment(name string, departmentID uint) (models.Employee, error)
	Delete(id uint) error
}
type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee models.Employee) (models.Employee, error) {
	err := r.db.Create(&employee).Error
	return employee, err
}
func (r *employeeRepository) FindByID(id uint) (models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("Department").First(&employee, id).Error
	return employee, err
}
func (r *employeeRepository) FindAllEmployee() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Preload("Department").Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) FindByEmployeeID(employeeID string) (models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("Department").Where("employee_id = ?", employeeID).First(&employee).Error
	return employee, err
}

func (r *employeeRepository) UpdateEmployee(employee models.Employee) (models.Employee, error) {
	err := r.db.Save(&employee).Error
	return employee, err
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

func (r *employeeRepository) FindByNameAndDepartment(name string, departmentID uint) (models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("name = ? AND department_id = ?", name, departmentID).First(&employee).Error
	return employee, err
}
