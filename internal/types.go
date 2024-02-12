package types

import "sync"

type Counter struct {
	Count int
}

type Systray struct {
	Notifications string
	Settings      string
}

type Settings struct {
	DatabaseLocation string
	SecretToken      string
}

type Notification struct {
	ID      int
	Title   string
	Message string
}

type NotificationData struct {
	NotificationsMutex  sync.Mutex
	ActiveNotifications []Notification
	NotificationChannel chan []Notification
	lastPolledIndex     int
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

func (n *NotificationData) GetLatestNotificationsSinceLastPoll() []Notification {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()
	newNotifications := n.ActiveNotifications[n.lastPolledIndex:]
	n.lastPolledIndex = len(n.ActiveNotifications)
	return newNotifications
}

func (n *NotificationData) GetAllNotifications() []Notification {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()
	newNotifications := n.ActiveNotifications[:n.lastPolledIndex]
	return newNotifications
}

func (n *NotificationData) AddNotification(notification Notification) {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()
	n.ActiveNotifications = append(n.ActiveNotifications, notification)
}

func (n *NotificationData) DeleteNotification(id int) error {
	n.NotificationsMutex.Lock()
	defer n.NotificationsMutex.Unlock()

	// Find the index of the notification with the given ID
	for i, notification := range n.ActiveNotifications {
		if notification.ID == id {
			n.ActiveNotifications = append(n.ActiveNotifications[:i], n.ActiveNotifications[i+1:]...)
			return nil
		}
	}

	return nil
}
