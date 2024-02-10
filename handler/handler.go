package handler

import (
	"context"
	"net/http"
	"strconv"

	"wails3-htmx-template/components"

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
			components.Systray(map[string]string{"notifications": "active", "settings": ""}, GetNotifications()).Render(ctx, w)
		}
	}
}

type Counter struct {
	count int
}

func CounterHandler(c *Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.count++
		w.Write([]byte("count is " + strconv.Itoa(c.count)))
	}
}

func GetNotifications() []map[string]string {
	ActiveNotifications := make([]map[string]string, 0)
	if "Query availble active notifications" != "" {
		ActiveNotifications = append(ActiveNotifications, map[string]string{"Title": "No Notifications", "Content": "There are no notifications at this time"})
		return ActiveNotifications
	}
	return ActiveNotifications
}
