package service

import (
	"context"
	"fmt"
	"time"

	"githup.ahmedramadan.4cashier/internal/models"
	"githup.ahmedramadan.4cashier/internal/repos"
)

// BadgeEvaluator هو العقد الذي يجب أن تتبعه كل قاعدة إنجاز
type BadgeEvaluator interface {
	Evaluate(ctx context.Context, profile *models.HRProfile) (*models.Badge, error)
	Name() string
}

// =================================================================
// ⭐️ تطبيق قاعدة شارة: Top Rated HR
// =================================================================

type TopRatedBadgeEvaluator struct {
	repo repos.HRRepository
}

func NewTopRatedBadgeEvaluator(repo repos.HRRepository) BadgeEvaluator {
	return &TopRatedBadgeEvaluator{repo: repo}
}

func (b *TopRatedBadgeEvaluator) Name() string {
	return "Top Rated HR"
}

func (b *TopRatedBadgeEvaluator) Evaluate(ctx context.Context, profile *models.HRProfile) (*models.Badge, error) {
	// قاعدة العمل: 50 تقييم على الأقل ومتوسط 4.5 أو أعلى
	if profile.TotalRatesCount >= 50 && *profile.Rate >= 4.5 {
		
		// تأكد من عدم تكرار الشارة (باستخدام الـ Repository)
		hasBadge, err := b.repo.CheckIfProfileHasBadge(ctx, profile.ID, b.Name())
		if err != nil {
			return nil, err
		}
		
		if !hasBadge {
			// بناء الشارة
			badge := &models.Badge{
				HRProfileID:     profile.ID,
				CreatedDate:     time.Now(),
				TotalRates:      profile.TotalRatesCount,
				Rate:            *profile.Rate,
				JobPosition:     b.Name(), 
				CurrentJobRoles: fmt.Sprintf("Achieved %d rates with %.2f avg", profile.TotalRatesCount, profile.Rate),
			}
			return badge, nil
		}
	}
	
	return nil, nil
}