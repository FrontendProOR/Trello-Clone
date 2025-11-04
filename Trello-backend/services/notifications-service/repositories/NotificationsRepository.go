package repositories

import (
	"notifications-service/models"

	"github.com/gocql/gocql"
)

type NotificationsRepository struct {
	session *gocql.Session
}

func NewNotificationsRepository(session *gocql.Session) *NotificationsRepository {
	return &NotificationsRepository{session: session}
}

// Kreiranje tabele ako ne postoji
func (r *NotificationsRepository) CreateTable() error {
	// Kreiraj tabelu
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		user_id TEXT,
		id UUID,
		message TEXT,
		created_at TIMESTAMP,
		is_read BOOLEAN,
		PRIMARY KEY (user_id, id)
	) WITH CLUSTERING ORDER BY (id DESC);`
	if err := r.session.Query(query).Exec(); err != nil {
		return err
	}

	// Sekundarni indeks na is_read (da možeš ALLOW FILTERING da izbegneš u nekim slučajevima)
	query = `CREATE INDEX IF NOT EXISTS is_read_index ON notifications (is_read)`
	return r.session.Query(query).Exec()
}

// Preuzimanje svih notifikacija po userID
func (r *NotificationsRepository) GetNotificationsByUserID(userID string) ([]models.Notification, error) {
	var notifications []models.Notification

	query := `SELECT id, message, created_at, is_read 
			  FROM notifications 
			  WHERE user_id = ?`
	iter := r.session.Query(query, userID).Iter()

	var notification models.Notification
	for iter.Scan(&notification.ID, &notification.Message, &notification.CreatedAt, &notification.IsRead) {
		notification.UserID = userID
		notifications = append(notifications, notification)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Preuzimanje nepročitanih notifikacija
func (r *NotificationsRepository) GetUnreadNotifications(userID string) ([]models.Notification, error) {
	var notifications []models.Notification

	query := `SELECT id, message, created_at, is_read 
			  FROM notifications 
			  WHERE user_id = ? AND is_read = false`
	iter := r.session.Query(query, userID).Iter()

	var notification models.Notification
	for iter.Scan(&notification.ID, &notification.Message, &notification.CreatedAt, &notification.IsRead) {
		notification.UserID = userID
		notifications = append(notifications, notification)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Kreiranje notifikacije
func (r *NotificationsRepository) CreateNotification(notification *models.Notification) error {
	query := `INSERT INTO notifications (user_id, id, message, created_at, is_read)
			  VALUES (?, ?, ?, ?, ?)`
	return r.session.Query(query,
		notification.UserID,
		notification.ID,
		notification.Message,
		notification.CreatedAt,
		notification.IsRead).Exec()
}

// Označavanje notifikacije kao pročitane
func (r *NotificationsRepository) MarkAsRead(userID string, notificationID string) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = ? AND id = ?`
	return r.session.Query(query, userID, notificationID).Exec()
}

// Brisanje notifikacije
func (r *NotificationsRepository) DeleteNotification(userID, notificationID string) error {
	query := `DELETE FROM notifications WHERE user_id = ? AND id = ?`
	return r.session.Query(query, userID, notificationID).Exec()
}
