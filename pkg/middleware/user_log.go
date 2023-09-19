package middleware

import (
	"context"
	"fmt"
	"myapp/pkg/configs"
	"myapp/src/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func SaveUserLogs(c *fiber.Ctx, user models.User) {
	ctxTodo := context.TODO()
	now := time.Now()

	lastLoginID := fmt.Sprintf("LastLoginID++%d", user.ID)

	// Delete Redis for testing purpose
	// _, err = configs.RedisClient.Del(ctxTodo, lastLoginID).Result()
	// if err != nil {
	// 	return nil
	// }

	logTTL := 15 * time.Minute
	lastLog, errRedis := configs.RedisClient.Get(ctxTodo, lastLoginID).Result()
	if errRedis == redis.Nil {
		// create redis instance
		// fmt.Println("RedisClient.Set")
		errRedisSet := configs.RedisClient.Set(ctxTodo, lastLoginID, now, logTTL).Err()
		if errRedisSet != nil {
			log.Errorf("RedisClient.Set Error: %s", errRedisSet.Error())
		}
	}

	// fmt.Println(lastLog)
	if lastLog == "" {
		// fmt.Println("Save new LOG")
		// update user last login and IP
		err := configs.DB.Model(&user).Select("LastLoginAt", "LastLoginIp").
			Updates(models.User{LastLoginAt: &now, LastLoginIp: c.IP()}).Error
		if err != nil {
			log.Errorf(fmt.Sprintf(err.Error()))
		}
	}
}
