package pgsender

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
)

var (
	uri = os.Getenv("pg_url")
	min = os.Getenv("pg_min")
	max = os.Getenv("pg_max")
)

// Handle errors
func Handle(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Print(err)
		panic(err)
	}
}

func waitForNotification(l *pq.Listener) {
	select {
	case <-l.Notify:
		log.Println("New work available")
		return
	case <-time.After(90 * time.Second):
		go l.Ping()
		// Check if there's more work available, just in case it takes
		// a while for the Listener to notice connection loss and
		// reconnect.
		log.Println("received no work for 90 seconds, checking for new work")
	}
}

type workChecker func(db *sql.DB)

// Connect sms.request
func Connect(event string, check workChecker) {
	db, err := sql.Open("postgres", uri)
	Handle(err, "Error opening db")
	_, err = db.Exec("set role sender")
	Handle(err, "Error swithing to sender role")
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	}
	min, err := strconv.Atoi(min)
	max, err := strconv.Atoi(max)
	minReconn := time.Duration(min) * time.Second
	maxReconn := time.Duration(max) * time.Minute
	listener := pq.NewListener(uri, minReconn, maxReconn, reportProblem)
	err = listener.Listen(event)
	Handle(err, "Error listening to db event:"+event)
	for {
		// process all available work before waiting for notifications
		check(db)
		waitForNotification(listener)
	}

}
