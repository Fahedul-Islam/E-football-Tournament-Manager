package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

func (s *service) GetNotifications(ctx context.Context, userID int, page int) ([]*domain.Notification, error) {
	return s.announcementRepo.GetNotifications(ctx, userID, page)
}

func (s *service) MarkNotificationAsRead(ctx context.Context, notificationID int, userID int) error {
	return s.announcementRepo.MarkNotificationAsRead(ctx, notificationID, userID)
}

func (s *service) MarkAllNotificationsAsRead(ctx context.Context, userID int) error {
	return s.announcementRepo.MarkAllNotificationsAsRead(ctx, userID)
}