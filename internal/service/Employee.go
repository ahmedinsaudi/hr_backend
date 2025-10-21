package service

import (
	"context"

	"github.com/rs/zerolog"
	"githup.ahmedramadan.4cashier/internal/models"
	"githup.ahmedramadan.4cashier/internal/repos"
)

// EmployeeService provides business logic for Employee entities.
type EmployeeService struct {
	Logger zerolog.Logger
	Repo   repos.EmployeeRepository // Changed to EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(logger zerolog.Logger, repo repos.EmployeeRepository) *EmployeeService {
	return &EmployeeService{Logger: logger.With().Str("layer", "service").Logger(), Repo: repo}
}

// GetEmployees retrieves a list of employees from the repository.
func (s *EmployeeService) GetEmployees(ctx context.Context, branchID int) ([]models.Employee, error) {
	data, err := s.Repo.GetEmployees(ctx, branchID) // Pass branchID to repository
	if err != nil {
		serviceLogger := s.Logger.With().Caller().Logger()
		serviceLogger.Error().Err(err).Int("branchID", branchID).Msg("Failed to fetch Employees from repository")
	}
	return data, err
}

// GetEmployee retrieves a single employee by ID from the repository.
func (s *EmployeeService) GetEmployee(ctx context.Context, id int) (models.Employee, error) {
	data, err := s.Repo.GetEmployee(ctx, id)
	if err != nil {
		serviceLogger := s.Logger.With().Caller().Logger()
		serviceLogger.Error().Err(err).Int("employeeID", id).Msg("Failed to fetch Employee from repository by ID")
	}
	return data, err
}

// AddEmployee adds a new employee to the repository.
func (s *EmployeeService) AddEmployee(ctx context.Context, employee models.Employee) (int, error) {
	data, err := s.Repo.AddEmployee(ctx, employee)
	if err != nil {
		serviceLogger := s.Logger.With().Caller().Logger()
		serviceLogger.Error().Err(err).Msg("Failed to Add Employee to repository")
	}
	return data, err
}

// UpdateEmployee updates an existing employee in the repository.
func (s *EmployeeService) UpdateEmployee(ctx context.Context, employee models.Employee) (int, error) {
	data, err := s.Repo.UpdateEmployee(ctx, employee)
	if err != nil {
		serviceLogger := s.Logger.With().Caller().Logger()
		serviceLogger.Error().Err(err).Int("employeeID", employee.ID).Msg("Failed to Update Employee in repository")
	}

	return data, err
}

// DeleteEmployee deletes an employee from the repository by ID.
func (s *EmployeeService) DeleteEmployee(ctx context.Context, id int) (int, error) {
	data, err := s.Repo.DeleteEmployee(ctx, id)
	if err != nil {
		serviceLogger := s.Logger.With().Caller().Logger()
		serviceLogger.Error().Err(err).Int("employeeID", id).Msg("Failed to delete Employee from repository")
	}
	return data, err
}
