package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/work-acc/Wildberries-L2/dev11/internal/config"
	"github.com/work-acc/Wildberries-L2/dev11/internal/model"
	"github.com/work-acc/Wildberries-L2/dev11/internal/service"
)

func New(cfg config.Config, service *service.Service) *Router {
	router := http.NewServeMux()
	r := Router{
		router:  router,
		service: service,
	}

	srv := &http.Server{
		Addr:           cfg.Api.Addr,
		Handler:        r.middlewareLogs(router),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
	}

	r.srv = srv

	r.router.HandleFunc("/events_for_day", r.eventsForDay)
	r.router.HandleFunc("/events_for_week", r.eventsForWeek)
	r.router.HandleFunc("/events_for_month", r.eventsForMonth)

	r.router.HandleFunc("/create_event", r.createEvent)
	r.router.HandleFunc("/update_event", r.updateEvent)
	r.router.HandleFunc("/delete_event", r.deleteEvent)

	return &r
}

type Router struct {
	srv     *http.Server
	router  *http.ServeMux
	service *service.Service
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) middlewareLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (rtr *Router) eventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	date, err := rtr.parseAndValidateDate(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	events, err := rtr.service.EventsForDay(date)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (rtr *Router) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	date, err := rtr.parseAndValidateDate(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	events, err := rtr.service.EventsForWeek(date)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (rtr *Router) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	date, err := rtr.parseAndValidateDate(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	events, err := rtr.service.EventsForMonth(date)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (rtr *Router) createEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	event, err := rtr.parseAndValidateParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := rtr.service.Create(event); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rtr *Router) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	event, err := rtr.parseAndValidateParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := rtr.service.Update(event); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rtr *Router) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	event, err := rtr.parseAndValidateParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := rtr.service.Delete(event); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Router) parseAndValidateParams(r *http.Request) (event model.Event, err error) {
	userID := r.FormValue("user_id")
	dateStr := r.FormValue("date")
	title := r.FormValue("title")

	path := r.URL.Path

	if path == "delete_event" {
		if userID == "" || dateStr == "" {
			return event, fmt.Errorf("the request parameters were passed incorrectly")
		}
	} else {
		if userID == "" || dateStr == "" || title == "" {
			return event, fmt.Errorf("the request parameters were passed incorrectly")
		}
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return event, fmt.Errorf("incorrect user_id")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return event, fmt.Errorf("incorrect формат даты")
	}

	return model.Event{
		UserID: userIDInt,
		Date:   date,
		Title:  title,
	}, nil
}

func (h *Router) parseAndValidateDate(r *http.Request) (time.Time, error) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("the request parameters were passed incorrectly")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	return date, nil
}
