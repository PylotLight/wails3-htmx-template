package handler

import (
	"context"
	"net/http"
	"strconv"

	"wails3-htmx-template/components"
	types "wails3-htmx-template/internal"

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
			components.Systray(types.Systray{Notifications: "active", Settings: ""}, types.Notifications.GetNotifications()).Render(ctx, w)
		}
	}
}

func CounterHandler(c *types.Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Count++
		w.Write([]byte("count is " + strconv.Itoa(c.Count)))
	}
}
