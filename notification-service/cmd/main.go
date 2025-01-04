package main

import (
	"fmt"
	"lld/pkg/notifications"
	"sync"
	"time"
)

func CryptoNotificationSystem() {

	// code goes here
	notificationModes := []notifications.NotificationMode{}
	notificationModes = append(notificationModes, notifications.NotificationModeSMS)
	user1 := notifications.NewUser("vaibhav", notificationModes)

	notification := notifications.Notification{
		Name:        "bicoin-price-change",
		Coin:        "bitcoin",
		Subscribers: make([]*notifications.User, 0),
	}

	notification2 := notifications.Notification{
		Name:        "matic-price-change",
		Coin:        "matic",
		Subscribers: make([]*notifications.User, 0),
	}

	//notification.Subscribers = append(notification.Subscribers, user1)

	service := notifications.NewNotificationService()
	service.CreateNotification(notification2)
	service.CreateNotification(notification)
	service.Run()

	service.Subscribe(user1, "bitcoin", "bicoin-price-change")

	//service.Notifications["bitcoin"] = map[string]*notifications.Notification{}

	price := 100000
	for i := 0; i < 10; i++ {
		price = price + i
		ni := notifications.NotificationInstance{
			Id:   fmt.Sprintf("%d", i),
			Name: "bicoin-price-change",
			Coin: "bitcoin",

			Data: map[string]interface{}{
				"price":     fmt.Sprintf("%d", price),
				"updatedAt": time.Now().Local().String(),
			},
		}
		service.PushNotification(ni)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		price := 100
		for i := 0; i < 50; i++ {
			price = price + i
			ni := notifications.NotificationInstance{
				Id:   fmt.Sprintf("%d", i),
				Name: "bicoin-price-change",
				Coin: "bitcoin",

				Data: map[string]interface{}{
					"price":     fmt.Sprintf("%d", price),
					"updatedAt": time.Now().Local().String(),
				},
			}
			service.PushNotification(ni)
			time.Sleep(100 * time.Millisecond) // simulate delay

		}
		wg.Done()
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		price = 100
		for i := 0; i < 50; i++ {
			price = price + i
			ni := notifications.NotificationInstance{
				Id:   fmt.Sprintf("%d", i),
				Name: "matic-price-change",
				Coin: "matic",
				Data: map[string]interface{}{
					"price":     fmt.Sprintf("%d", price),
					"updatedAt": time.Now().Local().String(),
				},
			}
			time.Sleep(100 * time.Millisecond)
			service.PushNotification(ni)

		}
		wg.Done()
	}(&wg)

	// Demonstrate unsubscribe functionality
	time.Sleep(1 * time.Second)
	service.Subscribe(user1, "matic", "matic-price-change")
	time.Sleep(1 * time.Second)

	service.Unsubscribe(user1, "bitcoin", "bicoin-price-change")
	time.Sleep(2 * time.Second)

	// Demonstrate change notification mode functionality
	newModes := []notifications.NotificationMode{notifications.NotificationModeEMAIL}
	service.ChangeNotificationMode(user1, newModes)

	wg.Wait()

}

func main() {

	// do not modify below here, readline is our function
	// that properly reads in the input for you
	CryptoNotificationSystem()

}
