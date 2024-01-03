package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"wails-template/components"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist components
var assets embed.FS
var version = "0.0.0"

func main() {
	r := NewChiRouter()

	app := application.New(application.Options{
		Name:        "wails3-htmx-template",
		Description: "A demo of using raw HTML & CSS",
		Assets: application.AssetOptions{
			FS: assets,
			Middleware: func(next http.Handler) http.Handler {
				r.NotFound(next.ServeHTTP)
				return r
			},
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F12": func(window *application.WebviewWindow) {
				println("teste")
				window.Show()
			},
		},
	})

	// Register for events
	app.Events.On("myevent", func(e *application.WailsEvent) {
		fmt.Println("event run")
		app.Logger.Info("[Go] WailsEvent received", "name", e.Name, "data", e.Data, "sender", e.Sender, "cancelled", e.Cancelled)
	})

	// Create window
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Main Window",
		CSS:   `body { background-color: rgba(255, 255, 255, 0); } .main { color: white; margin: 20%; }`,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F12": func(window *application.WebviewWindow) {
				println("teste")
				window.Show()
			},
		},
		URL: "/",
	})

	// Systray Window
	systemTray := app.NewSystemTray()
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:          "systray",
		Width:         500,
		Height:        800,
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F12": func(window *application.WebviewWindow) {
				systemTray.OpenMenu()
			},
		},
		URL: "/systray/",
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:  "view",
		Title: "View Window",
		// Width:         800,
		// Height:        800,
		Frameless:     false,
		AlwaysOnTop:   false,
		Hidden:        true,
		DisableResize: true,
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		// KeyBindings: map[string]func(window *application.WebviewWindow){
		// 	"F12": func(window *application.WebviewWindow) {
		// 		systemTray.OpenMenu()
		// 	},
		// },
		URL: "/",
	})

	// Systray Menu
	myMenu := app.NewMenu()
	systemTray.SetMenu(myMenu)
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// Attach extra windows
	systemTray.AttachWindow(window).WindowOffset(5)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func NewChiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := &Counter{}

	r.Get("/initial", templ.Handler(components.Pages([]struct {
		Path  string
		Label string
	}{
		{"/greet", "Greet form"},
		{"/events", "Events page"},
	}, struct {
		Version string
		Text    string
	}{
		version, "No update available",
	})).ServeHTTP)
	r.Get("/greet", templ.Handler(components.GreetForm("/greet")).ServeHTTP)
	r.Post("/greet", components.Greet)
	r.Get("/modal", templ.Handler(components.TestPage("#modal", "outerHTML")).ServeHTTP)
	r.Post("/modal", templ.Handler(components.Modal("Title for the modal", "Sample Data")).ServeHTTP)
	r.Get("/sidebar", templ.Handler(components.SideBar()).ServeHTTP)
	r.Get("/counter", CounterHandler(c))
	r.Get("/events", templ.Handler(components.Events()).ServeHTTP)
	// Custom Endpoints
	r.Get("/event", TestLoop)
	return r
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

func TestLoop(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.Modal("title", "test ")).ServeHTTP(w, r)
	fmt.Println("test event")
	application.Get().Events.Emit(&application.WailsEvent{
		Name: "myevent",
		Data: "hello",
	})
	fmt.Println("test loop")
	for i := 0; i < 3; i++ {
		w.Write([]byte("test loop"))
		time.Sleep(time.Second * 2)
	}
}

// func TestLoop(w http.ResponseWriter, r *http.Request) {
// 	// templ.Handler(components.Modal("title", "test ")).ServeHTTP(w, r)
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
// 		return
// 	}

// 	for i := 0; i < 3; i++ {
// 		w.Write([]byte("test loop"))
// 		flusher.Flush()
// 		time.Sleep(time.Second * 2)
// 	}
// }

/* TODO
UX Idea:
So while keeping the default wails template view, we can add a hover menu that opens an overlay grid with the options for different views,
that then uses htmx + astro to render the different views.

Systray usage:
New idea:
Change notifications to be purely for notifications, with settings and updates done in the main window via bottom/top/sidebar.
Old:
 - User/System settings0
 - Notifications
 - Realtime updates (time, background tasks)

Systray could be used for notifcations/modals, a links or options for navigation, or just a demo for opening alternate windows
settings/notifications/realtime updates for status information/widgets and action shortcuts, background tasks&activites


- Add support for htmx
no js - done (except for the 1 small demo)
add greet demo -
chi - done?
tailwind - done
daisyui
systray - added, but will be a view controller or something?
multiwindow - what for?
astro - done

add embedded files demo/user config/embedded binary - store in webview storage
streaming? sse/websockets/chunked/events?
notifications?
db

1. go to events page which loads the page
2. button to activate the call to events
3. return instant result then loop to send an event with updates.


*/
