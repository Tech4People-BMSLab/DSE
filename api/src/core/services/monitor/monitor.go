package monitor

import (
	"context"
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/datetime"
	"dse/src/utils/event"
	"dse/src/utils/hashmap"
	"time"

	"golang.org/x/time/rate"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type Map = map[string]any

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
)
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func save(timestamp string, data *hashmap.HashMap[string, any]) {
	b, _ := data.ToJSON()
	pool := db.GetConnection()
	pool.Exec(context.Background(), `INSERT INTO metrics (timestamp, metric) VALUES ($1, $2)`, timestamp, b)
}

// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func MonitorUsers() {
	users, err := db.GetUsers()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get users")
		return
	}

	m := hashmap.NewHashMap[string, any]()
	m.Set("version",      "1")
	m.Set("type",         "users")
	m.Set("all",          int64(users.Len()))
	m.Set("consented",    int64(0))
	m.Set("connected",    int64(0))
	m.Set("disconnected", int64(0))

	users.Each(func(index int, user *models.User) bool {
		if user.State.Client.User.Form.Age != "" {
			m.Set("consented", m.MustGet("consented").(int64)+1)
		}

		if user.State.Server.Online {
			m.Set("connected", m.MustGet("connected").(int64)+1)
		} else {
			m.Set("disconnected", m.MustGet("disconnected").(int64)+1)
		}
		return true
	})

	save(datetime.ToISO(datetime.Now()), m)
}

// ------------------------------------------------------------
// : Searches
// ------------------------------------------------------------
func MonitorSearches() {
	pool       := db.GetConnection()
	week_start := datetime.ToTime(datetime.StartOfWeek())
	week_end   := datetime.ToTime(datetime.EndOfWeek())

	query := `
	SELECT COUNT(*) FROM public.searches
	WHERE timestamp >= $1 AND timestamp < $2
	`

	ctx   := context.Background()
	row   := pool.QueryRow(ctx, query, week_start, week_end)
	count := int64(0)

	err := row.Scan(&count)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute query")
		return
	}

	m := hashmap.NewHashMap[string, any]()
	m.Set("version", "1")
	m.Set("type"   , "searches")
	m.Set("count"  , int64(count))

	save(datetime.ToISO(datetime.Now()), m)
}

func MonitorSearchesSize() {
	pool  := db.GetConnection()
	query := `SELECT pg_total_relation_size('public.searches') AS table_size`

	ctx  := context.Background()
	row  := pool.QueryRow(ctx, query)
	size := int64(0)

	err := row.Scan(&size)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute query")
		return
	}

	m := hashmap.NewHashMap[string, any]()
	m.Set("version", "1")
	m.Set("type"   , "searches_size")
	m.Set("size"   , int64(size))

	save(datetime.ToISO(datetime.Now()), m)
}

func MonitorSearchesTotal() {
	pool  := db.GetConnection()
	query := `SELECT COUNT(*) FROM public.searches;`

	ctx   := context.Background()
	row   := pool.QueryRow(ctx, query)
	count := int64(0)

	err := row.Scan(&count)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute query")
		return
	}

	m := hashmap.NewHashMap[string, any]()
	m.Set("version", "1")
	m.Set("type"   , "searches_total")
	m.Set("count"  , int64(count))

	save(datetime.ToISO(datetime.Now()), m)
}


// ------------------------------------------------------------
// : Monitor
// ------------------------------------------------------------
func Init() {
	go func() {
		limiter := rate.NewLimiter(rate.Every(5 * time.Second), 1)
		ch      := event.On(event.UserConnected)
		for range ch { 
			if limiter.Allow() {
				go MonitorUsers()
			}
		}	
	}()

	go func() {
		limiter := rate.NewLimiter(rate.Every(5 * time.Second), 1)
		ch      := event.On(event.UserDisconnected)
		for range ch { 
			if limiter.Allow() {
				go MonitorUsers() 
			}
		}	
	}()

	go func() {
		go MonitorSearches()

		limiter := rate.NewLimiter(rate.Every(5 * time.Second), 1)
		ch      := event.On(event.ExtractorItemDone)
		for range ch { 
			if limiter.Allow() {
				go MonitorSearches() 
			}
		}
	}()

	go func() {
		go MonitorSearchesSize()

		limiter := rate.NewLimiter(rate.Every(5 * time.Second), 1)
		ch      := event.On(event.ExtractorItemDone)
		for range ch { 
			if limiter.Allow() {
				go MonitorSearchesSize()
			}
		}
	}()
}
