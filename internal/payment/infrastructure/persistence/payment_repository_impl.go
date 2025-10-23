package persistence

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"obs-tools-usage/internal/payment/domain/entity"
	"obs-tools-usage/internal/payment/domain/repository"
)

// PaymentRepositoryImpl implements PaymentRepository interface using MariaDB
type PaymentRepositoryImpl struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewPaymentRepositoryImpl creates a new payment repository implementation
func NewPaymentRepositoryImpl(db *gorm.DB, logger *logrus.Logger) repository.PaymentRepository {
	return &PaymentRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// CreatePayment creates a new payment
func (r *PaymentRepositoryImpl) CreatePayment(payment *entity.Payment) error {
	r.logger.WithField("payment_id", payment.ID).Debug("Creating payment in database")

	if err := r.db.Create(payment).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", payment.ID).Error("Failed to create payment")
		return fmt.Errorf("failed to create payment: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"payment_id": payment.ID,
		"user_id":    payment.UserID,
		"amount":     payment.Amount,
		"status":     payment.Status,
	}).Debug("Successfully created payment")

	return nil
}

// GetPayment retrieves a payment by ID
func (r *PaymentRepositoryImpl) GetPayment(paymentID string) (*entity.Payment, error) {
	r.logger.WithField("payment_id", paymentID).Debug("Getting payment from database")

	var payment entity.Payment
	if err := r.db.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment not found: %s", paymentID)
		}
		r.logger.WithError(err).WithField("payment_id", paymentID).Error("Failed to get payment")
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	r.logger.WithField("payment_id", paymentID).Debug("Successfully retrieved payment")
	return &payment, nil
}

// UpdatePayment updates an existing payment
func (r *PaymentRepositoryImpl) UpdatePayment(payment *entity.Payment) error {
	r.logger.WithField("payment_id", payment.ID).Debug("Updating payment in database")

	payment.UpdatedAt = time.Now()
	if err := r.db.Save(payment).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", payment.ID).Error("Failed to update payment")
		return fmt.Errorf("failed to update payment: %w", err)
	}

	r.logger.WithField("payment_id", payment.ID).Debug("Successfully updated payment")
	return nil
}

// DeletePayment deletes a payment
func (r *PaymentRepositoryImpl) DeletePayment(paymentID string) error {
	r.logger.WithField("payment_id", paymentID).Debug("Deleting payment from database")

	if err := r.db.Where("id = ?", paymentID).Delete(&entity.Payment{}).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", paymentID).Error("Failed to delete payment")
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	r.logger.WithField("payment_id", paymentID).Debug("Successfully deleted payment")
	return nil
}

// GetPaymentsByUser retrieves payments by user ID
func (r *PaymentRepositoryImpl) GetPaymentsByUser(userID string) ([]*entity.Payment, error) {
	r.logger.WithField("user_id", userID).Debug("Getting payments by user from database")

	var payments []*entity.Payment
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error; err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to get payments by user")
		return nil, fmt.Errorf("failed to get payments by user: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"user_id":        userID,
		"payments_count": len(payments),
	}).Debug("Successfully retrieved payments by user")

	return payments, nil
}

// GetPaymentsByBasket retrieves payments by basket ID
func (r *PaymentRepositoryImpl) GetPaymentsByBasket(basketID string) ([]*entity.Payment, error) {
	r.logger.WithField("basket_id", basketID).Debug("Getting payments by basket from database")

	var payments []*entity.Payment
	if err := r.db.Where("basket_id = ?", basketID).Order("created_at DESC").Find(&payments).Error; err != nil {
		r.logger.WithError(err).WithField("basket_id", basketID).Error("Failed to get payments by basket")
		return nil, fmt.Errorf("failed to get payments by basket: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"basket_id":      basketID,
		"payments_count": len(payments),
	}).Debug("Successfully retrieved payments by basket")

	return payments, nil
}

// GetPaymentsByStatus retrieves payments by status
func (r *PaymentRepositoryImpl) GetPaymentsByStatus(status entity.PaymentStatus) ([]*entity.Payment, error) {
	r.logger.WithField("status", status).Debug("Getting payments by status from database")

	var payments []*entity.Payment
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&payments).Error; err != nil {
		r.logger.WithError(err).WithField("status", status).Error("Failed to get payments by status")
		return nil, fmt.Errorf("failed to get payments by status: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"status":         status,
		"payments_count": len(payments),
	}).Debug("Successfully retrieved payments by status")

	return payments, nil
}

// GetPaymentsByDateRange retrieves payments within a date range
func (r *PaymentRepositoryImpl) GetPaymentsByDateRange(startDate, endDate string) ([]*entity.Payment, error) {
	r.logger.WithFields(logrus.Fields{
		"start_date": startDate,
		"end_date":   endDate,
	}).Debug("Getting payments by date range from database")

	var payments []*entity.Payment
	if err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Order("created_at DESC").Find(&payments).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payments by date range")
		return nil, fmt.Errorf("failed to get payments by date range: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"start_date":     startDate,
		"end_date":       endDate,
		"payments_count": len(payments),
	}).Debug("Successfully retrieved payments by date range")

	return payments, nil
}

// CreatePaymentItem creates a payment item
func (r *PaymentRepositoryImpl) CreatePaymentItem(item *entity.PaymentItem) error {
	r.logger.WithField("payment_id", item.PaymentID).Debug("Creating payment item in database")

	if err := r.db.Create(item).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", item.PaymentID).Error("Failed to create payment item")
		return fmt.Errorf("failed to create payment item: %w", err)
	}

	r.logger.WithField("payment_id", item.PaymentID).Debug("Successfully created payment item")
	return nil
}

// GetPaymentItems retrieves payment items by payment ID
func (r *PaymentRepositoryImpl) GetPaymentItems(paymentID string) ([]*entity.PaymentItem, error) {
	r.logger.WithField("payment_id", paymentID).Debug("Getting payment items from database")

	var items []*entity.PaymentItem
	if err := r.db.Where("payment_id = ?", paymentID).Find(&items).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", paymentID).Error("Failed to get payment items")
		return nil, fmt.Errorf("failed to get payment items: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"payment_id":    paymentID,
		"items_count":   len(items),
	}).Debug("Successfully retrieved payment items")

	return items, nil
}

// DeletePaymentItems deletes payment items by payment ID
func (r *PaymentRepositoryImpl) DeletePaymentItems(paymentID string) error {
	r.logger.WithField("payment_id", paymentID).Debug("Deleting payment items from database")

	if err := r.db.Where("payment_id = ?", paymentID).Delete(&entity.PaymentItem{}).Error; err != nil {
		r.logger.WithError(err).WithField("payment_id", paymentID).Error("Failed to delete payment items")
		return fmt.Errorf("failed to delete payment items: %w", err)
	}

	r.logger.WithField("payment_id", paymentID).Debug("Successfully deleted payment items")
	return nil
}

// GetPaymentStats retrieves payment statistics for a user
func (r *PaymentRepositoryImpl) GetPaymentStats(userID string) (*repository.PaymentStats, error) {
	r.logger.WithField("user_id", userID).Debug("Getting payment stats from database")

	var stats repository.PaymentStats

	// Get total payments count
	if err := r.db.Model(&entity.Payment{}).Where("user_id = ?", userID).Count(&stats.TotalPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to get total payments count: %w", err)
	}

	// Get total amount
	if err := r.db.Model(&entity.Payment{}).Where("user_id = ?", userID).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total amount: %w", err)
	}

	// Get completed payments count
	if err := r.db.Model(&entity.Payment{}).Where("user_id = ? AND status = ?", userID, entity.PaymentStatusCompleted).Count(&stats.CompletedPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed payments count: %w", err)
	}

	// Get failed payments count
	if err := r.db.Model(&entity.Payment{}).Where("user_id = ? AND status = ?", userID, entity.PaymentStatusFailed).Count(&stats.FailedPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to get failed payments count: %w", err)
	}

	// Get pending payments count
	if err := r.db.Model(&entity.Payment{}).Where("user_id = ? AND status = ?", userID, entity.PaymentStatusPending).Count(&stats.PendingPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending payments count: %w", err)
	}

	// Calculate average amount
	if stats.TotalPayments > 0 {
		stats.AverageAmount = stats.TotalAmount / float64(stats.TotalPayments)
	}

	r.logger.WithFields(logrus.Fields{
		"user_id":           userID,
		"total_payments":    stats.TotalPayments,
		"total_amount":      stats.TotalAmount,
		"completed_payments": stats.CompletedPayments,
	}).Debug("Successfully retrieved payment stats")

	return &stats, nil
}

// GetTotalRevenue retrieves total revenue within a date range
func (r *PaymentRepositoryImpl) GetTotalRevenue(startDate, endDate string) (float64, error) {
	r.logger.WithFields(logrus.Fields{
		"start_date": startDate,
		"end_date":   endDate,
	}).Debug("Getting total revenue from database")

	var totalRevenue float64
	if err := r.db.Model(&entity.Payment{}).Where("status = ? AND created_at BETWEEN ? AND ?", entity.PaymentStatusCompleted, startDate, endDate).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get total revenue")
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"start_date":     startDate,
		"end_date":       endDate,
		"total_revenue":  totalRevenue,
	}).Debug("Successfully retrieved total revenue")

	return totalRevenue, nil
}

// GetPaymentCountByStatus retrieves payment count by status
func (r *PaymentRepositoryImpl) GetPaymentCountByStatus(status entity.PaymentStatus) (int64, error) {
	r.logger.WithField("status", status).Debug("Getting payment count by status from database")

	var count int64
	if err := r.db.Model(&entity.Payment{}).Where("status = ?", status).Count(&count).Error; err != nil {
		r.logger.WithError(err).WithField("status", status).Error("Failed to get payment count by status")
		return 0, fmt.Errorf("failed to get payment count by status: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"status": status,
		"count":  count,
	}).Debug("Successfully retrieved payment count by status")

	return count, nil
}

// Ping checks database connectivity
func (r *PaymentRepositoryImpl) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Ping()
}
