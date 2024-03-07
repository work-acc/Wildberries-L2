package service

import (
	"fmt"
	"time"

	"github.com/work-acc/Wildberries-L2/dev11/internal/model"
)

func New() *Service {
	return &Service{}
}

type Service struct{}

// эмуляция БД
var emulationDB = make(map[time.Time][]model.Event)

func (s *Service) Create(event model.Event) error {
	records, ok := emulationDB[event.Date]
	if ok {
		for i := 0; i < len(records); i++ {
			if records[i].UserID == event.UserID {
				return fmt.Errorf("the record already exists in the database")
			}
		}
	}

	emulationDB[event.Date] = append(emulationDB[event.Date], event)

	return nil
}

func (s *Service) Update(event model.Event) error {
	records, ok := emulationDB[event.Date]
	if ok {
		for i := 0; i < len(records); i++ {
			if records[i].UserID == event.UserID {
				records[i].Title = event.Title
			}
		}
	}

	return fmt.Errorf("the event was not found")
}

func (s *Service) Delete(event model.Event) error {
	records, ok := emulationDB[event.Date]
	if ok {
		for i := 0; i < len(records); i++ {
			if records[i].UserID == event.UserID {
				records[i] = records[len(records)-1]
				records = records[:len(records)-1]
			}
		}
	}

	return fmt.Errorf("the event was not found")
}

func (s *Service) EventsForDay(date time.Time) (events []model.Event, err error) {
	events, ok := emulationDB[date]
	if ok {
		return events, nil
	}

	return events, fmt.Errorf("the event was not found")
}

func (s *Service) EventsForWeek(date time.Time) (events []model.Event, err error) {
	startOfWeek := date.AddDate(0, 0, -int(date.Weekday())+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	for currentDay := startOfWeek; currentDay.Before(endOfWeek.AddDate(0, 0, 1)); currentDay = currentDay.AddDate(0, 0, 1) {
		events = append(events, emulationDB[currentDay]...)
	}

	if events == nil {
		return events, fmt.Errorf("the event was not found")
	}

	return events, nil
}

func (s *Service) EventsForMonth(date time.Time) (events []model.Event, err error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	for currentDay := startOfMonth; currentDay.Before(endOfMonth.AddDate(0, 0, 1)); currentDay = currentDay.AddDate(0, 0, 1) {
		events = append(events, emulationDB[currentDay]...)
	}

	if events == nil {
		return events, fmt.Errorf("the event was not found")
	}

	return events, nil
}
