package db

import (
	"context"
	"dse/src/core/models"
	"dse/src/utils"
	"dse/src/utils/arraylist"
	"dse/src/utils/cmap"
	"dse/src/utils/event"
	"dse/src/utils/gatekeeper"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type User   = models.User
type Search = models.Search

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func init() {
	
}

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger() // Initialize logger for DB operations

	// Database connection parameters with default values
	dsn  string = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	host string = "localhost"
	port string = "5432"
	user string = "postgres"
	pass string = "postgres"
	name string = "postgres"

	pool *pgxpool.Pool // PGX Pool instance

	gk    = gatekeeper.NewGateKeeper(true) // Locked on initialization
	users = cmap.New[*User]()
	ready = make(chan struct{})            // Channel to signal readiness

	mutex = &sync.Mutex{} // Write mutex

	ch_update = make(chan *User, 1000) // Channel for user updates

	// Errors
	ErrMissingToken   = errors.New("missing token")        // Error for missing token
	ErrTokenInvalid   = errors.New("invalid token")        // Error for invalid token
	ErrUserNotFound   = errors.New("user not found")       // Error when user is not found
	ErrTokenEmpty     = errors.New("token is empty")       // Error when token is empty
	ErrTokenIncorrect = errors.New("token length invalid") // Error for incorrect token length
)

// ------------------------------------------------------------
// : Methods	
// ------------------------------------------------------------
func Wait() {
	gk.Wait()
}
// ------------------------------------------------------------
// : Users
// ------------------------------------------------------------
func HasUser(token string) (bool, error) {
    Wait()
	if token == ""      { return false, ErrTokenEmpty }
	if len(token) != 12 { return false, ErrTokenIncorrect }

    var count int64
    query := `SELECT COUNT(*) FROM users WHERE token = $1`

	mutex.Lock()
	defer mutex.Unlock()

    err := pool.QueryRow(context.Background(), query, token).Scan(&count)
    if err != nil {
		return false, err
    }
	return count > 0, nil
}


func CreateUser(token string) (*User, error) {
	Wait()

	if token == ""       { return nil, ErrTokenEmpty }
	if len(token) != 12  { return nil, ErrTokenIncorrect }

	exists, err := HasUser(token)
	if err != nil { return nil, err }

	if exists { return nil, fmt.Errorf("User with token %s already exists", token) }
	
	mutex.Lock()
	defer mutex.Unlock()

	var user = &User{Token: token}
	user.Init()

	var query = `INSERT INTO users (token, state) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING`
	_, err = pool.Exec(context.Background(), query, token, user.State)
	if err != nil { return nil, err }

	event.Emit("user.created", user)

	users.Set(token, user)
	return user, nil
}

func UpdateUser(user *User) (*User, error) {
	Wait()
	defer recover()

	mutex.Lock()
	defer mutex.Unlock()

	var query = `UPDATE users SET state = $1 WHERE token = $2`
	_, err := pool.Exec(context.Background(), query, user.State, user.Token)
	if err != nil { return nil, err }

	return user, nil
}

func GetUser(token string) (*User, error) {
	Wait()
	if token == ""    	{ return nil, ErrTokenEmpty }
	if len(token) != 12 { return nil, ErrTokenIncorrect }

	if users.Has(token) {
		user, _ := users.Get(token)
		return user, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	rows, err := pool.Query(context.Background(), `SELECT token, state FROM users WHERE token = $1`, token)
	if err != nil { return nil, err }

	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Token, &user.State)
		if err != nil { return nil, err }

		user.Init()
		users.Set(user.Token, &user)

		return &user, nil
	}
	return nil, ErrUserNotFound
}

func GetUsers() (*arraylist.ArrayList[*User], error) {
// func GetUsers() ([]*User, error) {
	Wait()

	list := arraylist.NewArrayList[*User]()
	for _, user := range users.Items() {
		list.Add(user)
	}
	return list, nil

	// slice := make([]*User, 0, len(users.Items()))
	// for _, user := range users.Items() {
	// 	slice = append(slice, user)
	// }
	// return slice, nil
}

func StreamUsers() (<-chan *User, error) {
	Wait()
	var channel = make(chan *User, 1000)

	go func() {
		defer close(channel)
		for _, user := range users.Items() {
			channel <- user
		}
	}()

	return channel, nil
}
// ------------------------------------------------------------
// : Search
// ------------------------------------------------------------

func StreamSearches(ctx context.Context, days int) (<-chan *Search, error) {
	Wait()

	var offset  = 0 
	var limit   = 5_000
	var channel = make(chan *Search, 50_000)

	go func() {
		defer close(channel)

		var after = carbon.Now(carbon.UTC).SubDays(days).StdTime()
		var query = `
		SELECT searches.*, users.* 
		FROM   searches LEFT JOIN users ON searches.token = users.token 
		WHERE  searches.timestamp >= $1 
		ORDER  BY searches.timestamp DESC 
		OFFSET $2 
		LIMIT  $3
		`

		for {
			mutex.Lock()
			rows, err := pool.Query(context.Background(), query, after, offset, limit)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to stream searches")
				return
			}
			mutex.Unlock()

			count := 0
			for rows.Next() {
				count++
				var user   User
				var search Search
				var timestamp time.Time // Fixes: Cannot scan into string from binary format

				err = rows.Scan(
					&search.ID,
					&search.Token,
					&timestamp,
					&search.Metadata,
					&user.Token,
					&user.State,
				)
				search.Timestamp = carbon.CreateFromTimestamp(timestamp.Unix()).ToIso8601String()

				if err != nil {
					logger.Error().Err(err).Msg("Failed to scan search")
					return
				}

				select {
					case <-ctx.Done(): return
					case channel <- &search:
				}
			}

			rows.Close()

			if count < limit { break }
			
			offset += limit
		}
	}()

	return channel, nil
}

func CreateSearch(search *Search) (*Search, error) {
    Wait()

	if search.Token == ""      { return nil, ErrMissingToken }
	if len(search.Token) != 12 { return nil, ErrTokenInvalid }

    var err error
    var query = `INSERT INTO searches (token, timestamp, metadata) VALUES ($1, $2, $3)`

    metadata, err := json.Marshal(search.Metadata)
	if err != nil {
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()

    _, err = pool.Exec(context.Background(), query, search.Token, search.Timestamp, metadata)
    if err != nil { 
		return nil, err 
	}

	return search, nil
}

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func GetConnection() *pgxpool.Pool {
	Wait()
	return pool
}

func Query(query string, args ...interface{}) (pgx.Rows, error) {
	Wait()
	return pool.Query(context.Background(), query, args...)
}

func QueryWithContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	Wait()
	return pool.Query(ctx, query, args...)
}

func StreamQuery(ctx context.Context, query string, args ...interface{}) (<-chan map[string]interface{}, error) {
	Wait()
	
	const limit = 10_000
	results := make(chan map[string]interface{}, limit)
	
	go func() {
		defer close(results)
		
		offset := 0
		for {
			// Add pagination to query
			query := fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
			
			mutex.Lock()
			rows, err := pool.Query(ctx, query, args...)
			mutex.Unlock()
			
			if err != nil {
				logger.Error().Err(err).Msg("Query failed")
				return
			}
			
			count := 0
			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					logger.Error().Err(err).Msg("Failed to read row values") 
					rows.Close()
					return
				}
				
				// Build row map
				data := make(map[string]interface{}, len(rows.FieldDescriptions()))
				for i, col := range rows.FieldDescriptions() {
					data[string(col.Name)] = values[i]
				}
				
				select {
				case <-ctx.Done():
					rows.Close()
					return
				case results <- data:
					count++
				}
			}
			rows.Close()
			
			// Exit if we got less than limit
			if count < limit {
				break
			}
			
			offset += limit
		}
	}()

	return results, nil
}


func QueryWithLimit(ctx context.Context, query string, limit int, args ...interface{}) (<-chan []interface{}, error) {
    Wait()

    if limit <= 0 {
        return nil, fmt.Errorf("limit must be greater than 0")
    }

    out := make(chan []interface{}, limit)

    go func() {
        defer close(out)
        offset := 0

        for {
            // Clone the args and append limit and offset in the correct order
            page := append([]interface{}{}, args...)
            page = append(page, limit, offset) // Append limit first, then offset

            // Construct the paginated query with ORDER BY before LIMIT and OFFSET
            paginatedQuery := fmt.Sprintf("%s ORDER BY id DESC LIMIT $%d OFFSET $%d", query, len(page)-1, len(page))

            mutex.Lock()
            rows, err := pool.Query(ctx, paginatedQuery, page...)
            mutex.Unlock()
            if err != nil {
                logger.Error().Err(err).Msg("Failed to query with limit")
                return
            }

            count := 0
            for rows.Next() {
                values, err := rows.Values()
                if err != nil {
                    logger.Error().Err(err).Msg("Failed to read row values")
                    rows.Close()
                    return
                }

                select {
                case <-ctx.Done():
                    rows.Close()
                    return
                case out <- values:
                    count++
                }
            }
            rows.Close()

            // If we got fewer results than limit, we've reached the end
            if count < limit {
                break
            }
            offset += limit
        }
    }()

    return out, nil
}

// ------------------------------------------------------------
// : Internals
// ------------------------------------------------------------
func InitUsers() {
	var err error
	var query = `SELECT * FROM users WHERE token IS NOT NULL`
	var rows pgx.Rows

	rows, err = pool.Query(context.Background(), query)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load users")
		return
	}

	for rows.Next() {
		var user User
	
		err = rows.Scan(&user.Token, &user.State)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan user")
			return
		}
		user.Init()
		users.Set(user.Token, &user)
	}
}

func InitListeners() {
	go func() {
		ch := event.On("user.updated")

		for e := range ch {
			user, ok := e.Args[0].(*User)
			if !ok { continue }

			ch, ok := e.Args[1].(chan error)
			if !ok { continue }

			_, err := UpdateUser(user)
			ch <- err
			close(ch)
		}
	}()
}

func health() {
	var err error
	go func() {
		for {
			_, err = pool.Exec(context.Background(), "SELECT 1")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to ping database")
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()
	select {}
}
// ------------------------------------------------------------
// : Start
// ------------------------------------------------------------
func Start() {
	var err error

	if value, ok := os.LookupEnv("DB_HOST"); ok { host = value }
	if value, ok := os.LookupEnv("DB_PORT"); ok { port = value }
	if value, ok := os.LookupEnv("DB_USER"); ok { user = value }
	if value, ok := os.LookupEnv("DB_PASS"); ok { pass = value }
	if value, ok := os.LookupEnv("DB_NAME"); ok { name = value }

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, pass, name, port)
	
	logger.Debug().Str("dsn", dsn).Msg("Connecting to database")

	config, err := pgxpool.ParseConfig(dsn)
	config.MaxConns = 20
	config.MinConns = 1
	config.MaxConnIdleTime   = 5 * time.Minute
	config.HealthCheckPeriod = 5 * time.Second
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse database connection string")
		return
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		return
	}

	user_table := `
	CREATE TABLE IF NOT EXISTS users (
		token   VARCHAR(12) PRIMARY KEY,
		state   JSONB
	);`

	search_table := `
	CREATE TABLE IF NOT EXISTS searches (
		id          BIGSERIAL PRIMARY KEY,
		token       VARCHAR(12),
		timestamp   TIMESTAMP,
		metadata    JSONB
	);`

	metric_table := `
	CREATE TABLE IF NOT EXISTS metrics (
		id          BIGSERIAL PRIMARY KEY,
		timestamp   TIMESTAMP,
		metric      JSONB
	);`

	_, err = pool.Exec(context.Background(), user_table)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create users table")
	}

	_, err = pool.Exec(context.Background(), search_table)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create searches table")
	}

	_, err = pool.Exec(context.Background(), metric_table)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create metrics table")
	}

	go health()
	go InitListeners()
	
	InitUsers()
	
	logger.Info().Msg("Ready")
	gk.Unlock()
}
