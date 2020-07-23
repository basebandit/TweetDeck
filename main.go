package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var (
	errMissingQueryString = errors.New("missing order query string")
	errMissingCSVFile     = errors.New("missing csv file parameter")
	errInvalidOrder       = errors.New("no such order")
)

type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render sets the application-specific error code in AppCode.
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// errInvalidRequest returns status 422 Unprocessable Entity including error message.
func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// errInternalServerError returns status 500 Internal Server Error including error message.
func errInternalServerError(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		ErrorText:      err.Error(),
	}
}

type server struct {
	port    string
	router  *chi.Mux
	records []map[string]string //inmemory rows read from csv file.
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errMissingCSVFile)
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	srv := server{
		port:    ":5000",
		router:  r,
		records: csvToMap(bytes.NewReader(file)),
	}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	fmt.Println("[+]============== API ROUTES ===================")

	if err := chi.Walk(srv.routes(), walkFunc); err != nil {
		fmt.Printf("Logging err:%s\n", err.Error())
	}

	fmt.Println("==================================================")

	fmt.Printf("Server listening on %s\n", srv.port)

	if err := http.ListenAndServe(":5000", srv.routes()); err != nil {
		log.Fatalf("Server Error : %v", err)
	}
}

func (s *server) routes() *chi.Mux {
	s.router.Get("/api/v1/ping", s.handlePing)
	s.router.Get("/api/v1/feed/followers", s.handleSortFollowers)
	s.router.Get("/api/v1/feed/following", s.handleSortFollowing)
	s.router.Get("/api/v1/feed/tweets", s.handleSortTweets)
	s.router.Get("/api/v1/feed/likes", s.handleSortLikes)
	return s.router
}

func (s *server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (s *server) handleSortFollowers(w http.ResponseWriter, r *http.Request) {

	o := r.URL.Query().Get("o")

	if o == "" {
		log.Println(errMissingQueryString)
		render.Render(w, r, errInvalidRequest(errMissingQueryString))
		return
	}

	order, err := strconv.Atoi(o)
	if err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(errInvalidOrder))
		return
	}

	if err = sortByFollowers(s.records, order); err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Respond(w, r, s.records)
}

func (s *server) handleSortFollowing(w http.ResponseWriter, r *http.Request) {
	o := r.URL.Query().Get("o")
	if o == "" {
		log.Println(errMissingQueryString)
		render.Render(w, r, errInvalidRequest(errMissingQueryString))
		return
	}

	order, err := strconv.Atoi(o)
	if err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(errInvalidOrder))
		return
	}

	if err = sortByFollowing(s.records, order); err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Respond(w, r, s.records)
}

func (s *server) handleSortTweets(w http.ResponseWriter, r *http.Request) {

	o := r.URL.Query().Get("o")

	if o == "" {
		log.Println(errMissingQueryString)
		render.Render(w, r, errInvalidRequest(errMissingQueryString))
		return
	}

	order, err := strconv.Atoi(o)
	if err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(errInvalidOrder))
		return
	}

	if err = sortByTweets(s.records, order); err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Respond(w, r, s.records)
}

func (s *server) handleSortLikes(w http.ResponseWriter, r *http.Request) {

	o := r.URL.Query().Get("o")

	if o == "" {
		log.Println(errMissingQueryString)
		render.Render(w, r, errInvalidRequest(errMissingQueryString))
		return
	}

	order, err := strconv.Atoi(o)
	if err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(errInvalidOrder))
		return
	}

	if err = sortByLikes(s.records, order); err != nil {
		log.Println(err)
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Respond(w, r, s.records)
}

func csvToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)

	rows := []map[string]string{}
	var header []string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			//the first read/line is the header.
			header = record
		} else {
			dict := map[string]string{}
			for i := range header { //loop through the csv header
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

//sort by order 1: ascending  -1 : descending
func sortByFollowers(records []map[string]string, order int) error {
	switch order {
	case -1: //descending
		sortDescending(records, "followers")
		return nil
	case 1: //ascending
		sortAscending(records, "followers")
		return nil
	default:
		return errInvalidOrder
	}
}

//sort by order 1: ascending  -1 : descending
func sortByFollowing(records []map[string]string, order int) error {
	switch order {
	case -1: //descending
		sortDescending(records, "following")
		return nil
	case 1: //ascending
		sortAscending(records, "following")
		return nil
	default:
		return errInvalidOrder
	}
}

//sort by order 1: ascending  -1 : descending
func sortByTweets(records []map[string]string, order int) error {
	switch order {
	case -1: //descending
		sortDescending(records, "tweets")
		return nil
	case 1: //ascending
		sortAscending(records, "tweets")
		return nil
	default:
		return errInvalidOrder
	}
}

//sort by order 1: ascending  -1 : descending
func sortByLikes(records []map[string]string, order int) error {
	switch order {
	case -1: //descending
		sortDescending(records, "likes")
		return nil
	case 1: //ascending
		sortAscending(records, "likes")
		return nil
	default:
		return errInvalidOrder
	}
}

//sort by col in ascending order from lowest - highest
func sortAscending(in []map[string]string, col string) {
	sort.SliceStable(in, func(i, j int) bool {
		a, err := strconv.Atoi(in[i][col])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(in[j][col])
		if err != nil {
			log.Fatal(err)
		}
		return a < b
	})
}

//sort by col in descending order from highest - lowest
func sortDescending(in []map[string]string, col string) {
	sort.SliceStable(in, func(i, j int) bool {
		a, err := strconv.Atoi(in[i][col])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(in[j][col])
		if err != nil {
			log.Fatal(err)
		}
		return a > b
	})
}
