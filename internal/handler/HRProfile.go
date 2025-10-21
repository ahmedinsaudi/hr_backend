package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"githup.ahmedramadan.4cashier/internal/bootstrap"
	"githup.ahmedramadan.4cashier/internal/models"
	mylogger "githup.ahmedramadan.4cashier/internal/mylogger"
	"githup.ahmedramadan.4cashier/internal/service" 
)

// HRService Interface from the service layer is needed here to type the Service field correctly.
// Assuming your service package is named 'services' (plural) to avoid confusion with the interface.

type HRHandler struct {
	Logger  zerolog.Logger
	Service *service.HRService 
}

func NewHRHandler(logger zerolog.Logger, service *service.HRService) *HRHandler {
	return &HRHandler{
		Logger:  logger.With().Str("layer", "handler").Str("component", "HRHandler").Logger(),
		Service: service,
	}
}

// ------------------------------------------------------------------
// GET /hr-profiles (عرض ملفات الـ HR)
// ------------------------------------------------------------------
func (h *HRHandler) GetHRProfiles(ctx fiber.Ctx) error {
	pagination := bootstrap.GetPagination(ctx)

	// Note: job_position is now correctly used for filtering, not sorting.
	// You might want to introduce a separate "sort" query param for sorting.
	filters := map[string]interface{}{
		"searchText":   ctx.Query("searchText"),
		"company_name": ctx.Query("company_name"),
		"job_position": ctx.Query("job_position"),
		"verified":     parseBoolOrDefault(ctx.Query("verified"), false),
	}

	log.Println("Filters:", filters)

	items, err := h.Service.GetHRProfiles(ctx.Context(), pagination, filters)
	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to fetch HR profiles")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to fetch HR profiles"})
	}

	return ctx.JSON(fiber.Map{
		"items": items,
	})
}

// ------------------------------------------------------------------
// GET /rates (عرض التقييمات)
// ------------------------------------------------------------------
func (h *HRHandler) GetRates(ctx fiber.Ctx) error {
	pagination := bootstrap.GetPagination(ctx)

	// A slight adjustment to ensure correct type parsing for filters
	filters := map[string]interface{}{
		"hr_profile_id": parseIntOrDefault(ctx.Query("hr_profile_id"), 0),
		"employee_id":   parseIntOrDefault(ctx.Query("employee_id"), 0),
		"min_rate":      parseFloat32OrDefault(ctx.Query("min_rate"), 0.0), // Changed to float32
		"max_rate":      parseFloat32OrDefault(ctx.Query("max_rate"), 5.0), // Changed to float32
		"is_verified":   parseBoolOrDefault(ctx.Query("is_verified"), false),
		"review_text":   ctx.Query("review_text"),
		"sort_column":   ctx.Query("sort_column"),
	}

	log.Println("Filters:", filters)

	// items will be []models.RateWithDetails
	items, err := h.Service.GetRates(ctx.Context(), pagination, filters)
	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to fetch rates")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to fetch rates"})
	}

	return ctx.JSON(fiber.Map{
		"items": items,
	})
}

// ------------------------------------------------------------------
// POST /hr/rate (تقييم HR)
// ------------------------------------------------------------------
func (h *HRHandler) RateHR(ctx fiber.Ctx) error {
	var rate models.Rate
	if err := ctx.Bind().Body(&rate); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid rate data"})
	}

	// ⭐️ ENHANCEMENT: RateHR service returns the full HRProfile (potentially with new badges)
	profile, err := h.Service.RateHR(ctx.Context(), &rate) 
	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to rate HR")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to rate HR"})
	}

	// Return the updated profile which includes the new average rate and potentially badges
	return ctx.JSON(profile)
}

// ------------------------------------------------------------------
// POST /hr/rate/like (إعجاب/عدم إعجاب بتقييم)
// ------------------------------------------------------------------
func (h *HRHandler) LikeRate(ctx fiber.Ctx) error {
	var like models.RateLike
	if err := ctx.Bind().Body(&like); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid like data"})
	}

	// The service returns the new likes count, useful for frontend update
	newLikesCount, err := h.Service.LikeRate(ctx.Context(), &like)
	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to like rate")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to like rate"})
	}

	return ctx.JSON(fiber.Map{"likes_count": newLikesCount})
}

// ------------------------------------------------------------------
// POST /hr/:id/experience (إضافة خبرة)
// ------------------------------------------------------------------
func (h *HRHandler) AddExperience(ctx fiber.Ctx) error {
	hrID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid HR ID"})
	}

	var experiences []models.Experience
	if err := ctx.Bind().Body(&experiences); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid experience data"})
	}

	if err := h.Service.AddExperience(ctx.Context(), hrID, experiences); err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to add experiences")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to add experiences"})
	}

	return ctx.SendStatus(204) // No Content on successful update/creation
}

// ------------------------------------------------------------------
// POST /hr/:id/job-roles (إضافة أدوار وظيفية)
// ------------------------------------------------------------------
func (h *HRHandler) AddJobRoles(ctx fiber.Ctx) error {
	hrID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid HR ID"})
	}

	var roles []models.JobRole
	if err := ctx.Bind().Body(&roles); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid job role data"})
	}

	if err := h.Service.AddJobRoles(ctx.Context(), hrID, roles); err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to add job roles")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to add job roles"})
	}

	return ctx.SendStatus(204)
}

// ------------------------------------------------------------------
// POST /hr/badge/like (إعجاب/عدم إعجاب بشارة)
// ------------------------------------------------------------------
func (h *HRHandler) LikeBadge(ctx fiber.Ctx) error {
	var like models.BadgeLike
	if err := ctx.Bind().Body(&like); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid like data"})
	}

	// Assuming LikeBadge service returns error only
	if err := h.Service.LikeBadge(ctx.Context(), &like); err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to like badge")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to like badge"})
	}

	return ctx.SendStatus(204)
}

// ------------------------------------------------------------------
// POST /hr/badge (منح شارة يدوياً - اختياري)
// ------------------------------------------------------------------
func (h *HRHandler) AwardBadge(ctx fiber.Ctx) error {
	var badge models.Badge
	if err := ctx.Bind().Body(&badge); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid badge data"})
	}

	// Note: AwardBadge is usually done by the system (RateHR), 
	// but this manual handler remains for admin use.
	id, err := h.Service.AwardBadge(ctx.Context(), &badge)
	if err != nil {
		mylogger.HandleLogging(h.Logger, err, "Failed to award badge")
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to award badge"})
	}

	return ctx.JSON(fiber.Map{"id": id, "message": "Badge awarded manually"})
}

// ------------------------------------------------------------------
// ⭐️ Helper Functions (تم إضافة Float32)
// ------------------------------------------------------------------
func parseBoolOrDefault(str string, def bool) bool {
	if str == "true" {
		return true
	}
	if str == "false" {
		return false
	}
	return def
}

func parseIntOrDefault(str string, def int) int {
	if val, err := strconv.Atoi(str); err == nil {
		return val
	}
	return def
}

func parseFloat32OrDefault(str string, def float32) float32 {
	if val, err := strconv.ParseFloat(str, 32); err == nil {
		return float32(val)
	}
	return def
}


func (h *HRHandler) GetEmployeeStats(c fiber.Ctx) error {
    // 1. استخراج employee_id من المسار
    employeeIDStr := c.Params("employee_id")
    employeeID, err := strconv.Atoi(employeeIDStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
    }

    // 2. جلب الإحصائيات من Repository
    stats, err := h.Service.GetEmployeeStats(c.Context(), employeeID)
    if err != nil {
        // ... (تسجيل الخطأ)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve stats"})
    }

    // 3. الإرجاع إلى العميل
    return c.JSON(stats)
}