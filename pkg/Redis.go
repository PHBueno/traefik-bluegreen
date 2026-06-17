package pkg

// type ClientTenant struct {
// 	app    string
// 	slot   string
// 	tenant string
// }

// func NewRedisConnection(address string, password string, database int) (*redis.Client, error) {
// 	context := context.Background()

// 	redisClient := redis.NewClient(
// 		&redis.Options{
// 			Addr:     address,
// 			Password: password,
// 			DB:       database,
// 		},
// 	)

// 	_, err := redisClient.Ping(context).Result()

// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, "[REDIS ERROR] => ", err)
// 		return nil, err
// 	}

// 	fmt.Fprintln(os.Stdout, "Successfully connecting to Redis!")
// 	return redisClient, nil
// }
