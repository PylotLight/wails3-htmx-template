package main

import (
	"embed"
	"log"
	"net/http"

	"wails-template/components"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
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

	r.Get("/initial", templ.Handler(components.Pages([]struct {
		Path  string
		Label string
	}{
		{"/greet", "Greet form"},
	}, struct {
		Version string
		Text    string
	}{
		version, "No update available",
	})).ServeHTTP)
	r.Get("/greet", templ.Handler(components.GreetForm("/greet")).ServeHTTP)
	r.Post("/greet", components.Greet)
	r.Get("/modal", templ.Handler(components.TestPage("#modal", "outerHTML")).ServeHTTP)
	r.Post("/modal", templ.Handler(components.ModalPreview("Title for the modal", "Sample Data")).ServeHTTP)
	return r
}

/* TODO
UX Idea:
So while keeping the default wails template view, we can add a hover menu that opens an overlay grid with the options for different views,
that then uses htmx + astro to render the different views.


- Add support for htmx
no js
chi - done?
tailwind - done
daisyui
systray - added, but will be a view controller or something?
multiwindow
astro?


*/
