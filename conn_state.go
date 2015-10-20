package greene

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func New(duration time.Duration) func(net.Conn, http.ConnState) {
	// Keep a map of conn => time relationships to track
	// the age of each connection.
	mapping := make(map[net.Conn]time.Time)

	// Lock to guard against mutation of `mapping`.
	mutex := &sync.Mutex{}

	// Set `conn` to `mapping` with current time.
	set := func(conn net.Conn) {
		mutex.Lock()
		defer mutex.Unlock()
		mapping[conn] = time.Now()
	}

	// Remove `conn` from mapping.
	del := func(conn net.Conn) {
		mutex.Lock()
		defer mutex.Unlock()
		delete(mapping, conn)
	}

	return func(conn net.Conn, state http.ConnState) {
		switch state {

		case http.StateNew:
			// Store each new `conn` in our map with
			// the current time.
			set(conn)

		case http.StateIdle:
			// Close `conn` if its age exceeds the configured
			// duration.
			if t, ok := mapping[conn]; ok {
				if time.Now().Sub(t) >= duration {
					del(conn)
					conn.Close()
				}
			}

		case http.StateClosed, http.StateHijacked:
			// Save memory and discard disconnected clients
			del(conn)
		}
	}
}
