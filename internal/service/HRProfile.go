package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"githup.ahmedramadan.4cashier/internal/bootstrap"
	"githup.ahmedramadan.4cashier/internal/models"
	"githup.ahmedramadan.4cashier/internal/repos"
)

type HRService struct {
	log  zerolog.Logger
	repo repos.HRRepository
	badgeEvaluators []BadgeEvaluator
}

func NewHRService(log zerolog.Logger, repo repos.HRRepository) *HRService {
	return &HRService{
		log:  log.With().Str("layer", "service").Str("component", "HRService").Logger(),
		repo: repo,
	}
}



func (s *HRService) AddExperience(ctx context.Context, hrID int, exp []models.Experience) error {
	err := s.repo.AddExperience(ctx, hrID, exp)
	if err != nil {
		s.log.Error().Err(err).Msg("AddExperience failed")
	}
	return err
}

func (s *HRService) AddJobRoles(ctx context.Context, hrID int, roles []models.JobRole) error {
	err := s.repo.AddJobRoles(ctx, hrID, roles)
	if err != nil {
		s.log.Error().Err(err).Msg("AddJobRoles failed")
	}
	return err
}

func (s *HRService) GetEmployeeStats(ctx context.Context, employeeId int) (models.EmployeeStats, error) {
	return s.repo.GetEmployeeStats(ctx,employeeId)
	
}


func (s *HRService) LikeBadge(ctx context.Context, like *models.BadgeLike) error {
	err := s.repo.LikeBadge(ctx, like)
	if err != nil {
		s.log.Error().Err(err).Msg("LikeBadge failed")
	}
	return err
}

func (s *HRService) AwardBadge(ctx context.Context, badge *models.Badge) (int, error) {
	id, err := s.repo.AwardBadge(ctx, badge)
	if err != nil {
		s.log.Error().Err(err).Msg("AwardBadge failed")
	}
	return id, err
}



func (s * HRService) RateHR(ctx context.Context, rate *models.Rate) (*models.HRProfile, error) {
	
	// 1. تسجيل التقييم وتحديث المتوسط (Atomic)
	_, _, err := s.repo.RateHR(ctx, rate)
	if err != nil {
		return nil, fmt.Errorf("service failed to execute rate transaction: %w", err)
	}

	// 2. جلب البروفايل المحدث (لتقييم الشارات)
	profile, err := s.repo.GetHRProfileByID(ctx, rate.HRProfileID)
	if err != nil {
		return nil, fmt.Errorf("service failed to fetch updated profile: %w", err)
	}
    
	// 3. تشغيل محرك الشارات (The Engine)
	awardedBadges, err := s.runBadgeEngine(ctx, profile)
	if err != nil {
		fmt.Printf("Warning: Failed to run badge engine for profile %d: %v\n", profile.ID, err)
	}
	
	// (اختياري) إضافة الشارات الممنوحة حديثاً لنموذج الرد
	profile.Badges = awardedBadges 

	return profile, nil
}

func (s * HRService) runBadgeEngine(ctx context.Context, profile *models.HRProfile) ([]models.Badge, error) {
	var awardedBadges []models.Badge
	
	for _, evaluator := range s.badgeEvaluators {
		badge, err := evaluator.Evaluate(ctx, profile)
		
		if err != nil {
			return nil, err 
		}
		
		if badge != nil {
			// منح الشارة
			_, err := s.repo.AwardBadge(ctx, badge)
			if err == nil {
				awardedBadges = append(awardedBadges, *badge)
			}
		}
	}
	
	return awardedBadges, nil
}

func (s * HRService) LikeRate(ctx context.Context, like *models.RateLike) (int, error) {
	
	// 1. تسجيل اللايك وتحديث الـ count (Atomic)
	newLikesCount, err := s.repo.LikeRate(ctx, like)
	if err != nil {
		return 0, fmt.Errorf("service failed to execute like transaction: %w", err)
	}

	// 2. المصداقية للمُصوّت (Voter Credibility)
	if like.IsLike {
		// +1 نقطة للموظف الذي يساهم في الفلترة
		s.repo.UpdateEmployeePoints(ctx, like.EmployeeID, 1) 
	}

	// 3. المصداقية للكاتب (Author Credibility)
	// قاعدة العمل: عند الوصول لـ 50 لايك، الكاتب يأخذ نقاط إضافية
	if newLikesCount == 50 { 
		rateOwnerID, err := s.repo.GetRateOwner(ctx, like.RateID)
		if err == nil && rateOwnerID > 0 {
			// +10 نقاط إضافية لمساهمته بمحتوى موثوق
			s.repo.UpdateEmployeePoints(ctx, rateOwnerID, 10) 
		}
	}
	
	return newLikesCount, nil
}


func (s * HRService) GetHRProfiles(ctx context.Context, pagination bootstrap.Pagination, filters map[string]interface{}) ([]models.HRProfile, error) {
	return s.repo.GetHRProfiles(ctx, pagination, filters)
}

func (s * HRService) GetRates(ctx context.Context, pagination bootstrap.Pagination, filters map[string]interface{}) ([]models.RateWithDetails, error) {
	return s.repo.GetRates(ctx, pagination, filters)
}