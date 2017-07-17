package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tobyjsullivan/ues-sdk/event/reader"
	"github.com/urfave/negroni"
    "github.com/tobyjsullivan/ues-sdk/event"
    "encoding/json"
    "encoding/base64"
)

var (
	logger      *log.Logger
	eventReader *reader.EventReader
)

func init() {
	var err error
	logger = log.New(os.Stdout, "[svc] ", 0)

	eventReader, err = reader.New(&reader.EventReaderConfig{
        ServiceUrl: os.Getenv("EVENT_READER_API"),
    })
	if err != nil {
		logger.Println("Error initializing Event Reader API.", err.Error())
		panic(err.Error())
	}
}

func main() {
	r := buildRoutes()

	n := negroni.New()
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n.Run(":" + port)
}

func buildRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", statusHandler).Methods("GET")
	r.HandleFunc("/histories/{headId}", fetchHistoryHandler).Methods("GET")

	return r
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The service is online!\n")
}

func fetchHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramHeadId := vars["headId"]
	if paramHeadId == "" {
		http.Error(w, "Must supply headId in path.", http.StatusBadRequest)
		return
	}

    headId := event.EventID{}
    err := headId.Parse(paramHeadId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    results := getFullHistory(headId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    inverseHistory := make([]string, 0)
    for r := range results {
        if r.Error != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        e := r.Event
        id := e.ID()

        formatted := &eventFmt{
            ID: id.String(),
            Type: e.Type,
            Data: base64.StdEncoding.EncodeToString(e.Data),
        }

        js, err := json.Marshal(formatted)
        if err != nil {
            http.Error(w, fmt.Sprintf("%s: %s", id, err.Error()), http.StatusInternalServerError)
            return
        }

        inverseHistory = append(inverseHistory, string(js))
    }

    for i := len(inverseHistory) - 1; i >= 0; i-- {
        fmt.Fprintln(w, inverseHistory[i])
    }
}

type eventFmt struct {
    ID string `json:"id"`
    Type string `json:"type"`
    Data string `json:"data"`
}

type optionEvent struct {
    Event *event.Event
    Error error
}

func getFullHistory(id event.EventID) (chan *optionEvent) {
    results := make(chan *optionEvent, 10)

    go func(id event.EventID, results chan *optionEvent) {
        zero := event.EventID{}
        defer close(results)

        for id != zero {
            logger.Println("Getting event for history.", id.String())
            e, err := eventReader.GetEvent(id)
            if err != nil {
                results<- &optionEvent{Error: err}
                return
            }

            logger.Println("Got event of type:", e.Type)
            results<- &optionEvent{Event: e}
            id = e.PreviousEvent
            logger.Println("Next event:", id.String())
        }
    }(id, results)

    return results
}