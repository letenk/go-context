package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "foo", "Hello")
	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result: ", val)
	fmt.Println("took:", time.Since(start))
}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	val := ctx.Value("foo")
	fmt.Println(val)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respChannel := make(chan Response)
	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		respChannel <- Response{
			value: val,
			err:   err,
		}
	}()

	for {
		select {
		// If case context is done return erro
		case <-ctx.Done():
			return 0, fmt.Errorf("fething data from thirf party took to long, err: %s\n", ctx.Err())
		case resp := <-respChannel:
			return resp.value, resp.err
		}
	}
}

func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	time.Sleep(time.Millisecond * 150)
	return 666, nil
}
