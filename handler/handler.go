package handler

import (
	"context"
	"net/http"
	"strconv"

	"wails3-htmx-template/components"
	types "wails3-htmx-template/internal"

	"github.com/a-h/templ"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func InitContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		windowID := r.Header.Get("X-Wails-Window-Id")
		app := application.Get()
		app.Events.Emit(&application.WailsEvent{
			Name: "myevent",
			Data: "now",
		})
		if windowID == "1" {
			components.Index().Render(ctx, w)
		}
		if windowID == "2" {
			components.Systray(types.Systray{Notifications: "active", Settings: ""}, types.Notifications.GetAllNotifications(), types.Settings{DatabaseLocation: "", SecretToken: ""}).Render(ctx, w)
		}
	}
}

func CounterHandler(c *types.Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Count++
		w.Write([]byte("count is " + strconv.Itoa(c.Count)))
	}
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templ.Handler(components.Notifications(types.Notifications.GetLatestNotificationsSinceLastPoll())).ServeHTTP(w, r)
	case "DELETE":
		idValue := r.PathValue("id")
		ID, err := strconv.Atoi(idValue)
		if err != nil {
			types.Notifications.AddNotification(types.Notification{Title: "Failed to convert ID", Message: "Failed to convert delete ID from string to int"})
		}
		err = types.Notifications.DeleteNotification(ID)
		if err != nil {
			types.Notifications.AddNotification(types.Notification{Title: "Failed to delete ID", Message: "Failed to delete ID"})
		}
		// types.Notifications.AddNotification(types.Notification{Title: "Successfully deleted ID: " + idValue, Message: "Notification Cleared"})
		templ.Handler(components.Notifications(types.Notifications.GetAllNotifications())).ServeHTTP(w, r)
	}

}

func SysTrayHandler(w http.ResponseWriter, r *http.Request) {
	button := r.PathValue("button")
	activeStates := types.Systray{}
	switch button {
	case "notifications":
		activeStates.Notifications = "active"
		activeStates.Settings = ""
		templ.Handler(components.Systray(activeStates, types.Notifications.GetAllNotifications(), types.Settings{})).ServeHTTP(w, r)
	case "settings":
		activeStates.Notifications = ""
		activeStates.Settings = "active"
		templ.Handler(components.Systray(activeStates, nil, types.Settings{DatabaseLocation: "SamplePath", SecretToken: "****"})).ServeHTTP(w, r)
	}
}
