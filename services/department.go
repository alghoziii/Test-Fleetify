package services

import (
	"Test_Fleetify/domain/dto"
	"Test_Fleetify/domain/models"
	"Test_Fleetify/repositories"
)

type DepartmentService interface {
	CreateDepartment(request dto.DepartmentRequest) (dto.DepartmentResponse, error)
	GetAllDepartments() ([]dto.DepartmentResponse, error)
	GetDepartmentByID(id uint) (dto.DepartmentResponse, error)
	UpdateDepartment(id uint, request dto.DepartmentRequest) (dto.DepartmentResponse, error)
	DeleteDepartment(id uint) error
}

type departmentService struct {
	departmentRepo repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepo repositories.DepartmentRepository) DepartmentService {
	return &departmentService{departmentRepo: departmentRepo}
}

func (s *departmentService) CreateDepartment(request dto.DepartmentRequest) (dto.DepartmentResponse, error) {
	department := models.Department{
		DepartmentName:  request.DepartmentName,
		MaxClockInTime:  request.MaxClockInTime,
		MaxClockOutTime: request.MaxClockOutTime,
	}

	createdDepartment, err := s.departmentRepo.CreateDepartment(department)
	if err != nil {
		return dto.DepartmentResponse{}, err
	}

	return dto.DepartmentResponse{
		ID:              createdDepartment.ID,
		DepartmentName:  createdDepartment.DepartmentName,
		MaxClockInTime:  createdDepartment.MaxClockInTime,
		MaxClockOutTime: createdDepartment.MaxClockOutTime,
	}, nil
}

func (s *departmentService) GetAllDepartments() ([]dto.DepartmentResponse, error) {
	departments, err := s.departmentRepo.FindAllDepartment()
	if err != nil {
		return nil, err
	}

	var responses []dto.DepartmentResponse
	for _, department := range departments {
		responses = append(responses, dto.DepartmentResponse{
			ID:              department.ID,
			DepartmentName:  department.DepartmentName,
			MaxClockInTime:  department.MaxClockInTime,
			MaxClockOutTime: department.MaxClockOutTime,
		})
	}

	return responses, nil
}

func (s *departmentService) GetDepartmentByID(id uint) (dto.DepartmentResponse, error) {
	department, err := s.departmentRepo.FindDepartmentByID(id)
	if err != nil {
		return dto.DepartmentResponse{}, err
	}

	return dto.DepartmentResponse{
		ID:              department.ID,
		DepartmentName:  department.DepartmentName,
		MaxClockInTime:  department.MaxClockInTime,
		MaxClockOutTime: department.MaxClockOutTime,
	}, nil
}

func (s *departmentService) UpdateDepartment(id uint, request dto.DepartmentRequest) (dto.DepartmentResponse, error) {
	department, err := s.departmentRepo.FindDepartmentByID(id)
	if err != nil {
		return dto.DepartmentResponse{}, err
	}

	department.DepartmentName = request.DepartmentName
	department.MaxClockInTime = request.MaxClockInTime
	department.MaxClockOutTime = request.MaxClockOutTime

	updatedDepartment, err := s.departmentRepo.UpdateDepartment(department)
	if err != nil {
		return dto.DepartmentResponse{}, err
	}

	return dto.DepartmentResponse{
		ID:              updatedDepartment.ID,
		DepartmentName:  updatedDepartment.DepartmentName,
		MaxClockInTime:  updatedDepartment.MaxClockInTime,
		MaxClockOutTime: updatedDepartment.MaxClockOutTime,
	}, nil
}

func (s *departmentService) DeleteDepartment(id uint) error {
	return s.departmentRepo.Delete(id)
}
