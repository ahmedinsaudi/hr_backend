package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"githup.ahmedramadan.4cashier/internal/models"
	mylogger "githup.ahmedramadan.4cashier/internal/mylogger"
	"githup.ahmedramadan.4cashier/internal/service"
)

// EmployeeHandler handles HTTP requests related to Employee entities.
type EmployeeHandler struct {
	Logger  zerolog.Logger
	Service *service.EmployeeService
}

// NewEmployeeHandler creates a new instance of EmployeeHandler.
func NewEmployeeHandler(logger zerolog.Logger, serv *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{Logger: logger.With().Str("layer", "handler").Logger(), Service: serv}
}

// GetEmployees retrieves a list of employees, optionally filtered by branch ID.
// It expects a "branch_id" query parameter.
func (h *EmployeeHandler) GetEmployees(ctx fiber.Ctx) error {
	branchIdStr := ctx.Query("branch_id")
	branchID, _ := strconv.Atoi(branchIdStr) // Using _ for error, mirroring template

	log.Printf("dfffffffffffffffffffffffffffffffffffff") // Mirroring template's log
	data, err := h.Service.GetEmployees(ctx.Context(), branchID)

	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to fetch Employees from service")
		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"}) // Mirroring template's error response
	}

	return ctx.JSON(data)
}

// GetEmployee retrieves a single employee by their ID.
// It expects an "employee_id" query parameter.
func (h *EmployeeHandler) GetEmployee(ctx fiber.Ctx) error {
	employeeIdStr := ctx.Query("employee_id")
	id, _ := strconv.Atoi(employeeIdStr) // Using _ for error, mirroring template
	data, err := h.Service.GetEmployee(ctx.Context(), id)

	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to get Employee from service")
		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"}) // Mirroring template's error response
	}

	return ctx.JSON(data)
}

// AddEmployee adds a new employee.
// It expects the employee data in the request body, parsed via ctx.Locals("bodyParse").
func (h *EmployeeHandler) AddEmployee(ctx fiber.Ctx) error {
	employee, ok := ctx.Locals("bodyParse").(*models.Employee)
	if !ok || employee == nil {
		mylogger.HandleLogging(h.Logger, nil, "Failed to parse employee body for AddEmployee")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee data"})
	}

	insertedID, err := h.Service.AddEmployee(ctx.Context(), *employee)

	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to add Employee through service")
		return ctx.Status(401).JSON(fiber.Map{"error": err.Error()}) // Mirroring template's error response
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": insertedID})
}

// UpdateEmployee updates an existing employee.
// It expects the employee data in the request body, parsed via ctx.Locals("bodyParse").
func (h *EmployeeHandler) UpdateEmployee(ctx fiber.Ctx) error {
	employee, ok := ctx.Locals("bodyParse").(*models.Employee)
	if !ok || employee == nil {
		mylogger.HandleLogging(h.Logger, nil, "Failed to parse employee body for UpdateEmployee")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee data"})
	}

	rowsAffected, err := h.Service.UpdateEmployee(ctx.Context(), *employee)

	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to update Employee through service")
		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"}) // Mirroring template's error response
	}

	return ctx.JSON(fiber.Map{"rows_affected": rowsAffected})
}

// DeleteEmployee deletes an employee by their ID.
// It expects the employee ID as a path parameter (e.g., /employees/:id).
func (h *EmployeeHandler) DeleteEmployee(ctx fiber.Ctx) error {
	employeeIdStr := ctx.Params("id")
	id, _ := strconv.Atoi(employeeIdStr) // Using _ for error, mirroring template
	rowsAffected, err := h.Service.DeleteEmployee(ctx.Context(), id)

	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to delete Employee through service")
		return ctx.Status(401).JSON(fiber.Map{"error": "Failed to delete Employee from repository"}) // Mirroring template's error response
	}

	return ctx.JSON(fiber.Map{"rows_affected": rowsAffected})
}
