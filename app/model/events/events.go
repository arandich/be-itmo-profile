package events

import (
	"database/sql"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/lib/pq"
	"log"
	"time"
)

type Event struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	EventType string    `db:"event_type"`
	BeItmo    string    `db:"be_itmo"`
	Url       string    `db:"url"`
}

func (e Event) GetEventsByTag(tag string, db *sql.DB) ([]Event, error) {
	rows, err := db.Query("SELECT * FROM events WHERE be_itmo=$1", tag)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var events []Event
	err = sqlscan.ScanAll(&events, rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return events, nil
}

func (e Event) GetEventsByTagNoUser(tag string, db *sql.DB, idSlice []int) ([]Event, error) {
	rows, err := db.Query("SELECT * FROM events WHERE be_itmo=$1 and id != any($2) and start_date > current_date", tag, pq.Array(idSlice))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var events []Event
	err = sqlscan.ScanAll(&events, rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return events, nil
}
