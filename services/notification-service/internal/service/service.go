package service

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendEmail(from string, to []string, data []byte) error {
	// ToDo: business logic

	return nil
}
