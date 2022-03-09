package text_info

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	url := os.Getenv("REDISCLOUD_URL")
	if url == "" {
		log.Fatalln("REDISCLOUD_URL env variable not set")
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opt)

	//updateSchema()
}

// just gets each user and then upsert's the user. practically speaking it adds json fields that aren't there yet.
func updateSchema() {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		user, err := getUser(key)
		if err != nil {
			panic(err)
		}
		upsertUser(user)
	}
}

type User struct {
	PhoneNumber        string `json:"phoneNumber"`
	FreeMessagesUsed   int    `json:"freeMessages"`
	Private            bool   `json:"private"`
	Subscribed         bool   `json:"subscribed"`
	SubscriptionAmount int    `json:"subscriptionAmount"`
	PaidMessages       int    `json:"paidMessages"`
}

func getUser(phoneNumber string) (*User, error) {
	m, err := rdb.Get(ctx, phoneNumber).Result()
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal([]byte(m), &u)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func addNewUser(phoneNumber string) error {
	user := User{
		PhoneNumber:        phoneNumber,
		FreeMessagesUsed:   0,
		Private:            false,
		Subscribed:         false,
		SubscriptionAmount: 0,
	}
	return upsertUser(&user)
}

func upsertUser(user *User) error {

	jsonMichael, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, user.PhoneNumber, string(jsonMichael), 0).Err()
}
