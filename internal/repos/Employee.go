package repos

import (
	"context"
	"database/sql" // For sql.ErrNoRows

	"github.com/jmoiron/sqlx"
	"githup.ahmedramadan.4cashier/internal/models"
)

// EmployeeRepository defines the interface for employee data operations.
type EmployeeRepository interface {
	GetEmployees(ctx context.Context, branchID int) ([]models.Employee, error) // Added branchID for filtering
	GetEmployee(ctx context.Context, id int) (models.Employee, error)
	AddEmployee(ctx context.Context, employee models.Employee) (int, error)
	UpdateEmployee(ctx context.Context, employee models.Employee) (int, error)
	DeleteEmployee(ctx context.Context, id int) (int, error)
	CreateEmployeeTx(ctx context.Context, tx *sqlx.Tx, employee models.Employee) (int, error) // Transactional add
}

// PosEmployeeRepository implements EmployeeRepository for PostgreSQL.
type PosEmployeeRepository struct {
	DB *sqlx.DB
}

// NewPosEmployeeRepository creates a new instance of PosEmployeeRepository.
func NewPosEmployeeRepository(db *sqlx.DB) EmployeeRepository {
	return &PosEmployeeRepository{DB: db}
}

// GetEmployees retrieves a list of employees, filtered by branch ID.
func (r *PosEmployeeRepository) GetEmployees(ctx context.Context, branchID int) ([]models.Employee, error) {
	data := []models.Employee{}
	// Adjust query to filter by branch_id
	err := r.DB.SelectContext(ctx, &data, "SELECT * FROM employees WHERE branch_id = $1", branchID)
	if err != nil {
		return []models.Employee{}, err
	}
	return data, err
}

// GetEmployee retrieves a single employee by ID.
func (r *PosEmployeeRepository) GetEmployee(ctx context.Context, id int) (models.Employee, error) {
	data := new(models.Employee)
	err := r.DB.GetContext(ctx, &data, "SELECT * FROM employees WHERE id =$1", id)
	if err != nil {
		return models.Employee{}, err
	}
	return *data, err
}

// AddEmployee adds a new employee to the database.
func (r *PosEmployeeRepository) AddEmployee(ctx context.Context, employee models.Employee) (int, error) {
	var insertedID int

	query := `
		INSERT INTO employees (user_id, branch_id, employee_code, hire_date, job_title, salary, is_active, role,
			performance_score, last_performance_review, commission_rate, target_sales, achievements, supervisor_id,
			employment_status, created_at, updated_at)
		VALUES (:user_id, :branch_id, :employee_code, :hire_date, :job_title, :salary, :is_active, :role,
			:performance_score, :last_performance_review, :commission_rate, :target_sales, :achievements, :supervisor_id,
			:employment_status, :created_at, :updated_at)
		RETURNING id;
	`

	rows, err := r.DB.NamedQueryContext(ctx, query, employee)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&insertedID); err != nil {
			return 0, err
		}
	} else if err = rows.Err(); err != nil {
		return 0, err
	} else {
		return 0, sql.ErrNoRows // No rows returned, indicating insertion might have failed or RETURNING clause issue
	}

	return insertedID, nil
}

// UpdateEmployee updates an existing employee in the database.
func (r *PosEmployeeRepository) UpdateEmployee(ctx context.Context, employee models.Employee) (int, error) {
	query := `
		UPDATE employees SET user_id = :user_id, branch_id = :branch_id, employee_code = :employee_code,
		hire_date = :hire_date, job_title = :job_title, salary = :salary, is_active = :is_active,
		role = :role, performance_score = :performance_score, last_performance_review = :last_performance_review,
		commission_rate = :commission_rate, target_sales = :target_sales, achievements = :achievements,
		supervisor_id = :supervisor_id, employment_status = :employment_status, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.DB.NamedExecContext(ctx, query, employee)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

// CreateEmployeeTx adds a new employee within a transaction.
func (r *PosEmployeeRepository) CreateEmployeeTx(ctx context.Context, tx *sqlx.Tx, employee models.Employee) (int, error) {
	var id int
	query := `
		INSERT INTO employees (user_id, branch_id, employee_code, hire_date, job_title, salary, is_active, role,
			performance_score, last_performance_review, commission_rate, target_sales, achievements, supervisor_id,
			employment_status, created_at, updated_at)
		VALUES (:user_id, :branch_id, :employee_code, :hire_date, :job_title, :salary, :is_active, :role,
			:performance_score, :last_performance_review, :commission_rate, :target_sales, :achievements, :supervisor_id,
			:employment_status, :created_at, :updated_at)
		RETURNING id;
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &id, employee)
	return id, err
}

// DeleteEmployee deletes an employee from the database by ID.
func (r *PosEmployeeRepository) DeleteEmployee(ctx context.Context, employeeId int) (int, error) {
	result, err := r.DB.ExecContext(ctx, "delete from employees where id=$1", employeeId)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(id), err
}
