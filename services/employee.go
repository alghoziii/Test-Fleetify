package services

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/domain/models"
	"Test_Fleetify/repositories"
	"Test_Fleetify/utils"
	"fmt"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(request dto.EmployeeRequest) (dto.EmployeeResponse, error)
	GetAllEmployees() ([]dto.EmployeeResponse, error)
	GetEmployeeByID(id uint) (dto.EmployeeResponse, error)
	UpdateEmployee(id uint, request dto.EmployeeRequest) (dto.EmployeeResponse, error)
	DeleteEmployee(id uint) error
}

type employeeService struct {
	employeeRepo repositories.EmployeeRepository
	db           *gorm.DB
}

func NewEmployeeService(employeeRepo repositories.EmployeeRepository, db *gorm.DB) EmployeeService {
	return &employeeService{employeeRepo: employeeRepo, db: db}
}

func (s *employeeService) CreateEmployee(request dto.EmployeeRequest) (dto.EmployeeResponse, error) {
	existingEmployee, err := s.employeeRepo.FindByNameAndDepartment(request.Name, request.DepartmentID)
	if err == nil && existingEmployee.ID != 0 {
		return dto.EmployeeResponse{}, fmt.Errorf("employee with name '%s' already exists in this department", request.Name)
	}
	// Auto-generate employee_id jika kosong
	employeeID := request.EmployeeID
	if employeeID == "" {
		generatedID, err := utils.GenerateEmployeeID(s.db, "EMP", 4)
		if err != nil {
			return dto.EmployeeResponse{}, fmt.Errorf("failed to generate employee ID: %v", err)
		}
		employeeID = generatedID
	}

	employee := models.Employee{
		EmployeeID:   employeeID,
		DepartmentID: request.DepartmentID,
		Name:         request.Name,
		Address:      request.Address,
	}

	createdEmployee, err := s.employeeRepo.CreateEmployee(employee)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	return dto.EmployeeResponse{
		ID:           createdEmployee.ID,
		EmployeeID:   createdEmployee.EmployeeID,
		DepartmentID: createdEmployee.DepartmentID,
		Name:         createdEmployee.Name,
		Address:      createdEmployee.Address,
	}, nil
}

func (s *employeeService) GetAllEmployees() ([]dto.EmployeeResponse, error) {
	employees, err := s.employeeRepo.FindAllEmployee()
	if err != nil {
		return nil, err
	}

	var responses []dto.EmployeeResponse
	for _, employee := range employees {
		responses = append(responses, dto.EmployeeResponse{
			ID:           employee.ID,
			EmployeeID:   employee.EmployeeID,
			DepartmentID: employee.DepartmentID,
			Name:         employee.Name,
			Address:      employee.Address,
		})
	}

	return responses, nil
}

func (s *employeeService) GetEmployeeByID(id uint) (dto.EmployeeResponse, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	return dto.EmployeeResponse{
		ID:           employee.ID,
		EmployeeID:   employee.EmployeeID,
		DepartmentID: employee.DepartmentID,
		Name:         employee.Name,
		Address:      employee.Address,
	}, nil
}

func (s *employeeService) UpdateEmployee(id uint, request dto.EmployeeRequest) (dto.EmployeeResponse, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	employee.EmployeeID = request.EmployeeID
	employee.DepartmentID = request.DepartmentID
	employee.Name = request.Name
	employee.Address = request.Address

	updatedEmployee, err := s.employeeRepo.UpdateEmployee(employee)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	return dto.EmployeeResponse{
		ID:           updatedEmployee.ID,
		EmployeeID:   updatedEmployee.EmployeeID,
		DepartmentID: updatedEmployee.DepartmentID,
		Name:         updatedEmployee.Name,
		Address:      updatedEmployee.Address,
	}, nil
}

func (s *employeeService) DeleteEmployee(id uint) error {
	return s.employeeRepo.Delete(id)
}
