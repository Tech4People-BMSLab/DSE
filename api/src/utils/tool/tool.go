package tool

import (
	"context"
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/datetime"
	"dse/src/utils/object"
	"fmt"
	"os"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ------------------------------------------------------------
// : Alias
// ------------------------------------------------------------
type Packet = models.Packet
type State  = models.State
type User   = models.User
type Task   = models.Task
type Form   = models.Form
type Search = models.Search

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger  = utils.NewLogger()
)

func Prepare() {
	meta := struct {
		Users map[string]*models.User // From current table
		Forms map[string]*models.Form // From backup table
	}{}
	meta.Users = make(map[string]*models.User)
	meta.Forms = make(map[string]*models.Form)

	goto USERS

	USERS:
	{
		var err error
		query := `SELECT * FROM public.users`
		rows, err := db.Query(query)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to query users")
		}

		for rows.Next() {
			var user User
			err := rows.Scan(&user.Token, &user.State)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to scan user")
			}

			meta.Users[user.Token] = &user
		}

		goto FORMS
	}

	FORMS:
	{
		query := `SELECT * FROM public.users_2 as users`
		rows, err := db.Query(query)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to query users")
		}

		for rows.Next() {
			var err error
			var token string
			var data map[string]any
			var form models.Form

			err = rows.Scan(&token, &data)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to scan user")
			}

			o, has := object.Get(data, "client.user.form")
			if has == false {
				continue
			}

			err = object.ToStruct(o, &form)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to convert form")
			}

			if form.Sex == "" {
				continue
			}

			meta.Forms[token] = &form
		}

		goto MERGE
	}

	MERGE:
	{
		for token, form := range meta.Forms {
			user, has := meta.Users[token]
			if has == false {
				continue
			}

			user.State.Client.User.Form.Sex = form.Sex
		}

		cusers := 0

		for _, user := range meta.Users {
			if user.State.Client.User.Form.Sex != "" {
				cusers += 1
			}
		}

		goto OUTPUT
	}

	OUTPUT:
	{
		var err error
		counter := 0

		f, err := os.Create("searches.json")
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create output file")
		}
		defer f.Close()

		ctx := context.Background()
		query := `SELECT * FROM public.searches`
		ch, err := db.QueryWithLimit(ctx, query, 10000)

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to query searches")
		}

		for row := range ch {
			var ok bool
			var user *models.User
			search := &models.Search{}

			id, ok := row[0].(int64)
			if !ok {
				logger.Fatal().Msg("Failed to convert row[0] to int64")
			}

			token, ok := row[1].(string)
			if !ok {
				logger.Fatal().Msg("Failed to convert row[1] to string")
			}

			timestamp, ok := row[2].(time.Time)
			if !ok {
				logger.Fatal().Msg("Failed to convert row[2] to time.Time")
			}

			metadata, ok := row[3].(map[string]any)
			if !ok {
				logger.Fatal().Msg("Failed to convert row[3] to map[string]any")
			}

			search.ID        = uint64(id)
			search.Token     = token
			search.Timestamp = datetime.FromTime(timestamp).ToIso8601String()
			search.Metadata  = metadata

			user = meta.Users[token]

			output, _ := sjson.Set("", "id", int(search.ID))
			output, _ = sjson.Set(output, "token", search.Token)
			output, _ = sjson.Set(output, "timestamp", search.Timestamp)

			b, err := search.Metadata.Value()
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to convert metadata to bytes")
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
