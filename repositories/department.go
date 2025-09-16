package repositories

import (
	"Test_Fleetify/domain/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	CreateDepartment(department models.Department) (models.Department, error)
	FindAllDepartment() ([]models.Department, error)
	FindDepartmentByID(id uint) (models.Department, error)
	UpdateDepartment(department models.Department) (models.Department, error)
	Delete(id uint) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) CreateDepartment(department models.Department) (models.Department, error) {
	err := r.db.Create(&department).Error
	return department, err
}

func (r *departmentRepository) FindAllDepartment() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) FindDepartmentByID(id uint) (models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	return department, err
}

func (r *departmentRepository) UpdateDepartment(department models.Department) (models.Department, error) {
	err := r.db.Save(&department).Error
	return department, err
}

func (r *departmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}
