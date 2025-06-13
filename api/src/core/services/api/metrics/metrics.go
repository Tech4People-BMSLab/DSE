package metrics

import (
	"context"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/datetime"
	"dse/src/utils/hashmap"
	"dse/src/utils/json"
	"net/http"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/tidwall/sjson"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type Map = map[string]any

type Time      = time.Time
type Duration  = time.Duration
type Carbon    = carbon.Carbon
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
	hm     = hashmap.NewHashMap[string, []Map]()

	time_begin = datetime.ToTime(datetime.Now().SubHours(24))
	time_end   = datetime.ToTime(datetime.Now())
)
// ------------------------------------------------------------
// : Internals
// ------------------------------------------------------------
func write(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func execMetricUsers() {
	pool  := db.GetConnection()
	query := `
	SELECT
	to_timestamp(FLOOR(EXTRACT(EPOCH FROM timestamp) / 10) * 10) AS time_bucket,
	CAST(AVG((metric->>'all')::integer) AS integer)              AS avg_all,
	CAST(AVG((metric->>'connected')::integer) AS integer)        AS avg_connected,
	CAST(AVG((metric->>'disconnected')::integer) AS integer)     AS avg_disconnected,
	CAST(AVG((metric->>'consented')::integer) AS integer)        AS avg_consented
	FROM public.metrics
	WHERE
	metric->>'type' = 'users' AND
	metric->>'version' = '1' AND
	timestamp BETWEEN $1 AND $2
	GROUP BY time_bucket
	ORDER BY time_bucket DESC
	`

	ctx       := context.Background()
	rows, err := pool.Query(ctx, query, time_begin, time_end)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load users")
		return
	}
	defer rows.Close()
	if rows.Err() != nil {
		logger.Error().Err(err).Msg("Failed to load users")
		return
	}

	var metrics []Map
	for rows.Next() {
		var r1 Time
		var r2 int64
		var r3 int64
		var r4 int64
		var r5 int64
		var m  Map = make(Map)

		err := rows.Scan(&r1, &r2, &r3, &r4, &r5)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan metric")
			return
		}

		m["timestamp"]    = r1
		m["all"] 	      = int64(r2)
		m["connected"]    = int64(r3)
		m["disconnected"] = int64(r4)
		m["consented"]    = int64(r5)

		metrics = append(metrics, m)
	}

	hm.Set("users", metrics)
}

func execMetricSearch() {
	pool  := db.GetConnection()
	query := `
	SELECT
	to_timestamp(FLOOR(EXTRACT(EPOCH FROM timestamp) / 10) * 10) AS time_bucket,
	CAST(AVG((metric->>'count')::integer)                        AS integer) AS avg_count
	FROM public.metrics
	WHERE
	metric->>'type' = 'searches' AND
	metric->>'version' = '1' AND
	timestamp BETWEEN $1 AND $2
	GROUP BY time_bucket
	ORDER BY time_bucket DESC
	`

	ctx       := context.Background()
	rows, err := pool.Query(ctx, query, time_begin, time_end)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load searches")
		return
	}
	defer rows.Close()
	
	var metrics []Map
	for rows.Next() {
		var r1 Time  // Timestamp
		var r2 int64 // Count
		var m  Map = make(Map)

		err := rows.Scan(&r1, &r2)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan metric")
			return
		}

		m["timestamp"] = r1
		m["count"]     = r2

		metrics = append(metrics, m)
	}

	hm.Set("searches", metrics)
}

func execMetricSearchSize() {
	pool  := db.GetConnection()
	query := `
	SELECT
	to_timestamp(FLOOR(EXTRACT(EPOCH FROM timestamp) / 10) * 10) AS time_bucket,
	CAST(AVG((metric->>'size')::integer) AS integer)             AS avg_size
	FROM public.metrics
	WHERE
	metric->>'type'    = 'searches_size' AND
	metric->>'version' = '1'             AND
	timestamp BETWEEN $1 AND $2
	GROUP BY time_bucket
	ORDER BY time_bucket DESC
	`

	ctx       := context.Background()
	rows, err := pool.Query(ctx, query, time_begin, time_end)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load searches size")
		return
	}
	defer rows.Close()

	var metrics []Map
	for rows.Next() {
		var r1 Time  // Timestamp
		var r2 int64 // Size
		var m  Map = make(Map)

		err := rows.Scan(&r1, &r2)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan metric")
			return
		}

		m["timestamp"] = r1
		m["size"]      = r2

		metrics = append(metrics, m)
	}

	hm.Set("searches_size", metrics)
}

func execMetricSearchTotal() {
	pool  := db.GetConnection()
	query := `
	SELECT
	to_timestamp(FLOOR(EXTRACT(EPOCH FROM timestamp) / 10) * 10) AS time_bucket,
	CAST(SUM((metric->>'count')::integer) AS integer)            AS count
	FROM public.metrics
	WHERE
	metric->>'type'    = 'searches_total' AND
	metric->>'version' = '1'              AND
	timestamp BETWEEN $1 AND $2
	GROUP BY time_bucket
	ORDER BY time_bucket DESC
	`

	ctx       := context.Background()
	rows, err := pool.Query(ctx, query, time_begin, time_end)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load searches total")
		return
	}
	defer rows.Close()

	var metrics []Map
	for rows.Next() {
		var r1 Time  // Timestamp
		var r2 int64 // Count
		var m  Map = make(Map)

		err := rows.Scan(&r1, &r2)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan metric")
			return
		}

		m["timestamp"] = datetime.FromTime(r1)
		m["count"]     = r2

		metrics = append(metrics, m)
	}

	hm.Set("searches_total", metrics)
}

// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	payload, _ := sjson.Set("", "status", "OK")
	write(w, payload)
}

func GetMetricUsers(w http.ResponseWriter, r *http.Request) {
	if !hm.Has("users") {
		execMetricUsers()
	}

	arr  := hm.MustGet("users")
	b, _ := json.ToBytes(arr)
	write(w, string(b))
}

func GetMetricSearch(w http.ResponseWriter, r *http.Request) {
	if !hm.Has("searches") {
		execMetricSearch()
	}
	arr  := hm.MustGet("searches")
	b, _ := json.ToBytes(arr)
	write(w, string(b))
}

func GetMetricsSearchSize(w http.ResponseWriter, r *http.Request) {
	if !hm.Has("searches_size") {
		execMetricSearchSize()
	}
	arr    := hm.MustGet("searches_size")
	b  , _ := json.ToBytes(arr)
	write(w, string(b))
}

func GetMetricsSearchTotal(w http.ResponseWriter, r *http.Request) {
	if !hm.Has("searches_total") {
		execMetricSearchTotal()
	}
	arr    := hm.MustGet("searches_total")
	b  , _ := json.ToBytes(arr)
	write(w, string(b))
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func Init() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			time_begin = datetime.ToTime(datetime.Now().SubHours(24))
			time_end   = datetime.ToTime(datetime.Now())

			execMetricUsers()
			execMetricSearch()
			execMetricSearchSize()
			execMetricSearchTotal()
		}
	}
}
