package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"githup.ahmedramadan.4cashier/internal/bootstrap"
	"githup.ahmedramadan.4cashier/internal/models"
)

type HRRepository interface {
	// CRUD/Update Functions
	AddExperience(ctx context.Context, hrID int, exp []models.Experience) error
	AddJobRoles(ctx context.Context, hrID int, roles []models.JobRole) error
	AwardBadge(ctx context.Context, badge *models.Badge) (int, error)
	
	// Core Business Logic Handlers (Atomic Transactions)
	RateHR(ctx context.Context, rate *models.Rate) (int, float32, error)
	LikeRate(ctx context.Context, like *models.RateLike) (int, error) 
	LikeBadge(ctx context.Context, like *models.BadgeLike) error
    
	// Retrieval Functions
	GetHRProfiles(ctx context.Context, pagination bootstrap.Pagination, filters map[string]interface{}) ([]models.HRProfile, error)
	GetRates(ctx context.Context, pagination bootstrap.Pagination, filters map[string]interface{}) ([]models.RateWithDetails, error)
    
    // Helper Functions for Service Logic
	GetHRProfileByID(ctx context.Context, hrID int) (*models.HRProfile, error)
	CheckIfProfileHasBadge(ctx context.Context, profileID int, badgeName string) (bool, error)
	GetRateOwner(ctx context.Context, rateID int) (int, error)
	UpdateEmployeePoints(ctx context.Context, employeeID int, pointsToAdd int) error

	GetEmployeeStats(ctx context.Context, employeeID int) (models.EmployeeStats, error)
}

type PosHRRepository struct {
	DB *sqlx.DB
}

func NewPosHRRepository(db *sqlx.DB) HRRepository {
	return &PosHRRepository{DB: db}
}



func (r *PosHRRepository) GetHRProfiles(
	ctx context.Context,
	pagination bootstrap.Pagination,
	filters map[string]interface{},
) ([]models.HRProfile, error) {

	var hrs []models.HRProfile
	args := []interface{}{}
	conditions := []string{}
	argPos := 1

	// البحث العام (searchText)
	if searchText, ok := filters["searchText"].(string); ok && searchText != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR company_name ILIKE $%d OR email ILIKE $%d )", argPos, argPos))
		args = append(args, "%"+searchText+"%")
		argPos++
	}


	// فلترة بحسب الشركة
	if companyName, ok := filters["company_name"].(string); ok && companyName != "" {
		conditions = append(conditions, fmt.Sprintf("company_name ILIKE $%d", argPos))
		args = append(args, "%"+companyName+"%")
		argPos++
	}

	// فلترة بحسب الوظيفة
	if jobPosition, ok := filters["job_position"].(string); ok && jobPosition != "" {
		conditions = append(conditions, fmt.Sprintf("job_position ILIKE $%d", argPos))
		args = append(args, "%"+jobPosition+"%")
		argPos++
	}

	// فلترة بحسب حالة التوثيق
	if verified, ok := filters["verified"].(bool); ok {
		conditions = append(conditions, fmt.Sprintf("verified_profile = $%d", argPos))
		args = append(args, verified)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}


	

	sortColumn := "created_at"
sortOrder := "DESC"
// يجب استخدام حقل مخصص للفرز مثل "sort" أو "orderBy"
if sortFilter, ok := filters["sort"].(string); ok && sortFilter != "" { 
    sortParts := strings.Split(sortFilter, ":")
    if len(sortParts) >= 1 && sortParts[0] != "" {
        sortColumn = sortParts[0] // عمود الفرز
    }
    if len(sortParts) == 2 {
		sortColumn = sortParts[0]
		if sortParts[1] == "asc" || sortParts[1] == "ASC" {
			sortOrder = "ASC"
		}
    }
}
	
	sortClause := fmt.Sprintf("ORDER BY %s %s", sortColumn, sortOrder)

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	args = append(args, pagination.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, name, email,image, company_name, job_position, rate, total_rates_count, verified_profile,
		 created_at, updated_at FROM hr_profiles
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortClause, argPos, argPos+1)

	err := r.DB.SelectContext(ctx, &hrs, query, args...)
	if err != nil {
		return nil, err
	}

	return hrs, nil
}

func (r *PosHRRepository) GetRates(
	ctx context.Context,
	pagination bootstrap.Pagination,
	filters map[string]interface{},
) ([]models.RateWithDetails, error) {
	var rates []models.RateWithDetails
	args := []interface{}{}
	conditions := []string{}
	argPos := 1

	if employeeID, ok := filters["employee_id"].(int); ok && employeeID > 0 {
		conditions = append(conditions, fmt.Sprintf("r.employee_id = $%d", argPos))
		args = append(args, employeeID)
		argPos++
	}

	if hrProfileID, ok := filters["hr_profile_id"].(int); ok && hrProfileID > 0 {
		conditions = append(conditions, fmt.Sprintf("r.hr_profile_id = $%d", argPos))
		args = append(args, hrProfileID)
		argPos++
	}

	if isVerified, ok := filters["is_verified"].(bool); ok {
		conditions = append(conditions, fmt.Sprintf("r.is_verified = $%d", argPos))
		args = append(args, isVerified)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sortColumn := "r.created_at"
	sortOrder := "DESC"

	offset := (pagination.Page - 1) * pagination.Limit
	args = append(args, pagination.Limit, offset)

	limitArgPos := argPos
	offsetArgPos := argPos + 1

	query := fmt.Sprintf(`
		SELECT 
            r.id, r.hr_profile_id, r.employee_id, r.review_text, r.rate_value, 
            r.likes_count, r.is_verified, r.is_anonymous, r.created_at,
            
            p.id AS profile_id, p.name AS profile_name, p.company_name, 
            p.job_position, p.rate AS profile_rate, p.total_rates_count, p.verified_profile, 
            
            b.id AS badge_id, b.rate AS badge_rate,
            
            e.name AS employee_name, e.image AS employee_image
            
        FROM rates r
        JOIN hr_profiles p ON r.hr_profile_id = p.id
        LEFT JOIN employees e ON r.employee_id = e.id
        
        LEFT JOIN (
            SELECT 
                *,
                ROW_NUMBER() OVER(PARTITION BY hr_profile_id ORDER BY created_at DESC) as rn
            FROM badges
        ) b ON b.hr_profile_id = p.id AND b.rn = 1 

        %s
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d
    `, whereClause, sortColumn, sortOrder, limitArgPos, offsetArgPos)

	err := r.DB.SelectContext(ctx, &rates, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates with details: %w", err)
	}

	return rates, nil
}


func (r *PosHRRepository) AddExperience(ctx context.Context, hrID int, newExp []models.Experience) error {
	if len(newExp) == 0 {
		return nil
	}

	var existingExp []models.Experience

	// 1. جلب الخبرات الموجودة حاليًا
	querySelect := `SELECT experience FROM hr_profiles WHERE id = $1`
	err := r.DB.GetContext(ctx, &existingExp, querySelect, hrID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing experience: %w", err)
	}

	if existingExp == nil {
		existingExp = []models.Experience{}
	}


	updatedExp := append(existingExp, newExp...)

	// 3. تحديث العمود بالخبرات الجديدة المدموجة
	queryUpdate := `
		UPDATE hr_profiles
		SET experience = $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.DB.ExecContext(ctx, queryUpdate, updatedExp, hrID)
	if err != nil {
		return fmt.Errorf("failed to update experience JSONB: %w", err)
	}

	return nil
}


func (r *PosHRRepository) AddJobRoles(ctx context.Context, hrID int, newRoles []models.JobRole) error {
	if len(newRoles) == 0 {
		return nil
	}

	var existingRoles []models.JobRole

	querySelect := `SELECT job_roles FROM hr_profiles WHERE id = $1`
	err := r.DB.GetContext(ctx, &existingRoles, querySelect, hrID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing job roles: %w", err)
	}


	if existingRoles == nil {
		existingRoles = []models.JobRole{}
	}

	// 2. دمج الأدوار الجديدة مع الموجودة
	updatedRoles := append(existingRoles, newRoles...)

	// 3. تحديث العمود بالـ job_roles الجديدة المدموجة
	queryUpdate := `
		UPDATE hr_profiles
		SET job_roles = $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.DB.ExecContext(ctx, queryUpdate, updatedRoles, hrID)
	if err != nil {
		return fmt.Errorf("failed to update job roles JSONB: %w", err)
	}

	return nil
}


// =================================================================
// ⭐️ Core Business Logic Implementations (Transactions)
// =================================================================

func (r *PosHRRepository) RateHR(ctx context.Context, rate *models.Rate) (int, float32, error) {
	var newAverageRate float32
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert the new rate
	queryInsertRate := `
        INSERT INTO rates (hr_profile_id, employee_id, review_text, rate_value, rating_context, likes_count, is_verified, hr_response, is_anonymous, created_at)
        VALUES (:hr_profile_id, :employee_id, :review_text, :rate_value, :rating_context, 0, :is_verified, :hr_response, :is_anonymous, NOW())
        RETURNING id
    `
	// ⭐️ FIX 1: استخدام PrepareNamedContext ثم ExecContext لتسجيل البيانات في Transaction
	// PrepareNamedContext هي الطريقة الأصح للـ sqlx داخل الـ Transaction
	stmt, err := tx.PrepareNamedContext(ctx, queryInsertRate)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to prepare named statement: %w", err)
	}
	defer stmt.Close()
	
	rows, err := stmt.QueryxContext(ctx, rate)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute named query: %w", err)
	}
	
	if rows.Next() {
		if err := rows.Scan(&rate.ID); err != nil {
			rows.Close()
			return 0, 0, fmt.Errorf("failed to scan new rate ID: %w", err)
		}
	}
	rows.Close()

	// 2. Update the HRProfile's average rate and count (Atomic Calculation)
	queryUpdateProfile := `
        UPDATE hr_profiles
        SET 
            total_rates_count = total_rates_count + 1,
            rate = ((rate * total_rates_count) + $1) / (total_rates_count + 1),
            updated_at = NOW()
        WHERE id = $2
        RETURNING rate
    `
	err = tx.GetContext(ctx, &newAverageRate, queryUpdateProfile, rate.RateValue, rate.HRProfileID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update hr_profile average: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rate.ID, newAverageRate, nil
}


func (r *PosHRRepository) LikeRate(ctx context.Context, like *models.RateLike) (int, error) {
	var newCount int
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. UPSERT the like
	queryUpsertLike := `
        INSERT INTO rate_likes (rate_id, employee_id, is_like, created_at)
        VALUES (:rate_id, :employee_id, :is_like, NOW())
        ON CONFLICT (rate_id, employee_id) 
        DO UPDATE SET 
            is_like = EXCLUDED.is_like,
            created_at = NOW()
    `
	_, err = tx.NamedExecContext(ctx, queryUpsertLike, like)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert rate_like: %w", err)
	}

	// 2. Recalculate the 'likes_count' and return the new value
	queryRecount := `
        UPDATE rates
        SET likes_count = (
            SELECT COUNT(*) FROM rate_likes
            WHERE rate_id = $1 AND is_like = true
        )
        WHERE id = $1
        RETURNING likes_count
    `
	err = tx.GetContext(ctx, &newCount, queryRecount, like.RateID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to recount likes for rate: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	
	return newCount, nil 
}

func (r *PosHRRepository) LikeBadge(ctx context.Context, like *models.BadgeLike) error {

	query := `
        INSERT INTO badge_likes (badge_id, employee_id, is_like, created_at)
        VALUES (:badge_id, :employee_id, :is_like, NOW())
        ON CONFLICT (badge_id, employee_id) 
        DO UPDATE SET 
            is_like = EXCLUDED.is_like,
            created_at = NOW()
    `
	_, err := r.DB.NamedExecContext(ctx, query, like)
	if err != nil {
		return fmt.Errorf("failed to upsert badge_like: %w", err)
	}
	
	// (Note: You might want to update the badge's likes_count atomically here 
	// similar to LikeRate, but for simplicity, we keep it as a simple upsert.)
	
	return nil
}

// =================================================================
// ⭐️ Helper Functions Implementation
// =================================================================

func (r *PosHRRepository) GetHRProfileByID(ctx context.Context, hrID int) (*models.HRProfile, error) {
	var profile models.HRProfile
	query := `SELECT id, name, email, company_name, job_position, rate, total_rates_count, verified_profile, updated_at FROM hr_profiles WHERE id = $1`
	err := r.DB.GetContext(ctx, &profile, query, hrID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HR profile %d: %w", hrID, err)
	}
	return &profile, nil
}



func (r *PosHRRepository) CheckIfProfileHasBadge(ctx context.Context, profileID int, badgeName string) (bool, error) {
	var count int
	query := `SELECT COUNT(id) FROM badges WHERE hr_profile_id = $1 AND job_position = $2`
	err := r.DB.GetContext(ctx, &count, query, profileID, badgeName)
	if err != nil {
		return false, fmt.Errorf("failed to check for badge %s: %w", badgeName, err)
	}
	return count > 0, nil
}

func (r *PosHRRepository) GetRateOwner(ctx context.Context, rateID int) (int, error) {
	var employeeID int
	err := r.DB.GetContext(ctx, &employeeID, "SELECT employee_id FROM rates WHERE id = $1", rateID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rate owner for rate %d: %w", rateID, err)
	}
	return employeeID, nil
}

func (r *PosHRRepository) UpdateEmployeePoints(ctx context.Context, employeeID int, pointsToAdd int) error {
	// Assumes an 'employees' table with an 'active_points' column
	query := `UPDATE employees SET active_points = active_points + $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, pointsToAdd, employeeID)
	if err != nil {
		return fmt.Errorf("failed to update points for employee %d: %w", employeeID, err)
	}
	return nil
}

func (r *PosHRRepository) AwardBadge(ctx context.Context, badge *models.Badge) (int, error) {
	query := `
        INSERT INTO badges (hr_profile_id, created_date, total_rates_number, rate, job_position, current_job_roles, created_at, updated_at)
        VALUES (:hr_profile_id, :created_date, :total_rates_number, :rate, :job_position, :current_job_roles, NOW(), NOW())
        RETURNING id
    `
	rows, err := r.DB.NamedQueryContext(ctx, query, badge)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&badge.ID); err != nil {
			return 0, err
		}
	}
	return badge.ID, nil
}


// internal/repos/employee_repo.go

func (r *PosHRRepository) GetEmployeeStats(ctx context.Context, employeeID int) (models.EmployeeStats, error) {
    stats := models.EmployeeStats{}

    // 1. حساب إجمالي التقييمات والإعجابات
    query1 := `
        SELECT
            COUNT(r.id) AS total_ratings_count,
            COALESCE(SUM(r.likes_count), 0) AS total_likes_count
        FROM rates r
        WHERE r.employee_id = $1
    `
    // ⚠️ يجب استخدام QueryRow أو DB.GetContext إذا كنت تستخدم sqlx
    if err := r.DB.QueryRowContext(ctx, query1, employeeID).Scan(
        &stats.TotalRatingsCount, 
        &stats.TotalLikesCount,
    ); err != nil {
        return models.EmployeeStats{}, fmt.Errorf("failed to fetch rate counts: %w", err)
    }

    // 2. حساب التصنيف المئوي (ContributorRank)
    // هذا مثال بسيط، المنطق الحقيقي أكثر تعقيداً ويتطلب نافذة وظيفية (Window Function)
    
    // لتبسيط الأمر الآن، نفترض أن أي شخص لديه أكثر من 5 تقييمات هو Top 10%
    if stats.TotalRatingsCount >= 5 {
        stats.ContributorRank = "Top 10%"
    } else if stats.TotalRatingsCount > 0 {
        stats.ContributorRank = "Contributor"
    } else {
        stats.ContributorRank = "Newbie"
    }


    return stats, nil
}