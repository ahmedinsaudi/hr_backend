package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// Experience history for HR profiles
type Experience struct {
	Name        *string    `json:"name"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	JobPosition *string    `json:"job_position"`
}

// Tasks or roles assigned to HR profiles
type JobRole struct {
	Name            *string    `json:"name"`
	RoleDescription *string    `json:"role_description"`
	StartDate       *time.Time `json:"start_date"`
	DoneRate        *int       `json:"done_rate"`  // Changed from DoneRate time.Time to int rating
	Visible         *bool      `json:"visible"`
}

// HRProfile represents an HR user on the platform.
type HRProfile struct {
	ID               int          `db:"id" json:"id"`
	Name             *string      `db:"name" json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email            *string      `db:"email" json:"email,omitempty" validate:"omitempty,email"`
	Image            *string      `db:"image" json:"image" validate:"omitempty"`
	PasswordHash  string    `db:"password_hash" json:"-"  validate:"omitempty"` 
	Skills *[]string `db:"skills" json:"skills,omitempty"`
	CompanyName      *string      `db:"company_name" json:"company_name,omitempty" validate:"omitempty,min=2,max=100"`
	JobPosition      *string      `db:"job_position" json:"job_position,omitempty" validate:"omitempty,min=2,max=100"`
	Experience       []Experience `db:"experience" json:"experience,omitempty"`
	JobRoles         []JobRole    `db:"job_roles" json:"job_roles,omitempty"`
	Rate             *float32     `db:"rate" json:"rate,omitempty"`
	TotalRatesCount  int          `db:"total_rates_count" json:"total_rates_count"`
	Verified         bool         `db:"verified_profile" json:"verified_profile"`
	CreatedAt        time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at" json:"updated_at"`


	Badges []Badge `db:"-" json:"badges,omitempty"`
}

// Optional legacy task model (if needed separately)
type HRProfileTask struct {
	ID             int       `db:"id" json:"id"`
	HRProfileID    int       `db:"hr_profile_id" json:"hr_profile_id" validate:"required,gt=0"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	TaskTitle      string    `db:"task_title" json:"task_title" validate:"required,min=2,max=200"`
	Description    string    `db:"description" json:"description" validate:"required,min=5"`
	Visible        bool      `db:"visible_for_people" json:"visible_for_people"`
	LikesCount     int       `db:"likes_count" json:"likes_count"`
}

// Employee is a regular user who can rate HRs.
type Employee struct {
	ID           int        `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Image            *string      `db:"image" json:"image" validate:"omitempty"`
	JobField     string     `db:"job_field" json:"job_field"`
	Address      *string    `db:"address" json:"address"`
	City         *string    `db:"city" json:"city"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"password_hash"`
	ActivePoints bool       `db:"active_points" json:"active_points"`
	IsVerified   bool       `db:"is_verified" json:"is_verified"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}


// Badge is awarded to high-performing HRs
type Badge struct {
	ID               int       `db:"id" json:"id"`
	HRProfileID      int       `db:"hr_profile_id" json:"hr_profile_id" validate:"required,gt=0"`
	CreatedDate      time.Time `db:"created_date" json:"created_date"` // fixed key from created_add
	TotalRates       int       `db:"total_rates_number" json:"total_rates_number"`
	Rate             float32   `db:"rate" json:"rate"`
	JobPosition      string    `db:"job_position" json:"job_position"`
	CurrentJobRoles  string    `db:"current_job_roles" json:"current_job_roles"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// Rating by employee to HR
type Rate struct {
	ID            int       `db:"id" json:"id"`
	HRProfileID   int       `db:"hr_profile_id" json:"hr_profile_id" validate:"required,gt=0"`
	EmployeeID    int       `db:"employee_id" json:"employee_id" validate:"required,gt=0"`
	ReviewText    string    `db:"review_text" json:"review_text" validate:"required,min=5,max=2000"`
	RateValue     float32   `db:"rate_value" json:"rate_value" validate:"gte=0,lte=5"`
	RatingContext *string   `db:"rating_context" json:"rating_context,omitempty"`
	LikesCount    int       `db:"likes_count" json:"likes_count"`
	IsVerified    bool       `db:"is_verified" json:"is_verified"` 
	HRResponse    *string   `db:"hr_response" json:"hr_response,omitempty"`
	IsAnonymous   bool      `db:"is_anonymous" json:"is_anonymous"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`

	

    
}

type RateWithDetails struct {
    Rate 

    ProfileID         int       `db:"profile_id" json:"profile_id"`
	ProfileName       string    `db:"profile_name" json:"profile_name"` 
	ProfileEmail      string    `db:"profile_email" json:"profile_email"`
	CompanyName       string    `db:"company_name" json:"company_name"`
	JobPosition       string    `db:"job_position" json:"job_position"`
	ProfileRate       float32   `db:"profile_rate" json:"profile_rate"`
	TotalRatesCount   int       `db:"total_rates_count" json:"total_rates_count"`
	VerifiedProfile   bool      `db:"verified_profile" json:"verified_profile"`

	BadgeID         int       `db:"badge_id" json:"badge_id"`
    BadgeRate       float32   `db:"badge_rate" json:"badge_rate"`

    EmployeeName      *string    `db:"employee_name" json:"employee_name"`
    EmployeeImage *string    `db:"employee_image" json:"employee_image"`
}

type EmployeeStats struct {
    TotalRatingsCount int     `json:"total_ratings_count"`
    ContributorRank   string  `json:"contributor_rank"` // يمكن أن تكون "Top 10%"
    TotalLikesCount   int     `json:"total_likes_count"`
}

type RateWithEmployee struct {
    Rate 

    EmployeeName   string    `db:"employee_name" json:"employee_name"`
    EmployeeAvatar string    `db:"employee_avatar" json:"employee_avatar"`
}

// Like/dislike on a rate
type RateLike struct {
	ID         int       `db:"id" json:"id"`
	RateID     int       `db:"rate_id" json:"rate_id" validate:"required,gt=0"`
	EmployeeID int       `db:"employee_id" json:"employee_id" validate:"required,gt=0"`
	IsLike     bool      `db:"is_like" json:"is_like"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Like/dislike on a badge (fixed typo)
type BadgeLike struct {
	ID         int       `db:"id" json:"id"`
	BadgeID    int       `db:"badge_id" json:"badge_id" validate:"required,gt=0"`
	EmployeeID int       `db:"employee_id" json:"employee_id" validate:"required,gt=0"`
	IsLike     bool      `db:"is_like" json:"is_like"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
