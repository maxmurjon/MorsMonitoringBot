package postgres

import (
	"context"
	"fmt"
	"morc/models"
	"morc/pkg/helper"

	"github.com/jackc/pgx/v5/pgxpool"
)

type barrelRepo struct {
	db *pgxpool.Pool
}

func (b *barrelRepo) Create(ctx context.Context, req *models.CreateBarrel) (*models.Barrel, error) {
	query := `INSERT INTO barrels (
		name,
		volume_liters,
		current_volume,
		location_name,
		latitude,
		longitude,
		assigned_seller_id,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now()) RETURNING id, created_at, updated_at`

	barrel := &models.Barrel{}
	err := b.db.QueryRow(ctx, query,
		req.Name,
		req.VolumeLiters,
		req.CurrentVolume,
		req.LocationName,
		req.Latitude,
		req.Longitude,
		req.AssignedSellerId,
	).Scan(&barrel.Id, &barrel.CreatedAt, &barrel.UpdatedAt)
	if err != nil {
		return nil, err
	}

	barrel.Name = req.Name
	barrel.VolumeLiters = req.VolumeLiters
	barrel.CurrentVolume = req.CurrentVolume
	barrel.LocationName = req.LocationName
	barrel.Latitude = req.Latitude
	barrel.Longitude = req.Longitude
	barrel.AssignedSellerId = req.AssignedSellerId

	return barrel, nil
}

func (b *barrelRepo) GetByID(ctx context.Context, id int64) (*models.Barrel, error) {
	barrel := &models.Barrel{}
	query := `SELECT
		id, name, volume_liters, current_volume, location_name, latitude, longitude, assigned_seller_id, created_at, updated_at
		FROM barrels WHERE id = $1`

	err := b.db.QueryRow(ctx, query, id).Scan(
		&barrel.Id, &barrel.Name, &barrel.VolumeLiters, &barrel.CurrentVolume, &barrel.LocationName, &barrel.Latitude, &barrel.Longitude, &barrel.AssignedSellerId, &barrel.CreatedAt, &barrel.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return barrel, nil
}

func (b *barrelRepo) GetList(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error) {
	res := &models.GetListBarrelResponse{}
	params := make(map[string]interface{})
	var arr []interface{}

	if req.Limit == 0 {
		req.Limit = 10
	}

	query := `SELECT id, name, volume_liters, current_volume, location_name, latitude, longitude, assigned_seller_id, created_at, updated_at FROM barrels`
	filter := " WHERE 1=1"
	order := " ORDER BY created_at DESC"
	offset := fmt.Sprintf(" OFFSET %d", req.Offset)
	limit := fmt.Sprintf(" LIMIT %d", req.Limit)

	if len(req.Search) > 0 {
		params["search"] = "%" + req.Search + "%"
		filter += " AND (LOWER(name) ILIKE :search OR LOWER(location_name) ILIKE :search)"
	}

	cQ := `SELECT count(1) FROM barrels` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err := b.db.QueryRow(ctx, cQ, arr...).Scan(&res.Count)
	if err != nil {
		return res, err
	}

	q := query + filter + order + offset + limit
	q, arr = helper.ReplaceQueryParams(q, params)

	rows, err := b.db.Query(ctx, q, arr...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		barrel := &models.Barrel{}
		err = rows.Scan(
			&barrel.Id, &barrel.Name, &barrel.VolumeLiters, &barrel.CurrentVolume,
			&barrel.LocationName, &barrel.Latitude, &barrel.Longitude, &barrel.AssignedSellerId,
			&barrel.CreatedAt, &barrel.UpdatedAt,
		)
		if err != nil {
			return res, err
		}
		res.Barrels = append(res.Barrels, barrel)
	}

	return res, nil
}

func (b *barrelRepo) Update(ctx context.Context, req *models.UpdateBarrel) (int64, error) {
	query := `UPDATE barrels SET `
	params := []interface{}{}
	counter := 1
	updated := false

	if req.Name != "" {
		query += fmt.Sprintf("name = $%d, ", counter)
		params = append(params, req.Name)
		counter++
		updated = true
	}
	if req.VolumeLiters != 0 {
		query += fmt.Sprintf("volume_liters = $%d, ", counter)
		params = append(params, req.VolumeLiters)
		counter++
		updated = true
	}
	if req.CurrentVolume != 0 {
		query += fmt.Sprintf("current_volume = $%d, ", counter)
		params = append(params, req.CurrentVolume)
		counter++
		updated = true
	}
	if req.LocationName != "" {
		query += fmt.Sprintf("location_name = $%d, ", counter)
		params = append(params, req.LocationName)
		counter++
		updated = true
	}
	if req.Latitude != 0 {
		query += fmt.Sprintf("latitude = $%d, ", counter)
		params = append(params, req.Latitude)
		counter++
		updated = true
	}
	if req.Longitude != 0 {
		query += fmt.Sprintf("longitude = $%d, ", counter)
		params = append(params, req.Longitude)
		counter++
		updated = true
	}
	if req.AssignedSellerId != nil {
		query += fmt.Sprintf("assigned_seller_id = $%d, ", counter)
		params = append(params, req.AssignedSellerId)
		counter++
		updated = true
	}
	if updated {
		query = query[:len(query)-2] + fmt.Sprintf(", updated_at = now() WHERE id = $%d", counter)
		params = append(params, req.Id)
		result, err := b.db.Exec(ctx, query, params...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected(), nil
	}
	return 0, fmt.Errorf("no updates provided")
}

func (b *barrelRepo) Delete(ctx context.Context, id int64) (int64, error) {
	query := `DELETE FROM barrels WHERE id = $1`
	result, err := b.db.Exec(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (b *barrelRepo) GetListSellerId(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error) {
	res := &models.GetListBarrelResponse{}
	params := make(map[string]interface{})
	var arr []interface{}

	if req.Limit == 0 {
		req.Limit = 10
	}

	// Asosiy so‘rov
	query := `SELECT id, name, volume_liters, current_volume, location_name, latitude, longitude, assigned_seller_id, created_at, updated_at FROM barrels`

	// Bo‘sh bochkalarni filtrlash (NULL yoki bo‘sh qator)
	filter := " WHERE (assigned_seller_id = 0 OR assigned_seller_id IS NULL)"
	order := " ORDER BY created_at DESC"
	offset := fmt.Sprintf(" OFFSET %d", req.Offset)
	limit := fmt.Sprintf(" LIMIT %d", req.Limit)

	// Agar qidiruv so‘rovi berilgan bo‘lsa
	if len(req.Search) > 0 {
		params["search"] = "%" + req.Search + "%"
		filter += " AND (LOWER(name) ILIKE :search OR LOWER(location_name) ILIKE :search)"
	}

	// Umumiy sonni olish uchun so‘rov
	cQ := `SELECT count(1) FROM barrels` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)

	err := b.db.QueryRow(ctx, cQ, arr...).Scan(&res.Count)
	if err != nil {
		return res, err
	}

	// Ma'lumotlarni olish uchun asosiy so‘rov
	q := query + filter + order + offset + limit
	q, arr = helper.ReplaceQueryParams(q, params)

	rows, err := b.db.Query(ctx, q, arr...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// Ma'lumotlarni o‘qish
	for rows.Next() {
		barrel := &models.Barrel{}
		err = rows.Scan(
			&barrel.Id, &barrel.Name, &barrel.VolumeLiters, &barrel.CurrentVolume,
			&barrel.LocationName, &barrel.Latitude, &barrel.Longitude, &barrel.AssignedSellerId,
			&barrel.CreatedAt, &barrel.UpdatedAt,
		)
		if err != nil {
			return res, err
		}
		res.Barrels = append(res.Barrels, barrel)
	}

	return res, nil
}


func (b *barrelRepo) GetBarrelBySellerId(ctx context.Context, sellerId string) (*models.Barrel, error) {
	barrel := &models.Barrel{}
	query := `SELECT
		id, name, volume_liters, current_volume, location_name, latitude, longitude, assigned_seller_id, created_at, updated_at
		FROM barrels WHERE assigned_seller_id = $1`

	err := b.db.QueryRow(ctx, query, sellerId).Scan(
		&barrel.Id, &barrel.Name, &barrel.VolumeLiters, &barrel.CurrentVolume, &barrel.LocationName, &barrel.Latitude, &barrel.Longitude, &barrel.AssignedSellerId, &barrel.CreatedAt, &barrel.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return barrel, nil
}

func (b *barrelRepo) GetListEmpty(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error) {
	res := &models.GetListBarrelResponse{}
	params := make(map[string]interface{})
	var arr []interface{}

	if req.Limit == 0 {
		req.Limit = 10
	}

	query := `SELECT id, name, volume_liters, current_volume, location_name, latitude, longitude, assigned_seller_id, created_at, updated_at FROM barrels WHERE current_volume = 0`
	order := " ORDER BY created_at DESC"
	offset := fmt.Sprintf(" OFFSET %d", req.Offset)
	limit := fmt.Sprintf(" LIMIT %d", req.Limit)

	if len(req.Search) > 0 {
		params["search"] = "%" + req.Search + "%"
		query += " AND (LOWER(name) ILIKE :search OR LOWER(location_name) ILIKE :search)"
	}

	cQ := `SELECT count(1) FROM barrels WHERE current_volume = 0`
	if len(req.Search) > 0 {
		cQ += " AND (LOWER(name) ILIKE :search OR LOWER(location_name) ILIKE :search)"
	}
	cQ, arr = helper.ReplaceQueryParams(cQ, params)

	err := b.db.QueryRow(ctx, cQ, arr...).Scan(&res.Count)
	if err != nil {
		return res, err
	}

	q := query + order + offset + limit
	q, arr = helper.ReplaceQueryParams(q, params)

	rows, err := b.db.Query(ctx, q, arr...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		barrel := &models.Barrel{}
		err = rows.Scan(
			&barrel.Id, &barrel.Name, &barrel.VolumeLiters, &barrel.CurrentVolume,
			&barrel.LocationName, &barrel.Latitude, &barrel.Longitude, &barrel.AssignedSellerId,
			&barrel.CreatedAt, &barrel.UpdatedAt,
		)
		if err != nil {
			return res, err
		}
		res.Barrels = append(res.Barrels, barrel)
	}

	return res, nil
}
