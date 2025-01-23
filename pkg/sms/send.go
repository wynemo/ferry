package sms

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

const URL = "172.16.100.2:8513"

func sendSMS(phone, messageContent, code string) error {
	// Define the SMS gateway URL and query parameters
	baseURL := "http://%s/sms/Api/Send.do?%s"
	spCode := "1001"
	loginName := "test"
	password := "111111"
	currentTime := time.Now().Format("20060102150405")

	// Construct the query parameters
	params := url.Values{}
	params.Add("SpCode", spCode)
	params.Add("LoginName", loginName)
	params.Add("Password", password)
	params.Add("MessageContent", fmt.Sprintf("%s%s", messageContent, code))
	params.Add("UserNumber", phone)
	params.Add("ScheduleTime", currentTime)

	// Combine base URL with query parameters
	fullURL := fmt.Sprintf(baseURL, URL, params.Encode())

	// Make the HTTP GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	fmt.Println("SMS sent successfully!")
	return nil
}

// GenerateSMSCode generates a random 6-digit SMS verification code and stores it in Redis
func GenerateSMSCode(phone string) (string, error) {
	// 初始化 Redis 客户端
	// 从配置中获取 Redis URL，并直接替换掉前缀 "redis://"
	redisURL := viper.GetString("settings.redis.url")
	redisAddr := strings.Replace(redisURL, "redis://", "", 1) // 替换前缀

	fmt.Println("---------------")
	fmt.Println(redisAddr)
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Generate a 6-digit random code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	codeStr := fmt.Sprintf("%06d", code)

	// Define the Redis key and expiration time
	redisKey := fmt.Sprintf("sms_%s", phone)
	expiration := 5 * time.Minute // Set expiration time for the key

	// Store the code in Redis
	err := client.Set(redisKey, codeStr, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store SMS code in Redis: %w", err)
	}

	fmt.Printf("Generated SMS code for %s: %s\n", phone, codeStr)
	return codeStr, nil
}
