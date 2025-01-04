package notifications

import (
	"fmt"
	"sync"
)

type NotificationService struct {
	Notifications     map[string]map[string]*Notification //Nofications of the coin
	NotificationStore map[string][]*NotificationInstance
	requests          chan NotificationInstance
	nimutext          sync.RWMutex
	nLock             sync.RWMutex
}

func NewNotificationService() NotificationService {
	return NotificationService{
		Notifications:     map[string]map[string]*Notification{},
		requests:          make(chan NotificationInstance),
		NotificationStore: map[string][]*NotificationInstance{},
	}
}

func (ns *NotificationService) CreateNotification(notification Notification) {
	ns.nLock.Lock()
	defer ns.nLock.Unlock()
	_, exists := ns.Notifications[notification.Coin]
	if !exists {
		ns.Notifications[notification.Coin] = make(map[string]*Notification)
	}

	ns.Notifications[notification.Coin][notification.Name] = &notification

	//ns.Notifications[notification.Coin][notification.Name] = append(ns.Notifications[notification.Coin][notification.Name], &notification)

}

// Push the notification to the store
func (ns *NotificationService) PushNotification(instance NotificationInstance) {
	ns.nimutext.Lock()
	_, exists := ns.NotificationStore[instance.Coin]
	if !exists {
		ns.NotificationStore[instance.Coin] = []*NotificationInstance{}
	}
	ns.NotificationStore[instance.Coin] = append(ns.NotificationStore[instance.Coin], &instance)
	ns.nimutext.Unlock()
	ns.requests <- instance
}

func (ns *NotificationService) Run() {
	go func() {
		for {
			select {
			case request := <-ns.requests:
				// Process the request
				err := ns.ProcessRequest(request)
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(request)

				// case <-e.stopChan:
				// 	return
			}
		}
	}()
}

func (ns *NotificationService) ProcessRequest(ni NotificationInstance) error {

	notifications, exists := ns.Notifications[ni.Coin]
	if !exists {
		return fmt.Errorf("notifications does not exist for the %s so skipping it", ni.Coin+ni.Name)
	}

	notification, exits := notifications[ni.Name]
	if !exits {
		return fmt.Errorf("Notification does not exist for the coin so skipping it")
	}
	for _, subscriber := range notification.Subscribers {

		for _, mode := range subscriber.NotificationModes {
			ns.sendNotification(mode, subscriber, ni)
		}

	}
	return nil
}

func (ns *NotificationService) sendNotification(mode NotificationMode, user *User, ni NotificationInstance) {
	fmt.Printf("sending notification for the user: %s on mode %s", user.UserName, mode)
	fmt.Println(ni)

}

// Subscribe a user to a notification
func (ns *NotificationService) Subscribe(user *User, coin string, notificationName string) {
	ns.nLock.Lock()
	defer ns.nLock.Unlock()
	if notifications, exists := ns.Notifications[coin]; exists {
		if notification, exists := notifications[notificationName]; exists {
			notification.Subscribers = append(notification.Subscribers, user)
			fmt.Printf("Subscribed the user %s to the notification %s of coin %s \n", user.UserName, notificationName, coin)
		}
	}
}

// Unsubscribe a user from a notification
func (ns *NotificationService) Unsubscribe(user *User, coin string, notificationName string) {
	ns.nLock.Lock()
	defer ns.nLock.Unlock()
	if notifications, exists := ns.Notifications[coin]; exists {
		if notification, exists := notifications[notificationName]; exists {
			for i, subscriber := range notification.Subscribers {
				if subscriber.UserName == user.UserName {
					fmt.Printf("Unsubscribed the user %s from the notification %s of coin %s", user.UserName, notificationName, coin)
					notification.Subscribers = append(notification.Subscribers[:i], notification.Subscribers[i+1:]...)
					break
				}
			}
		}
	}
}

// ChangeNotificationMode allows a user to change their notification mode
func (ns *NotificationService) ChangeNotificationMode(user *User, newModes []NotificationMode) {
	user.NotificationModes = newModes
	fmt.Printf("Changed the notification mode for the user %s\n", user.UserName)
}
