package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "wails3-htmx-template",
		Description: "A demo of using raw HTML & CSS",
		Bind: []any{
			&GreetService{},
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Events.Emit(&application.WailsEvent{
				Name: "time",
				Data: now,
			})
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}

// func NewChiRouter() *chi.Mux {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	c := &Counter{}

// 	// r.Get("/initial", templ.Handler(components.Pages([]struct {
// 	// 	Path  string
// 	// 	Label string
// 	// }{
// 	// 	{"/greet", "Greet form"},
// 	// 	{"/events", "Events page"},
// 	// }, struct {
// 	// 	Version string
// 	// 	Text    string
// 	// }{
// 	// 	version, "No update available",
// 	// })).ServeHTTP)
// 	r.Get("/init", InitContent())
// 	r.Get("/greet", templ.Handler(components.GreetForm("/greet")).ServeHTTP)
// 	r.Post("/greet", components.Greet)
// 	r.Get("/modal", templ.Handler(components.TestPage("#modal", "outerHTML")).ServeHTTP)
// 	r.Post("/modal", templ.Handler(components.Modal("Title for the modal", "Sample Data")).ServeHTTP)
// 	// r.Get("/systray", InitContent())
// 	// r.Get("/sidebar", templ.Handler(components.SideBar()).ServeHTTP)
// 	r.Get("/counter", CounterHandler(c))
// 	r.Get("/events", templ.Handler(components.Events()).ServeHTTP)
// 	// Custom Endpoints
// 	r.Get("/event", TestLoop)
// 	return r
// }

/* TODO
1. Add systray
2. Add user settings
3. Add notifications
4. Add realtime updates
5. Add multiwindow
6. Add mouse bound window for context interaction/frameless/global keybindings
7. Add daisyui/other
8. Add DB + local storage (bbolt/pocketbase)
9. Add dialog
10. Add Auth
11. Embedded files + audio notification


UX Idea:
So while keeping the default wails template view, we can add a hover menu that opens an overlay grid with the options for different views,
that then uses htmx + astro to render the different views.

Systray usage:
New idea 2:
Use bottom nav menu in systray to default to notfications view but change to settings/updates. Use events and wml to navigate from notifications to main window.

New idea:
Change notifications to be purely for notifications, with settings and updates done in the main window via bottom/top/sidebar.
Old:
 - User/System settings0
 - Notifications
 - Realtime updates (time, background tasks)

Systray could be used for notifcations/modals, a links or options for navigation, or just a demo for opening alternate windows
settings/notifications/realtime updates for status information/widgets and action shortcuts, background tasks&activites

###Stuck###
- stuck on how to manage route for different windows i.e systray for loading the relevant components/html/css for each window - fixed


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
