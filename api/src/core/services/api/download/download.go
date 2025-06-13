package download

import (
	"bufio"
	"bytes"
	"context"
	"dse/src/core/log"
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils/cast"
	"dse/src/utils/datetime"
	"dse/src/utils/file"
	"dse/src/utils/gatekeeper"
	"dse/src/utils/httpio"
	"dse/src/utils/json"
	"dse/src/utils/object"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cohesivestack/valgo"
	"github.com/dromara/carbon/v2"
	"github.com/sourcegraph/conc"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type OutputUser struct {
	Token                string `json:"token"`

	Age                  string `json:"age"`
	Sex                  string `json:"sex"`
	Income               string `json:"income"`
	Resident             string `json:"resident"`
	Education            string `json:"education"`
	Political            string `json:"political"`
	Employment           string `json:"employment"`
	Postcode             string `json:"postcode"`

	SocialTv             bool `json:"social_tv"`
	SocialRadio          bool `json:"social_radio"`
	SocialAnders         bool `json:"social_anders"`
	SocialReddit         bool `json:"social_reddit"`
	SocialTwitter        bool `json:"social_twitter"`
	SocialYoutube        bool `json:"social_youtube"`
	SocialDeKrant        bool `json:"social_de-krant"`
	SocialFacebook       bool `json:"social_facebook"`
	SocialLinkedin       bool `json:"social_linkedin"`
	SocialTelegram       bool `json:"social_telegram"`
	SocialWhatsapp       bool `json:"social_whatsapp"`
	SocialInstagram      bool `json:"social_instagram"`
	SocialUnselected     bool `json:"social_unselected"`
	SocialNieuwswebsites bool `json:"social_nieuwswebsites"`

	BrowserBrave         bool `json:"browser_brave"`
	BrowserOpera         bool `json:"browser_opera"`
	BrowserChrome        bool `json:"browser_chrome"`
	BrowserSafari        bool `json:"browser_safari"`
	BrowserFirefox       bool `json:"browser_firefox"`
	BrowserUnselected    bool `json:"browser_unselected"`
	BrowserMicrosoftEdge bool `json:"browser_microsoft-edge"`

	LanguageDuits        bool `json:"language_duits"`
	LanguageFrans        bool `json:"language_frans"`
	LanguageEngels       bool `json:"language_engels"`
	LanguageSpaans       bool `json:"language_spaans"`
	LanguageItaliaans    bool `json:"language_italiaans"`
	LanguageNederlands   bool `json:"language_nederlands"`
	LanguageUnselected   bool `json:"language_unselected"`

	SearchEngineBing       bool `json:"search_engine_bing"`
	SearchEngineYahoo      bool `json:"search_engine_yahoo"`
	SearchEngineAnders     bool `json:"search_engine_anders"`
	SearchEngineEcosia     bool `json:"search_engine_ecosia"`
	SearchEngineGoogle     bool `json:"search_engine_google"`
	SearchEngineStartpage  bool `json:"search_engine_startpage"`
	SearchEngineDuckduckgo bool `json:"search_engine_duckduckgo"`
	SearchEngineUnselected bool `json:"search_engine_unselected"`
}

type OutputSearch struct {
	ID          int    `json:"id"`
	Token       string `json:"token"`
	Timestamp   string `json:"timestamp"`

	Url         string `json:"url"`
	Website     string `json:"website"`
	Keyword     string `json:"keyword"`

	Browser       map[string]any `json:"browser"`
	Localization  string          `json:"localization"`

	Form     *models.Form `json:"form"`
	Results  map[string]any `json:"results"`
}

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	gk = gatekeeper.NewGateKeeper(true)

	// Errors
	ErrInvalidToken         = fmt.Errorf("invalid token")
	ErrInvalidDays          = fmt.Errorf("invalid days parameter")
	ErrInternal             = fmt.Errorf("internal server error")
)

// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func parseTokenParameter(r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")

	if token == "dse2024" {
		return token, nil
	}

	return token, ErrInvalidToken
}

func parseDaysParameter(r *http.Request) (int, error) {
	days := r.URL.Query().Get("days")

	if days == "" {
		return 1000, nil
	}

	value, err := strconv.Atoi(days)
	if err != nil {
		return 0, ErrInvalidDays
	}

	if value < 0 {
		return 0, ErrInvalidDays
	}

	return value, nil
}

func parseStartEndParameters(r *http.Request) (time.Time, time.Time, error) {
    const layout0 = time.RFC3339
    const layout1 = "2006-01-02"

    start := r.URL.Query().Get("start")
    end   := r.URL.Query().Get("end")

    var timestart, timeend time.Time
    var err error

	if start == "" || end == "" {
		return time.Time{}, time.Time{}, nil
	}

    // Parse start time
    timestart, err = time.Parse(layout0, start)
    if err != nil {
        timestart, err = time.Parse(layout1, start)
        if err != nil {
            return time.Time{}, time.Time{}, fmt.Errorf("invalid start time format")
        }
        // Set to start of the day
        timestart = timestart.Truncate(24 * time.Hour)
    }

    // Parse end time
    timeend, err = time.Parse(layout0, end)
    if err != nil {
        timeend, err = time.Parse(layout1, end)
        if err != nil {
            return time.Time{}, time.Time{}, fmt.Errorf("invalid end time format")
        }
        // Set to end of the day
        timeend = timeend.Add(24*time.Hour - time.Second)
    }

    return timestart, timeend, nil
}

func parseTimestampFromLogLine(line string) (time.Time, error) {
	timestr := gjson.Get(line, "time").String()
	if timestr == "" {
		return time.Time{}, fmt.Errorf("missing time field")
	}

	timestamp, err := time.Parse(time.RFC3339, timestr)
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}


func loadUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)
	stream, err := db.StreamUsers()
	if err != nil { return nil, err }

	for user := range stream {
		if user == nil { break }
		users = append(users, user)
	}

	return users, nil
}

func loadUserMap() (map[string]*models.User, error) {
	users := make(map[string]*models.User, 0)
	stream, err := db.StreamUsers()
	if err != nil { return nil, err }

	for user := range stream {
		if user == nil { break }
		users[user.Token] = user
	}

	return users, nil
}

func transformSearch(search *models.Search, user *models.User) (*OutputSearch, error) {
	var output = &OutputSearch{}
	var skip = false

	func() {
		b, err := search.Metadata.Value()
		if err != nil {
			skip = true
		}

		var parsed = gjson.ParseBytes(b)

		output.ID        = int(search.ID)
		output.Token     = search.Token
		output.Timestamp = search.Timestamp

		output.Url     = parsed.Get("url").String()
		output.Website = parsed.Get("website").String()
		output.Keyword = parsed.Get("keyword").String()

		output.Browser, _ = parsed.Get("browser").Value().(map[string]any)

		output.Localization = parsed.Get("localization").String()
		output.Results, _   = parsed.Get("results").Value().(map[string]any)

		output.Form = user.State.Client.User.Form

		if output.Results == nil || len(output.Results) == 0 {
			skip = true
		}
	}()

	if skip {
		return nil, ErrInternal
	}

	return output, nil
}

func LoadData() {
	defer gk.Unlock()

	if file.Exists("searches.json") {
		t, err := file.GetModified("searches.json")
		if err != nil {
			log.Error().Err(err).Msg("Failed to get modified time")
			return
		}

		if time.Since(t) < 10 * time.Minute {
			return
		}
	}

	meta := struct {
		users map[string]*models.User // From current table
		forms map[string]*models.Form // From backup table
	}{}
	meta.users = make(map[string]*models.User)
	meta.forms = make(map[string]*models.Form)

	{ // Users
		var err error
		query     := `SELECT * FROM public.users`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to query users")
		}

		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.Token, &user.State)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to scan user")
			}

			meta.users[user.Token] = &user
		}
	}

	{ // Forms
		query := `SELECT * FROM public.users_2 as users`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to query users")
		}

		for rows.Next() {
			var err error
			var token string
			var data map[string]any
			var form models.Form

			err = rows.Scan(&token, &data)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to scan user")
			}

			o, has := object.Get(data, "client.user.form")
			if has == false {
				continue
			}

			err = object.ToStruct(o, &form)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to convert form")
			}

			if form.Sex == "" {
				continue
			}

			meta.forms[token] = &form
		}
	}

	{ // Merge
		for token, form := range meta.forms {
			user, has := meta.users[token]
			if has == false {
				continue
			}

			user.State.Client.User.Form.Sex = form.Sex
		}

		cusers := 0

		for _, user := range meta.users {
			if user.State.Client.User.Form.Sex != "" {
				cusers += 1
			}
		}
	}

	{ // Output
		var err error
		counter := 0

		f, err := os.Create("searches.json")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create output file")
		}
		defer f.Close()

		ctx := context.Background()
		query := `SELECT * FROM public.searches`
		ch, err := db.QueryWithLimit(ctx, query, 10000)

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to query searches")
		}

		for row := range ch {
			var ok bool
			var user *models.User
			search := &models.Search{}

			id, ok := row[0].(int64)
			if !ok {
				log.Fatal().Msg("Failed to convert row[0] to int64")
			}

			token, ok := row[1].(string)
			if !ok {
				log.Fatal().Msg("Failed to convert row[1] to string")
			}

			timestamp, ok := row[2].(time.Time)
			if !ok {
				log.Fatal().Msg("Failed to convert row[2] to time.Time")
			}

			metadata, ok := row[3].(map[string]any)
			if !ok {
				log.Fatal().Msg("Failed to convert row[3] to map[string]any")
			}

			search.ID        = uint64(id)
			search.Token     = token
			search.Timestamp = datetime.FromTime(timestamp).ToIso8601String()
			search.Metadata  = metadata

			user = meta.users[token]

			output, _ := sjson.Set("", "id", int(search.ID))
			output, _ = sjson.Set(output, "token", search.Token)
			output, _ = sjson.Set(output, "timestamp", search.Timestamp)

			b, err := search.Metadata.Value()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to convert metadata to bytes")
			}

			parsed := gjson.ParseBytes(b)

			output, _ = sjson.Set(output, "url", parsed.Get("url").String())
			output, _ = sjson.Set(output, "website", parsed.Get("website").String())
			output, _ = sjson.Set(output, "keyword", parsed.Get("keyword").String())

			browser, _ := parsed.Get("browser").Value().(map[string]any)
			output, _ = sjson.Set(output, "browser", browser)

			output, _ = sjson.Set(output, "localization", parsed.Get("localization").String())

			mapping, ok := parsed.Get("results").Value().(map[string]interface{})
			var results []interface{}
			if ok {
				results, _ = mapping["search_result"].([]interface{})
			}

			output, _ = sjson.Set(output, "results", results)
			output, _ = sjson.Set(output, "form", user.State.Client.User.Form)

			if results == nil || len(results) == 0 {
				continue
			}

			f.WriteString(fmt.Sprintf("%s\n", output))
			counter += 1
		}
	}
}

// ------------------------------------------------------------
// : Download Users
// ------------------------------------------------------------
func GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := parseTokenParameter(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users, err := loadUsers()
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	first := true

	for _, user := range users {
		b, err := json.ToBytes(user)
		if err != nil {
			log.Error().Err(err).Msg("")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if first {
			first = false
			w.Write([]byte(fmt.Sprintf("%s", b)))
		} else {
			w.Write([]byte(fmt.Sprintf(",\n%s", b)))
		}
	}
}

// ------------------------------------------------------------
// : Download Searches
// ------------------------------------------------------------
// http://localhost:5000/api/download/searches?token=dse2024&start=2024-11-01&end=2024-11-02
// https://static.33.56.161.5.clients.your-server.de/dse/api/download/searches?token=dse2024&days=1
// https://static.33.56.161.5.clients.your-server.de/dse/api/download/searches?token=dse2024&start=2024-11-01&end=2024-11-02
// https://static.33.56.161.5.clients.your-server.de/dse/api/download/searches?token=dse2024&start=2024-11-01T00:00:00Z&end=2024-11-02T00:00:00Z
func HandleDownloadSearches(w http.ResponseWriter, r *http.Request) {
    // Parse parameters
    _, err := parseTokenParameter(r)
    if err != nil {
        log.Error().Err(err).Msg("Invalid token")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    days, err := parseDaysParameter(r)
    if err != nil {
        log.Error().Err(err).Msg("Invalid days parameter")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    timestart, timeend, err := parseStartEndParameters(r)
    if err != nil {
        log.Error().Err(err).Msg("Invalid start or end parameter")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Set response headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)

    // Load user data
    users, err := loadUserMap()
    if err != nil {
        log.Error().Err(err).Msg("Failed to load users")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
    defer cancel()

    mutex   := &sync.Mutex{}
    offset  := 0
    limit   := 5_000
    channel := make(chan *models.Search, 50_000)

    go func() {
        defer close(channel)
        for {
            count := 0 // Reset count at the start of each loop iteration

            var params []interface{}
            param_index := 1

            var builder strings.Builder
            builder.WriteString("SELECT searches.*, users.* FROM searches LEFT JOIN users ON searches.token = users.token")

            // Build WHERE clause based on parameters
            if !timestart.IsZero() && !timeend.IsZero() {
                builder.WriteString(fmt.Sprintf(" WHERE timestamp BETWEEN $%d AND $%d", param_index, param_index+1))
                params = append(params, timestart, timeend)
                param_index += 2
            } else if days > 0 {
                builder.WriteString(fmt.Sprintf(" WHERE timestamp >= NOW() - INTERVAL '%d days'", days))
                // No additional parameters needed
            }

            // Add ORDER BY, OFFSET, and LIMIT clauses
            builder.WriteString(" ORDER BY timestamp DESC")
            builder.WriteString(fmt.Sprintf(" OFFSET $%d LIMIT $%d", param_index, param_index+1))
            params = append(params, offset, limit)

            query := builder.String()

            mutex.Lock()
            rows, err := db.QueryWithContext(ctx, query, params...)
            mutex.Unlock()
            if err != nil {
                log.Error().Err(err).Msg("Failed to query searches")
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            for rows.Next() {
                count++

                var user      models.User
                var search    models.Search
                var timestamp time.Time

                err = rows.Scan(
                    &search.ID,
                    &search.Token,
                    &timestamp,
                    &search.Metadata,
                    &user.Token,
                    &user.State,
                )
                if err != nil {
                    log.Error().Err(err).Msg("Failed to scan search")
                    return
                }

                search.Timestamp = datetime.FromTime(timestamp).ToIso8601String()

                select {
                case <-ctx.Done():
                    return
                case channel <- &search:
                }
            }

            rows.Close()

            if count < limit {
                break
            }
            offset += limit
        }
    }()

    for search := range channel {
        select {
        case <-ctx.Done():
            return
        default:
            user := users[search.Token]
            if user == nil {
                continue
            }

            output, err := transformSearch(search, user)
            if err != nil {
                continue
            }

            encoded, err := json.ToBytes(output)
            if err != nil {
                continue
            }

            w.Write(encoded)
            w.Write([]byte("\n"))
        }
    }
}

// http://localhost:5000/api/download/searches/merge?token=dse2024&start=2024-11-01&end=2024-11-02
func HandleDownloadSearches2(w http.ResponseWriter, r *http.Request) {
	request := struct {
		token   string
		start   string
		end     string
		days	string
	}{}

	request.token = r.URL.Query().Get("token")
	request.start = r.URL.Query().Get("start")
	request.end   = r.URL.Query().Get("end")
	request.days  = r.URL.Query().Get("days")
	
	validation := valgo.New()
	validation.Is(valgo.String(request.token, "token").EqualTo("dse2024"))

	switch {
		case request.days != "": {
			validation.Is(valgo.String(request.days, "days").GreaterThan("0").LessThan("10000"))
		}

		case request.start != "" && request.end != "": {
			fn := func(value string) bool { return datetime.Parse(value).IsZero() == false }
			validation.Is(valgo.String(request.start, "start").Passing(fn))
			validation.Is(valgo.String(request.end  , "end")  .Passing(fn))
		}

		default: {
			request.days = "9999"
		}
	}

	if validation.IsValid("token") == false { validation.AddErrorMessage("token", "Token is invalid") }
	if validation.IsValid("start") == false { validation.AddErrorMessage("start", "Start is invalid") }
	if validation.IsValid("end")   == false { validation.AddErrorMessage("end"  , "End is invalid")   }
	if validation.IsValid("days")  == false { validation.AddErrorMessage("days" , "Days is invalid")  }

	if validation.Valid() == false {
		httpio.WriteValidationError(w, validation)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	state := struct {
		Type  int

		Start carbon.Carbon
		End   carbon.Carbon
		Days  int

		Offset int
		Limit  int
		Rows   int
		
		Query string
		
		Users map[string]*models.User // From current table
		Forms map[string]*models.Form // From backup table

		Mutex *sync.Mutex
		Channel chan *models.Search
	}{}

	state.Mutex   = &sync.Mutex{}
	state.Channel = make(chan *models.Search, 50_000)

	state.Users = make(map[string]*models.User, 0)
	state.Forms = make(map[string]*models.Form, 0)

	state.Offset  = 0     // Current page
	state.Limit   = 5_000 // Items per page
	state.Rows    = 0

	switch {
		case request.start != "" && request.end != "": {
			state.Type  = 1
			state.Start = datetime.Parse(request.start)
			state.End   = datetime.Parse(request.end)
		}

		case request.days != "": {
			state.Type    = 2
			state.Days, _ = strconv.Atoi(request.days)
		}
	}

	wg := conc.NewWaitGroup()

	// Load users from current table
	wg.Go(func() {
		query := `
			SELECT 
				users.token, 
				users.state 
			FROM 
				public.users as users
		`

		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			user := &models.User{}
			err  := rows.Scan(&user.Token, &user.State)
			if err != nil {
				panic(err)
			}

			state.Users[user.Token] = user
		}
	})

	// Load users from backup table
	wg.Go(func() {
		query := `
			SELECT
				users.token,
				users.state
			FROM
				public.users_2 as users
		`

		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			usertoken := ""
			userstate := map[string]any{}
			err  := rows.Scan(&usertoken, &userstate)
			if err != nil {
				panic(err)
			}

			o, has := object.Get(userstate, "client.user.form")
			if has == false {
				continue
			}

			form := &models.Form{}
			
			b, err := json.ToBytes(o)
			json.FromBytes(b, form)

			if form.Age == "" {
				continue
			}

			state.Forms[usertoken] = form
		}
	})

	err := wg.WaitAndRecover().AsError()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load users")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Merge users
	wg.Go(func() {
		for token, form := range state.Forms {
			if form.Age == "" {
				continue
			}

			if user, exists := state.Users[token]; exists {
				user.State.Client.User.Form = form
			}
		}
	})

	err = wg.WaitAndRecover().AsError()
	if err != nil {
		log.Error().Err(err).Msg("Failed to merge users")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
    defer cancel()

	go func() {
		defer close(state.Channel)
		for {
			state.Rows = 0 // Reset row counter
			
			builder := strings.Builder{}
            builder.WriteString(`SELECT * FROM public.searches`)

			switch {
				// Start and End specified
				case state.Type == 1: { 
					start := state.Start.SetTimezone("UTC").ToIso8601ZuluString()
					end   := state.End.SetTimezone("UTC").ToIso8601ZuluString()
					builder.WriteString(fmt.Sprintf(" WHERE timestamp BETWEEN '%s' AND '%s'\n", start, end))
				}

				// Days specified
				case state.Type == 2: { 
					days := state.Days
					builder.WriteString(fmt.Sprintf(" WHERE timestamp >= NOW() - INTERVAL '%d days'\n", days))
				}
			}

			builder.WriteString(" ORDER BY timestamp DESC\n")
			builder.WriteString(fmt.Sprintf(" OFFSET %d LIMIT %d", state.Offset, state.Limit))

			state.Query = builder.String()

			println(state.Query)

			state.Mutex.Lock()
			rows, err := db.QueryWithContext(ctx, state.Query)
			state.Mutex.Unlock()

			for rows.Next() {
				state.Rows += 1

				var search models.Search
				var timestamp time.Time

				err = rows.Scan(
					&search.ID,
					&search.Token,
					&timestamp,
					&search.Metadata,
				)

				if err != nil {
					log.Error().Err(err).Msg("Failed to scan search")
					return
				}

				search.Timestamp = datetime.FromTime(timestamp).ToIso8601String()

				select {
					case <-ctx.Done(): return
					case state.Channel <- &search:
				}
			}
			rows.Close()

			if state.Rows < state.Limit {
				break
			}

			state.Offset += state.Limit
		}
	}()

	for search := range state.Channel {
        select {
        case <-ctx.Done():
            return
        default:
			user := state.Users[search.Token]
            if user == nil {
                continue
            }

            output, err := transformSearch(search, user)
            if err != nil {
                continue
            }

            encoded, err := json.ToBytes(output)
            if err != nil {
                continue
            }

            w.Write(encoded)
            w.Write([]byte("\n"))
        }
    }
}


// ------------------------------------------------------------
// : Download Log File
// ------------------------------------------------------------
func GetLogs(w http.ResponseWriter, r *http.Request) {
	// Parse the 'hours' parameter
	param_hours := r.URL.Query().Get("hours")
	if param_hours == "" {
		param_hours = "24000" // Default to 1000 days * 24 hours = 24000 hours
	}

	hours, err := strconv.Atoi(param_hours)
	if err != nil || hours < 0 {
		log.Error().Err(err).Msg("Invalid hours parameter")
		http.Error(w, "Invalid hours parameter", http.StatusBadRequest)
		return
	}

	// Open the log file
	file, err := os.Open("./logs/api.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to open log file")
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set response headers for streaming JSON
	w.Header().Set("Content-Type"     , "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control"    , "no-cache")
	w.WriteHeader(http.StatusOK)

	// Calculate the cutoff time based on hours
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)

	// Start the JSON array
	w.Write([]byte("[\n"))

	// Initialize a buffer to collect log lines
	var lines []string

	// Use a buffer to read chunks from the end
	const buffersize = 4096
	fileinfo, err := file.Stat()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file info")
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	filesize := fileinfo.Size()
	var offset int64 = 0
	var remaining = filesize
	var partialline bytes.Buffer

	for remaining > 0 {
		if remaining < buffersize {
			offset = 0
		} else {
			offset = remaining - buffersize
		}

		chunksize := buffersize
		if remaining < buffersize {
			chunksize = int(remaining)
		}

		// Read the chunk
		buf := make([]byte, chunksize)
		_, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			log.Error().Err(err).Msg("Error reading log file")
			break
		}

		// Process the chunk in reverse
		for i := chunksize - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				// If there's a partial line collected, append it
				if partialline.Len() > 0 {
					line := reverseBytes(partialline.Bytes())
					lineStr := string(line)
					partialline.Reset()

					// Prepend the line to lines
					lines = append(lines, lineStr)
				}
			} else {
				partialline.WriteByte(buf[i])
			}
		}

		remaining -= int64(chunksize)
	}

	// Add the last line if any
	if partialline.Len() > 0 {
		line := reverseBytes(partialline.Bytes())
		lineStr := string(line)
		lines = append(lines, lineStr)
	}

	// Now, iterate over the collected log lines
	first := true
	for _, line := range lines {
		// Parse the timestamp
		logTime, err := parseTimestampFromLogLine(line)
		if err != nil {
			// Skip lines with invalid or missing timestamps
			continue
		}

		// Check if the log is within the specified timeframe
		if logTime.Before(cutoff) {
			// Since we're reading from most recent to oldest, we can stop here
			break
		}

		// Write a comma before each log entry except the first
		if !first {
			w.Write([]byte(",\n"))
		} else {
			first = false
		}

		// Write the JSON log entry
		w.Write([]byte(line))

		// Flush the response to send data to the client immediately
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	// End the JSON array
	w.Write([]byte("\n]"))

	// Final flush to ensure all data is sent
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// http://localhost/api/download/searches?days=2
// http://localhost/api/download/searches?start=2025-01-28
// http://localhost/api/download/searches?start=2025-01-28&end=2025-01-29
// https://static.33.56.161.5.clients.your-server.de/dse/api/download/searches?days=2
// https://static.33.56.161.5.clients.your-server.de/dse/api/download/searches/full?days=2
func GetSearches(w http.ResponseWriter, r *http.Request) {
	gk.Wait()

	request := struct {
		start string
		end   string
		days  string

		kind string
	}{}

	request.kind  = ""
	request.days  = r.URL.Query().Get("days")
	request.start = r.URL.Query().Get("start")
	request.end   = r.URL.Query().Get("end")

	v := valgo.New()

	switch {
		case request.days != "" && (request.start != "" || request.end != ""):
			v.AddErrorMessage("parameters", "Cannot specify both days and start/end")

		case request.days != "":
			d := cast.ToInt(request.days)
			v.Is(valgo.Int(d, "days").GreaterThan(0).LessThan(10000))
			request.kind = "days"

		case request.start != "" && request.end != "":
			fn := func(value string) bool { return datetime.Parse(value).IsZero() == false }
			v.Is(valgo.String(request.start, "start").Passing(fn))
			v.Is(valgo.String(request.end, "end").Passing(fn))
			request.kind = "start_end"

		case request.start != "":
			fn := func(value string) bool { return datetime.Parse(value).IsZero() == false }
			v.Is(valgo.String(request.start, "start").Passing(fn))
			request.end = datetime.Now().ToIso8601String()
			request.kind = "start"

		case request.end != "":
			v.AddErrorMessage("parameters", "Start must be specified if end is provided")

		default:
			request.days = "9999"
			request.kind = "days"
	}

	if v.IsValid("days")  == false { v.AddErrorMessage("days" , `days is invalid`)  }

	if v.Valid() == false {
		httpio.WriteValidationError(w, v)
		return
	}

	f, err := file.Open("searches.json")
	if err != nil {
		httpio.WriteError(w, http.StatusInternalServerError, "File not found")
		return 
	}
	defer f.Close()

	ctx  := r.Context()
	buf  := bytes.Buffer{}
	size := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			line      := scanner.Text()
			timestamp := json.Get(line, "timestamp").String()
			var skip bool

			switch request.kind {
			case "days":
				cutoff := datetime.Now().SubDays(cast.ToInt(request.days))
				skip    = cutoff.Gte(datetime.Parse(timestamp))
			case "start":
				start := datetime.Parse(request.start)
				skip   = start.Gte(datetime.Parse(timestamp))
			case "start_end":
				start := datetime.Parse(request.start)
				end   := datetime.Parse(request.end)
				skip   = start.Gte(datetime.Parse(timestamp)) || end.Lte(datetime.Parse(timestamp))
			}

			if skip {
				break
			}

			line = fmt.Sprintf("%s\n", line)
			buf.WriteString(line)
			size++

			if size >= 10000 {
				w.Write(buf.Bytes())
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				buf.Reset()
				size = 0
			}
		}
	}

	if buf.Len() > 0 {
		w.Write(buf.Bytes())
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
