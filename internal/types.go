package types

import "sync"

type Counter struct {
	Count int
}

type Systray struct {
	Notifications string
	Settings      string
}

type Notification struct {
	Title   string
	Content string
}

type NotificationData struct {
	NotificationsMutex  sync.Mutex
	ActiveNotifications []Notification
	NotificationChannel chan []Notification
}

var Notifications = NotificationData{
	NotificationsMutex:  sync.Mutex{},
	ActiveNotifications: []Notification{},
	NotificationChannel: make(chan []Notification),
}

type NotificationAccessor interface {
	GetNotifications() []Notification
	AddNotification(notification Notification)
}

func (n *NotificationData) GetNotifications() []Notification {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()
	return n.ActiveNotifications
}

func (n *NotificationData) AddNotification(notification Notification) {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()
	n.ActiveNotifications = append(n.ActiveNotifications, notification)
}
