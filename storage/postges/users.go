package postgres

import (
	"context"
	"fmt"
	"morc/models"
	"morc/pkg/helper"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func (u *userRepo) Create(ctx context.Context, req *models.CreateUser) (*models.UserPrimaryKey, error) {

	query := `INSERT INTO users (
		first_name,
		last_name,
		username,
		phone,
		telegram_id,
		role,
		is_verified,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, false, now(), now())`

	_, err := u.db.Exec(ctx, query,
		req.FirstName,
		req.LastName,
		req.UserName,
		req.PhoneNumber,
		req.TelegramId,
		req.Role,
	)

	pKey := &models.UserPrimaryKey{Id: "1"}
	return pKey, err
}

func (u *userRepo) GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {
	res := &models.User{}
	query := `SELECT
		id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at
		FROM users WHERE id = $1`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.UserName, &res.PhoneNumber, &res.TelegramId, &res.Role, &res.IsVerified, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *userRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (*models.GetListUserResponse, error) {
	res := &models.GetListUserResponse{}
	params := make(map[string]interface{})
	var arr []interface{}

	// Agar Limit 0 bo'lsa, default 10 qilib qo'yamiz
	if req.Limit == 0 {
		req.Limit = 10
	}

	fmt.Println("DEBUG | Limit:", req.Limit, "Offset:", req.Offset)

	query := `SELECT id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at FROM users`
	filter := " WHERE 1=1"
	order := " ORDER BY created_at DESC"
	offset := fmt.Sprintf(" OFFSET %d", req.Offset)
	limit := fmt.Sprintf(" LIMIT %d", req.Limit)

	// Qidiruv sharti
	if len(req.Search) > 0 {
		params["search"] = "%" + req.Search + "%"
		filter += " AND (LOWER(first_name) ILIKE :search OR LOWER(last_name) ILIKE :search OR phone ILIKE :search OR telegram_id::TEXT ILIKE :search OR LOWER(role) ILIKE :search)"
	}

	// Debug: Queryni chiqarish
	cQ := `SELECT count(1) FROM users` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	fmt.Println("DEBUG | Count Query:", cQ)
	
	err := u.db.QueryRow(ctx, cQ, arr...).Scan(&res.Count)
	if err != nil {
		fmt.Println("❌ User count olishda xatolik:", err)
		return res, err
	}

	q := query + filter + order + offset + limit
	q, arr = helper.ReplaceQueryParams(q, params)
	fmt.Println("DEBUG | Final Query:", q)

	rows, err := u.db.Query(ctx, q, arr...)
	if err != nil {
		fmt.Println("❌ User ma'lumotlarini olishda xatolik:", err)
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		obj := &models.User{}
		err = rows.Scan(
			&obj.Id, &obj.FirstName, &obj.LastName, &obj.UserName,
			&obj.PhoneNumber, &obj.TelegramId, &obj.Role, &obj.IsVerified,
			&obj.CreatedAt, &obj.UpdatedAt,
		)
		if err != nil {
			fmt.Println("❌ User ma'lumotlarini o'qishda xatolik:", err)
			return res, err
		}
		res.Users = append(res.Users, obj)
	}

	fmt.Println("✅ Userlar muvaffaqiyatli olindi:", res)
	return res, nil
}


func (u *userRepo) Update(ctx context.Context, req *models.UpdateUser) (int64, error) {
	query := `UPDATE users SET `
	params := []interface{}{}
	counter := 1
	updated := false

	if req.FirstName != "" {
		query += fmt.Sprintf("first_name = $%d, ", counter)
		params = append(params, req.FirstName)
		counter++
		updated = true
	}
	if req.LastName != "" {
		query += fmt.Sprintf("last_name = $%d, ", counter)
		params = append(params, req.LastName)
		counter++
		updated = true
	}
	if req.PhoneNumber != "" {
		query += fmt.Sprintf("phone = $%d, ", counter)
		params = append(params, req.PhoneNumber)
		counter++
		updated = true
	}
	if req.TelegramId != "" {
		query += fmt.Sprintf("telegram_id = $%d, ", counter)
		params = append(params, req.TelegramId)
		counter++
		updated = true
	}
	if req.Role != "" {
		query += fmt.Sprintf("role = $%d, ", counter)
		params = append(params, req.Role)
		counter++
		updated = true
	}
	if updated {
		query = query[:len(query)-2] + fmt.Sprintf(", updated_at = now() WHERE id = $%d", counter)
		params = append(params, req.Id)
		result, err := u.db.Exec(ctx, query, params...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected(), nil
	}
	return 0, fmt.Errorf("hech qanday yangilanish kiritilmagan")
}

func (u *userRepo) Delete(ctx context.Context, req *models.UserPrimaryKey) (int64, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (u *userRepo) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	res := &models.User{}
	query := `SELECT id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at FROM users WHERE phone = $1`
	err := u.db.QueryRow(ctx, query, phone).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserName, &res.PhoneNumber, &res.TelegramId, &res.Role, &res.IsVerified, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (u *userRepo) GetUserByTelegramID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {
	res := &models.User{}
	query := `SELECT
		id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at
		FROM users WHERE telegram_id = $1`

	err := u.db.QueryRow(ctx, query, req.TelegramId).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.UserName, &res.PhoneNumber, &res.TelegramId, &res.Role, &res.IsVerified, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetByRole - ma'lum bir rolga ega foydalanuvchilarni olish
func (s *userRepo) GetByRole(ctx context.Context, role string) ([]models.User, error) {
	var users []models.User
	query := `SELECT id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at FROM users WHERE role = $1`
	rows, err := s.db.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.UserName,&user.PhoneNumber,&user.TelegramId,&user.Role,&user.IsVerified,&user.CreatedAt,&user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUnconfirmedCouriers - tasdiqlanmagan kuryerlarni olish
func (s *userRepo) GetUnconfirmedCouriers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT
		id, first_name, last_name, username, phone, telegram_id, role, is_verified, created_at, updated_at
		FROM users WHERE role != 'admin' AND is_verified = FALSE`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName,&user.UserName, &user.PhoneNumber,&user.TelegramId,&user.Role,&user.IsVerified,&user.CreatedAt,&user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


func (s *userRepo) Approve(ctx context.Context, userID string) error {
    query := `UPDATE users SET is_verified = TRUE WHERE id = $1`
    _, err := s.db.Exec(ctx, query, userID)
    return err
}

func (s *userRepo) Reject(ctx context.Context, userID string) error {
    query := `UPDATE users SET is_verified = FALSE WHERE id = $1`
    _, err := s.db.Exec(ctx, query, userID)
    return err
}



