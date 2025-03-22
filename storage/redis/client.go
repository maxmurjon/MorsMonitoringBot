package redis

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"morc/models"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Redis serverga ulanishni tekshirish
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis-ga ulanishda xatolik: %v", err)
	}

	return &RedisClient{
		client: rdb,
	}
}

// Foydalanuvchi ma'lumotlarini Redis-ga saqlash
func (r *RedisClient) SetUserState(ctx context.Context, chatID int64, user *models.CreateUser, ttl time.Duration) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, getUserKey(chatID), data, ttl).Err()
}

// Foydalanuvchi ma'lumotlarini Redis-dan olish
func (r *RedisClient) GetUserState(ctx context.Context, chatID int64) (*models.CreateUser, error) {
	data, err := r.client.Get(ctx, getUserKey(chatID)).Result()
	if err == redis.Nil {
		// Kalit topilmasa, foydalanuvchi hali saqlanmagan deb hisoblaymiz
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user models.CreateUser
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Foydalanuvchini Redis-dan o‘chirish
func (r *RedisClient) DeleteUserState(ctx context.Context, chatID int64) error {
	return r.client.Del(ctx, getUserKey(chatID)).Err()
}

// Foydalanuvchi uchun Redis kalitini yaratish
func getUserKey(chatID int64) string {
	return "user_state:" + strconv.FormatInt(chatID, 10)
}
