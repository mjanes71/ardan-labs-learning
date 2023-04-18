// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux // this is an example of an embedded pointer, allows every part of the contextmux to promote up to the parent struct (App)
	shutdown chan os.Signal // so you basically steal everythign it means to be an httptreemux and add it to your app object
	mw       []Middleware // all the httptreemux receiver functions still work
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw: 		mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux. this overrides the ttpptreemux Handle function
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := context.WithValue(r.Context(), key, &v)

		if err := handler(ctx, w, r); err != nil {

			// HANDLE ERROR
			return
		}

		// INJECT BUSINESS LAYER CODE
	}

	a.ContextMux.Handle(method, path, h)
}