package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/QwaQ-dev/servicesSubscription/internal/structures"
	"github.com/QwaQ-dev/servicesSubscription/pkg/sl"
)

type SubscriptionRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewSubsriptionRepo(
	db *sql.DB,
	log *slog.Logger,
) *SubscriptionRepo {
	return &SubscriptionRepo{
		db:  db,
		log: log,
	}
}

func (r *SubscriptionRepo) InsertSub(subscription *structures.Subscription) (int, error) {
	const op = "repository.subscriptionRepo.InsertSub"
	log := r.log.With("op", op)

	log.Info("Inserting subscription", slog.Any("subscription", subscription))

	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING ID
	`

	var id int

	err := r.db.QueryRow(
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(&id)

	if err != nil {
		log.Error("Failed to insert sub", sl.Err(err))
		return 0, err
	}

	log.Debug("Inserted successfully", slog.Int("id", id))
	return id, nil
}

func (r *SubscriptionRepo) SelectAllSubs() ([]structures.Subscription, error) {
	const op = "repository.subscriptionRepo.SelectAllSubs"
	log := r.log.With("op", op)

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Error("Failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	var subscriptions []structures.Subscription

	for rows.Next() {
		var subscription structures.Subscription

		err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
		)
		if err != nil {
			log.Error("Failed to fetch subscriptions", sl.Err(err))
			continue
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		log.Error("Rows iteration error", sl.Err(err))
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepo) SelectSubById(id int) (structures.Subscription, error) {
	const op = "repository.subscriptionRepo.SelectSubById"
	log := r.log.With("op", op)

	var subscription structures.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		log.Error("Failed to select sub", sl.Err(err))
		return subscription, fmt.Errorf("%s:%v", op, err)
	}

	return subscription, nil
}

func (r *SubscriptionRepo) UpdateSub(subscription *structures.Subscription, id int) error {
	const op = "repository.subscriptionsRepo.UpdateSub"
	log := r.log.With("op", op)

	query := `
		UPDATE subscriptions
		SET service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		id,
	)

	if err != nil {
		log.Error("Failed to update sub", sl.Err(err))
		return fmt.Errorf("%s:%v", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get affected rows", sl.Err(err))
		return fmt.Errorf("%s:%v", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no subs with id:%d", op, subscription.ID)
	}

	log.Info("Subscription updated", slog.Int("id", subscription.ID))
	return nil
}

func (r *SubscriptionRepo) DeleteSub(id int) error {
	const op = "repository.subscriptionsRepo.DeleteSub"
	log := r.log.With("op", op)

	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("Failed to delete sub", sl.Err(err))
		return fmt.Errorf("%s:%v", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get affected rows", sl.Err(err))
		return fmt.Errorf("%s:%v", op, err)
	}

	if rowsAffected == 0 {
		log.Info("No subscription found with ID", slog.Int("id", id))
	} else {
		log.Info("Subscription deleted", slog.Int("id", id))
	}

	return nil
}

func (r *SubscriptionRepo) SelectSum(data *structures.Counting) (int, error) {
	const op = "repository.subscriptionRepo.SelectSum"
	log := r.log.With("op", op)

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE to_date(start_date, 'MM-YYYY') 
			BETWEEN to_date($1, 'MM-YYYY') AND to_date($2, 'MM-YYYY')
		  AND ($3 = '' OR user_id = $3::uuid)
		  AND ($4 = '' OR service_name = $4)
	`

	var total int

	err := r.db.QueryRow(query, data.StartDate, data.EndDate, data.UserID, data.ServiceName).Scan(&total)
	if err != nil {
		log.Error("Failed to select sum", sl.Err(err))
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return total, nil
}
