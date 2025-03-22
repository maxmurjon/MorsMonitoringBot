package storage

import (
	"context"
	"morc/models"
	"time"
)

type StorageRepoI interface {
	User() UserRepoI
	Barrel() BarrelRepository
	Redis() RedisRepoI
	// Sale() SaleRepository
	// CourierRequest() CourierRequestRepository
	// Refill() RefillRepository
}

type RedisRepoI interface {
	SetUserState(ctx context.Context, chatID int64, user *models.CreateUser, ttl time.Duration) error
	GetUserState(ctx context.Context, chatID int64) (*models.CreateUser, error)
	DeleteUserState(ctx context.Context, chatID int64) error
}

type UserRepoI interface {
	Create(ctx context.Context, req *models.CreateUser) (*models.UserPrimaryKey, error)
	GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error)
	GetByPhone(ctx context.Context, req string) (*models.User, error)
	GetList(ctx context.Context, req *models.GetListUserRequest) (resp *models.GetListUserResponse, err error)
	Update(ctx context.Context, req *models.UpdateUser) (int64, error)
	Delete(ctx context.Context, req *models.UserPrimaryKey) (int64, error)
	GetUserByTelegramID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error)
	GetUnconfirmedCouriers(ctx context.Context) ([]models.User, error)
	GetByRole(ctx context.Context, role string) ([]models.User, error)
	Approve(ctx context.Context, userID string) error
	Reject(ctx context.Context, userID string) error
}

type BarrelRepository interface {
	Create(ctx context.Context, barrel *models.CreateBarrel) (*models.Barrel, error)
	GetByID(ctx context.Context, id int64) (*models.Barrel, error)
	GetList(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error)
	Update(ctx context.Context, barrel *models.UpdateBarrel) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	GetListSellerId(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error)
	GetBarrelBySellerId(ctx context.Context, sellerId string) (*models.Barrel, error)
	GetListEmpty(ctx context.Context, req *models.GetListBarrelRequest) (*models.GetListBarrelResponse, error)
}

// type SaleRepository interface {
// 	CreateSale(ctx context.Context, sale *models.Sale) (int64, error)
// 	GetSaleByID(ctx context.Context, id int64) (*models.Sale, error)
// 	GetSalesBySeller(ctx context.Context, sellerID int64) ([]*models.Sale, error)
// 	GetSalesByBarrel(ctx context.Context, barrelID int64) ([]*models.Sale, error)
// 	GetSaleList(ctx context.Context, limit, offset int) ([]*models.Sale, error)
// 	DeleteSale(ctx context.Context, id int64) error
// }

// type CourierRequestRepository interface {
// 	CreateCourierRequest(ctx context.Context, request *models.CourierRequest) (int64, error)
// 	GetCourierRequestByID(ctx context.Context, id int64) (*models.CourierRequest, error)
// 	GetPendingRequests(ctx context.Context) ([]*models.CourierRequest, error)
// 	UpdateCourierRequestStatus(ctx context.Context, id int64, status string) error
// 	DeleteCourierRequest(ctx context.Context, id int64) error
// }

// type RefillRepository interface {
// 	CreateRefill(ctx context.Context, refill *models.Refill) (int64, error)
// 	GetRefillByID(ctx context.Context, id int64) (*models.Refill, error)
// 	GetRefillsByCourier(ctx context.Context, courierID int64) ([]*models.Refill, error)
// 	GetRefillList(ctx context.Context, limit, offset int) ([]*models.Refill, error)
// 	DeleteRefill(ctx context.Context, id int64) error
// }
