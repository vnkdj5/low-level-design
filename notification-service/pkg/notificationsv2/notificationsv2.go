package notificationsv2

import (
	"fmt"
	"sync"
)

type NotificationMode string

const (
	NotificationModeSMS   NotificationMode = "SMS"
	NotificationModeEMAIL NotificationMode = "EMAIL"
	NotificationModePUSH  NotificationMode = "PUSH"
)

// Represents a user in the system
type User struct {
	UserName          string
	NotificationModes []NotificationMode
}

// Creates a new user
func NewUser(name string, modes []NotificationMode) *User {
	return &User{
		UserName:          name,
		NotificationModes: modes,
	}
}

// Represents the schema of a notification (e.g., "price change")
type NotificationSchema struct {
	Name       string                 // e.g., "price-change"
	DataSchema map[string]interface{} // Defines the structure of notification data
}

// Represents an actual notification instance (event data)
type NotificationInstance struct {
	ID         string
	SchemaName string
	Coin       string
	Data       map[string]interface{}
}

// Stores the schema and user subscriptions for notifications
type Notification struct {
	Schema      NotificationSchema
	Coin        string
	Subscribers []*User
}

// Manages all notifications and user subscriptions
type NotificationService struct {
	Notifications     map[string]map[string]*Notification // coin -> schema name -> notification
	NotificationStore map[string][]*NotificationInstance  // coin -> list of notification instances
	requests          chan NotificationInstance
	muNotifications   sync.RWMutex
	muStore           sync.RWMutex
}

// Creates a new NotificationService instance
func NewNotificationService() *NotificationService {
	return &NotificationService{
		Notifications:     make(map[string]map[string]*Notification),
		NotificationStore: make(map[string][]*NotificationInstance),
		requests:          make(chan NotificationInstance),
	}
}

// Adds a new notification schema
func (ns *NotificationService) AddNotificationSchema(coin string, schema NotificationSchema) {
	ns.muNotifications.Lock()
	defer ns.muNotifications.Unlock()

	if _, exists := ns.Notifications[coin]; !exists {
		ns.Notifications[coin] = make(map[string]*Notification)
	}

	ns.Notifications[coin][schema.Name] = &Notification{
		Schema:      schema,
		Coin:        coin,
		Subscribers: []*User{},
	}

	fmt.Printf("Added schema '%s' for coin '%s'\n", schema.Name, coin)
}

// Subscribes a user to a notification schema for a specific coin
func (ns *NotificationService) Subscribe(user *User, coin, schemaName string) error {
	ns.muNotifications.Lock()
	defer ns.muNotifications.Unlock()

	if schemas, exists := ns.Notifications[coin]; exists {
		if notification, exists := schemas[schemaName]; exists {
			notification.Subscribers = append(notification.Subscribers, user)
			fmt.Printf("User '%s' subscribed to '%s' notification of '%s'\n", user.UserName, schemaName, coin)
			return nil
		}
		return fmt.Errorf("notification schema '%s' does not exist for coin '%s'", schemaName, coin)
	}
	return fmt.Errorf("coin '%s' does not exist", coin)
}

// Unsubscribes a user from a notification schema for a specific coin
func (ns *NotificationService) Unsubscribe(user *User, coin, schemaName string) error {
	ns.muNotifications.Lock()
	defer ns.muNotifications.Unlock()

	if schemas, exists := ns.Notifications[coin]; exists {
		if notification, exists := schemas[schemaName]; exists {
			for i, subscriber := range notification.Subscribers {
				if subscriber.UserName == user.UserName {
					notification.Subscribers = append(notification.Subscribers[:i], notification.Subscribers[i+1:]...)
					fmt.Printf("User '%s' unsubscribed from '%s' notification of '%s'\n", user.UserName, schemaName, coin)
					return nil
				}
			}
			return fmt.Errorf("user '%s' is not subscribed to '%s' notification of '%s'", user.UserName, schemaName, coin)
		}
		return fmt.Errorf("notification schema '%s' does not exist for coin '%s'", schemaName, coin)
	}
	return fmt.Errorf("coin '%s' does not exist", coin)
}

// Pushes a new notification instance (event data)
func (ns *NotificationService) PushNotification(instance NotificationInstance) {
	ns.muStore.Lock()
	if _, exists := ns.NotificationStore[instance.Coin]; !exists {
		ns.NotificationStore[instance.Coin] = []*NotificationInstance{}
	}
	ns.NotificationStore[instance.Coin] = append(ns.NotificationStore[instance.Coin], &instance)
	ns.muStore.Unlock()

	// Add to the processing queue
	ns.requests <- instance
}

// Processes notification events asynchronously
func (ns *NotificationService) Run() {
	go func() {
		for instance := range ns.requests {
			err := ns.processInstance(instance)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}

// Processes a single notification instance
func (ns *NotificationService) processInstance(instance NotificationInstance) error {
	ns.muNotifications.RLock()
	defer ns.muNotifications.RUnlock()

	if schemas, exists := ns.Notifications[instance.Coin]; exists {
		if notification, exists := schemas[instance.SchemaName]; exists {
			for _, user := range notification.Subscribers {
				for _, mode := range user.NotificationModes {
					ns.sendNotification(mode, user, instance)
				}
			}
			return nil
		}
		return fmt.Errorf("notification schema '%s' does not exist for coin '%s'", instance.SchemaName, instance.Coin)
	}
	return fmt.Errorf("coin '%s' does not exist", instance.Coin)
}

// Sends a notification to a user
func (ns *NotificationService) sendNotification(mode NotificationMode, user *User, instance NotificationInstance) {
	fmt.Printf("Sending '%s' notification to user '%s' via '%s' for coin '%s'\n", instance.SchemaName, user.UserName, mode, instance.Coin)
}
