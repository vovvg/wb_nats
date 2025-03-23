package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
	"wb_nats/internal/schema"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) InsertMessage(message schema.Request) error {
	ctx := context.Background()

	var deliveryID int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO delivery (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, message.Delivery.Name, message.Delivery.Phone, message.Delivery.Zip, message.Delivery.City,
		message.Delivery.Address, message.Delivery.Region, message.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return fmt.Errorf("error add delivery: %w", err)
	}

	paymentTime := time.Unix(int64(message.Payment.PaymentDt), 0).UTC()
	var paymentID int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, message.Payment.Transaction, message.Payment.RequestId, message.Payment.Currency,
		message.Payment.Provider, message.Payment.Amount, paymentTime, message.Payment.Bank,
		message.Payment.DeliveryCost, message.Payment.GoodsTotal, message.Payment.CustomFee).Scan(&paymentID)
	if err != nil {
		return fmt.Errorf("error add payment: %w", err)
	}

	var requestID int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO request (order_uid, track_number, entry, locale, internal_signature, customer_id, 
		                    delivery_service, shard_key, sm_id, date_created, oof_shard, delivery_id, payment_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`, message.OrderUid, message.TrackNumber, message.Entry, message.Locale, message.InternalSignature,
		message.CustomerId, message.DeliveryService, message.ShardKey, message.SmId, message.DateCreated,
		message.OofShard, deliveryID, paymentID).Scan(&requestID)
	if err != nil {
		return fmt.Errorf("error add request: %w", err)
	}

	for _, item := range message.Items {
		_, err = s.pool.Exec(ctx, `
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, request_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmId, item.Brand, item.Status, requestID)
		if err != nil {
			return fmt.Errorf("error add item: %w", err)
		}
	}

	log.Printf("Success, request ID: %d\n", requestID)
	return nil
}

func (s *Storage) GetMessage(orderId string) (schema.Request, error) {
	var request schema.Request
	var payment schema.Payment
	var delivery schema.Delivery
	var paymentDt time.Time

	err := s.pool.QueryRow(context.Background(), `
        SELECT r.id, r.order_uid, r.track_number, r.entry, r.locale, r.internal_signature, 
               r.customer_id, r.delivery_service, r.shard_key, r.sm_id, r.date_created, 
               r.oof_shard, r.delivery_id, r.payment_id,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
               p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email
        FROM request r
        LEFT JOIN payment p ON r.payment_id = p.id
        LEFT JOIN delivery d ON r.delivery_id = d.id
        WHERE r.order_uid = $1
    `, orderId).Scan(
		&request.Id, &request.OrderUid, &request.TrackNumber, &request.Entry, &request.Locale,
		&request.InternalSignature, &request.CustomerId, &request.DeliveryService, &request.ShardKey,
		&request.SmId, &request.DateCreated, &request.OofShard, &request.DeliveryId, &request.PaymentId,
		&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount,
		&paymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
		&delivery.Region, &delivery.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return request, fmt.Errorf("order with ID %s not found", orderId)
		}
		return request, fmt.Errorf("error fetching order details for order %s: %w", orderId, err)
	}

	request.Payment = payment
	request.Delivery = delivery
	request.Payment.PaymentDt = int(paymentDt.Unix())

	rows, err := s.pool.Query(context.Background(), `
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items
        WHERE request_id = $1
    `, request.Id)

	if err != nil {
		return request, fmt.Errorf("error fetching items for order %s: %w", orderId, err)
	}
	defer rows.Close()

	var items []schema.Item
	for rows.Next() {
		var item schema.Item
		err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return request, fmt.Errorf("error scanning item for order %s: %w", orderId, err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return request, fmt.Errorf("error iterating over items for order %s: %w", orderId, err)
	}

	request.Items = items

	return request, nil
}
