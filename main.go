package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal" //"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/meshhq/golang-html-template-tutorial/assets"
)

// Templates
var homepageTpl *template.Template
var signinTpl *template.Template
var signupTpl *template.Template
var feedTpl *template.Template

func init() {
	homepageHTML := assets.MustAssetString("templates/index.html")
	homepageTpl = template.Must(template.New("homepage_view").Parse(homepageHTML))

	signupHTML := assets.MustAssetString("templates/signup.html")
	signupTpl = template.Must(template.New("signup_view").Parse(signupHTML))

	signinHTML := assets.MustAssetString("templates/signin.html")
	signinTpl = template.Must(template.New("signin_view").Parse(signinHTML))

	feedHTML := assets.MustAssetString("templates/feed.html")
	feedTpl = template.Must(template.New("feed_view").Parse(feedHTML))
}

func main() {
	serverCfg := Config{
		Host:         "localhost:8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	htmlServer := Start(serverCfg)
	defer htmlServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("main : shutting down")
}

// Config provides basic configuration
type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// HTMLServer represents the web service that serves up HTML
type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

// Start launches the HTML Server
func Start(cfg Config) *HTMLServer {
	// Setup Context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Handlers
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/index", HomeHandler)
	router.HandleFunc("/signup", SignupHandler)
	router.HandleFunc("/signin", SigninHandler)
	router.HandleFunc("/feed", FeedHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Create the HTML Server
	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           cfg.Host,
			Handler:        router,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}

	// Add to the WaitGroup for the listener goroutine
	htmlServer.wg.Add(1)

	// Start the listener
	go func() {
		fmt.Printf("\nHTMLServer : Service started : Host=%v\n", cfg.Host)
		htmlServer.server.ListenAndServe()
		htmlServer.wg.Done()
	}()

	return &htmlServer
}

// Stop turns off the HTML Server
func (htmlServer *HTMLServer) Stop() error {
	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("\nHTMLServer : Service stopping\n")

	// Attempt the graceful shutdown by closing the listener
	// and completing all inflight requests
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		// Looks like we timed out on the graceful shutdown. Force close.
		if err := htmlServer.server.Close(); err != nil {
			fmt.Printf("\nHTMLServer : Service stopping : Error=%v\n", err)
			return err
		}
	}

	// Wait for the listener to report that it is closed.
	htmlServer.wg.Wait()
	fmt.Printf("\nHTMLServer : Stopped\n")
	return nil
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

// Push the given resource to the client.
func push(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)
	if ok {
		if err := pusher.Push(resource, nil); err == nil {
			return
		}
	}
}

// Route Handlers

// HomeHandler renders the homepage view template
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/bootstrap.min.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	RegisterUser(w, r)

	fullData := map[string]interface{}{}
	render(w, r, homepageTpl, "homepage_view", fullData)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/bootstrap.min.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fullData := map[string]interface{}{}
	render(w, r, signupTpl, "signup_view", fullData)
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/bootstrap.min.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fullData := map[string]interface{}{}
	render(w, r, signinTpl, "signin_view", fullData)
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/bootstrap.min.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	SearchUser(w, r)

	fullData := map[string]interface{}{}
	render(w, r, feedTpl, "feed_view", fullData)
}
