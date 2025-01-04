package notifications

import "fmt"

type NotificationMode string

const (
	NotificationModeSMS   = "SMS"
	NotificationModeEMAIL = "EMAIL"
	NotificationModePUSH  = "PUSH"
)

type User struct {
	UserName          string
	NotificationModes []NotificationMode
	Subscriptions     []*Notification //TODO: check if actually need this
}

func NewUser(name string, modes []NotificationMode) *User {
	return &User{
		UserName:          name,
		NotificationModes: modes,
	}
}

type Subscriber interface {
	Subscribe(coin string, notificationName string)
}

func (u *User) Subscribe(coin string, notificationName string) {
	fmt.Printf("Subscribed to the notification %s of coin %s", notificationName, coin)
}

// Schema for creating notification
type Notification struct {
	Name        string //bitcoin-price-change
	Coin        string //bitcoin
	DataSchema  map[string]interface{}
	Subscribers []*User //
}

type NotificationInstance struct {
	Id   string
	Name string //bitcoin-price-change
	Coin string //bitcoin
	Data map[string]interface{}
}
