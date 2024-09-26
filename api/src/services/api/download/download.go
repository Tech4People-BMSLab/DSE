package download

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"bms.dse/src/services/db"
	"bms.dse/src/utils/httputil"
	"bms.dse/src/utils/logutil"
	"github.com/golang-module/carbon"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type JSON    = map[string]any
type Writer  = http.ResponseWriter
type Request = http.Request
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("API")
// ------------------------------------------------------------
// : Handler
// ------------------------------------------------------------
func HandleDownloadUsers(w Writer, r *Request) {
	envToken   := os.Getenv("API_TOKEN")
	queryToken := r.URL.Query().Get("token")

	if queryToken != envToken {
		httputil.WriteJSON(w, http.StatusUnauthorized, JSON{"error": "Unauthorized"})
		return
	}

	logger.Info().Msg("Request to Download Users")

	w.Header().Set("Content-Type", "application/json")

	currLimit   := 10_000
	currOffset  := 0
	currSkipped := 0
	currTotal   := 0

	bufferUsers := make([]map[string]any, 0, currLimit)

	for { // Per 10k
		db.Table("users").Offset(currOffset).Limit(currLimit).Find(&bufferUsers)

		if len(bufferUsers) == 0 { break }

		for _, user := range bufferUsers { // Per user
			b, err := json.Marshal(user)
			if err != nil {
				httputil.WriteJSON(w, http.StatusInternalServerError, JSON{"error": "Failed to marshal JSON"})
				continue
			}

			currTotal++

			parsed         := gjson.ParseBytes(b)
			parsedForm     := gjson.Parse(parsed.Get("form").String())
			parsedSocial   := gjson.Parse(parsed.Get("social").String())
			parsedBrowser  := gjson.Parse(parsed.Get("browser").String())
			parsedLanguage := gjson.Parse(parsed.Get("language").String())
			parsedEngine   := gjson.Parse(parsed.Get("search_engine").String())

			if parsedForm.Value() == nil { 
				currSkipped++
				continue 
			}

			user["age"]         = parsedForm.Get("age").String()
			user["sex"]         = parsedForm.Get("sex").String()
			user["income"]      = parsedForm.Get("income").String()
			user["resident"]    = parsedForm.Get("resident").String()
			user["education"]   = parsedForm.Get("education").String()
			user["political"]   = parsedForm.Get("political").String()
			user["employment"]  = parsedForm.Get("employment").String()
			user["postcode"]    = parsedForm.Get("postcode.value").String()

			user["social_tv"]             = parsedSocial.Get("tv").String()
			user["social_radio"]          = parsedSocial.Get("radio").String()
			user["social_anders"]         = parsedSocial.Get("anders").String()
			user["social_reddit"]         = parsedSocial.Get("reddit").String()
			user["social_twitter"]        = parsedSocial.Get("twitter").String()
			user["social_youtube"]        = parsedSocial.Get("youtube").String()
			user["social_de-krant"]       = parsedSocial.Get("de-krant").String()
			user["social_facebook"]       = parsedSocial.Get("facebook").String()
			user["social_linkedin"]       = parsedSocial.Get("linkedin").String()
			user["social_telegram"]       = parsedSocial.Get("telegram").String()
			user["social_whatsapp"]       = parsedSocial.Get("whatsapp").String()
			user["social_instagram"]      = parsedSocial.Get("instagram").String()
			user["social_unselected"]     = parsedSocial.Get("unselected").String()
			user["social_nieuwswebsites"] = parsedSocial.Get("nieuwswebsites").String()

			user["browser_brave"]          = parsedBrowser.Get("brave").Bool()
			user["browser_opera"]          = parsedBrowser.Get("opera").Bool()
			user["browser_chrome"]         = parsedBrowser.Get("chrome").Bool()
			user["browser_safari"]         = parsedBrowser.Get("safari").Bool()
			user["browser_firefox"]        = parsedBrowser.Get("firefox").Bool()
			user["browser_unselected"]     = parsedBrowser.Get("unselected").Bool()
			user["browser_microsoft-edge"] = parsedBrowser.Get("microsoft-edge").Bool()

			user["language_duits"]      = parsedLanguage.Get("duits").Bool()
			user["language_frans"]      = parsedLanguage.Get("frans").Bool()
			user["language_engels"]     = parsedLanguage.Get("engels").Bool()
			user["language_spaans"]     = parsedLanguage.Get("spaans").Bool()
			user["language_italiaans"]  = parsedLanguage.Get("italiaans").Bool()
			user["language_nederlands"] = parsedLanguage.Get("nederlands").Bool()
			user["language_unselected"] = parsedLanguage.Get("unselected").Bool()

			user["search_engine_bing"]       = parsedEngine.Get("bing").Bool()
			user["search_engine_yahoo"]      = parsedEngine.Get("yahoo").Bool()
			user["search_engine_anders"]     = parsedEngine.Get("anders").Bool()
			user["search_engine_ecosia"]     = parsedEngine.Get("ecosia").Bool()
			user["search_engine_google"]     = parsedEngine.Get("google").Bool()
			user["search_engine_startpage"]  = parsedEngine.Get("startpage").Bool()
			user["search_engine_duckduckgo"] = parsedEngine.Get("duckduckgo").Bool()
			user["search_engine_unselected"] = parsedEngine.Get("unselected").Bool()

			

			b, err = json.Marshal(user)
			if err != nil {
				httputil.WriteJSON(w, http.StatusInternalServerError, JSON{"error": "Failed to marshal JSON"})
				continue
			}

			w.Write(b)
			w.Write([]byte("\n"))
		}

		bufferUsers = bufferUsers[:0]
		currOffset += currLimit

		logger.Info().
		Int("offset", currOffset).
		Int("total", currTotal).
		Int("skipped", currSkipped).
		Msg("Download Users Progress")
	}

	logger.Info().
	Int("total", currTotal).
	Int("skipped", currSkipped).
	Msg("Download Users Complete")
}

func HandleDownloadSearches(w Writer, r *Request) {
	// http://localhost:5000/api/download/searches?token=dse2023&days=1
	// http://dev.bmslab.utwente.nl/dse/api/download/searches?token=dse2023&days=1

	envToken   := os.Getenv("API_TOKEN")
	queryToken := r.URL.Query().Get("token")

	if queryToken != envToken {
		httputil.WriteJSON(w, http.StatusUnauthorized, JSON{"error": "Unauthorized"})
		return
	}

	logger.Info().Msg("Request to Download Users")

	w.Header().Set("Content-Type", "application/json")

	queryDays, _ := strconv.Atoi(r.URL.Query().Get("days"))  // Days
	queryStart   := carbon.Parse(r.URL.Query().Get("start")) // Start time
	queryEnd     := carbon.Parse(r.URL.Query().Get("end"))   // End time

	hasDays  := queryDays != 0
	hasStart := !queryStart.IsZero()
	hasEnd   := !queryEnd.IsZero()

	currLimit   := 10_000
	currOffset  := 0
	currTotal   := 0
	currSkipped := 0

	transaction := &gorm.DB{}
	tmpUsers    := make([]JSON, 0, currLimit)
	tmpSearches := make([]JSON, 0, currLimit)
	
	hmUsers := make(map[string]JSON)

	switch (true) {
		case hasDays: {
			timeNow := carbon.Now()
			timeTo  := timeNow.SubDays(queryDays)

			transaction = db.Table("searches").
				Select("searches.*").
				Where("timestamp >= ?", timeTo.ToStdTime()).
				Joins("left join users on searches.user = users.token").
				Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}

		case hasStart && hasEnd: {
			transaction = db.Table("searches").
				Where("timestamp BETWEEN ? AND ?", queryStart.ToStdTime(), queryEnd.ToStdTime()).
				Joins("left join users on searches.user = users.token").
				Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}

		default: {
			transaction = db.Table("searches").
				Joins("left join users on searches.user = users.token").
				Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}
	}


	db.Table("users").Offset(0).Limit(currLimit).Find(&tmpUsers)

	for _, user := range tmpUsers { // Per user
		b, err := json.Marshal(user)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal JSON")
			continue
		}

		parsed := gjson.ParseBytes(b)
		parsedForm      := gjson.Parse(parsed.Get("form").String())
		parsedSocial    := gjson.Parse(parsedForm.Get("social").String())
		parsedBrowser   := gjson.Parse(parsedForm.Get("browser").String())
		parsedLanguage  := gjson.Parse(parsedForm.Get("language").String())
		parsedEngine    := gjson.Parse(parsedForm.Get("search_engine").String())

		user["age"]         = parsedForm.Get("age").String()
		user["sex"]         = parsedForm.Get("sex").String()
		user["income"]      = parsedForm.Get("income").String()
		user["resident"]    = parsedForm.Get("resident").String()
		user["education"]   = parsedForm.Get("education").String()
		user["political"]   = parsedForm.Get("political").String()
		user["employment"]  = parsedForm.Get("employment").String()
		user["postcode"]    = parsedForm.Get("postcode.value").String()

		user["social_tv"]             = parsedSocial.Get("tv").String()
		user["social_radio"]          = parsedSocial.Get("radio").String()
		user["social_anders"]         = parsedSocial.Get("anders").String()
		user["social_reddit"]         = parsedSocial.Get("reddit").String()
		user["social_twitter"]        = parsedSocial.Get("twitter").String()
		user["social_youtube"]        = parsedSocial.Get("youtube").String()
		user["social_de-krant"]       = parsedSocial.Get("de-krant").String()
		user["social_facebook"]       = parsedSocial.Get("facebook").String()
		user["social_linkedin"]       = parsedSocial.Get("linkedin").String()
		user["social_telegram"]       = parsedSocial.Get("telegram").String()
		user["social_whatsapp"]       = parsedSocial.Get("whatsapp").String()
		user["social_instagram"]      = parsedSocial.Get("instagram").String()
		user["social_unselected"]     = parsedSocial.Get("unselected").String()
		user["social_nieuwswebsites"] = parsedSocial.Get("nieuwswebsites").String()

		user["browser_brave"]          = parsedBrowser.Get("brave").Bool()
		user["browser_opera"]          = parsedBrowser.Get("opera").Bool()
		user["browser_chrome"]         = parsedBrowser.Get("chrome").Bool()
		user["browser_safari"]         = parsedBrowser.Get("safari").Bool()
		user["browser_firefox"]        = parsedBrowser.Get("firefox").Bool()
		user["browser_unselected"]     = parsedBrowser.Get("unselected").Bool()
		user["browser_microsoft-edge"] = parsedBrowser.Get("microsoft-edge").Bool()

		user["language_duits"]      = parsedLanguage.Get("duits").Bool()
		user["language_frans"]      = parsedLanguage.Get("frans").Bool()
		user["language_engels"]     = parsedLanguage.Get("engels").Bool()
		user["language_spaans"]     = parsedLanguage.Get("spaans").Bool()
		user["language_italiaans"]  = parsedLanguage.Get("italiaans").Bool()
		user["language_nederlands"] = parsedLanguage.Get("nederlands").Bool()
		user["language_unselected"] = parsedLanguage.Get("unselected").Bool()

		user["search_engine_bing"]       = parsedEngine.Get("bing").Bool()
		user["search_engine_yahoo"]      = parsedEngine.Get("yahoo").Bool()
		user["search_engine_anders"]     = parsedEngine.Get("anders").Bool()
		user["search_engine_ecosia"]     = parsedEngine.Get("ecosia").Bool()
		user["search_engine_google"]     = parsedEngine.Get("google").Bool()
		user["search_engine_startpage"]  = parsedEngine.Get("startpage").Bool()
		user["search_engine_duckduckgo"] = parsedEngine.Get("duckduckgo").Bool()
		user["search_engine_unselected"] = parsedEngine.Get("unselected").Bool()

		hmUsers[user["token"].(string)] = user
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteJSON(w, http.StatusInternalServerError, JSON{"error": "Failed to flush"})
		return
	}

	for { // Per 10k
		transaction.Offset(currOffset).Limit(currLimit).Find(&tmpSearches)
		logger.Debug().Msg(strconv.Itoa(len(tmpSearches)))

		if len(tmpSearches) == 0 { break }

		for _, search := range tmpSearches { // Per search
			search["user"] = hmUsers[search["user"].(string)]
			hasConsented  := search["user"].(JSON)["form"].(string) != "null"

			currTotal++

			if !hasConsented {
				currSkipped++
				continue // Skip user without consent
			} 

			b, err := json.Marshal(search)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal JSON")
				continue
			}


			w.Write(b)
			w.Write([]byte("\n"))

			flusher.Flush()
		}

		tmpSearches = tmpSearches[:0]
		currOffset += currLimit

		logger.Info().
		Int("offset" , currOffset).
		Int("total"  , currTotal).
		Int("skipped", currSkipped).
		Msg("Download Searches Progress")
	}

	logger.Info().
	Int("total", currTotal).
	Int("skipped", currSkipped).
	Msg("Download Searches Complete")
}

func HandleDownloadSearchesBypass(w Writer, r *Request) {
	// http://localhost:5000/api/download/searches/bypass?token=dse2023
	// http://dev.bmslab.utwente.nl/dse/api/download/searches/bypass?token=dse2023&days=1

	envToken   := os.Getenv("API_TOKEN")
	queryToken := r.URL.Query().Get("token")

	if queryToken != envToken {
		httputil.WriteJSON(w, http.StatusUnauthorized, JSON{"error": "Unauthorized"})
		return
	}

	logger.Info().Msg("Request to Download Users")

	w.Header().Set("Content-Type", "application/json")

	queryDays, _ := strconv.Atoi(r.URL.Query().Get("days"))  // Days
	queryStart   := carbon.Parse(r.URL.Query().Get("start")) // Start time
	queryEnd     := carbon.Parse(r.URL.Query().Get("end"))   // End time

	hasDays  := queryDays != 0
	hasStart := !queryStart.IsZero()
	hasEnd   := !queryEnd.IsZero()

	currLimit   := 10_000
	currOffset  := 0
	currTotal   := 0
	currSkipped := 0

	transaction := &gorm.DB{}
	tmpUsers    := make([]JSON, 0, currLimit)
	tmpSearches := make([]JSON, 0, currLimit)
	
	hmUsers := make(map[string]JSON)

	switch (true) {
		case hasDays: {
			timeNow := carbon.Now()
			timeTo  := timeNow.SubDays(queryDays)

			transaction = db.Table("searches").
				Select("searches.*").
				Where("timestamp >= ?", timeTo.ToStdTime()).
				Joins("left join users on searches.user = users.token").
				Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}

		case hasStart && hasEnd: {
			transaction = db.Table("searches").
				Where("timestamp BETWEEN ? AND ?", queryStart.ToStdTime(), queryEnd.ToStdTime()).
				Joins("left join users on searches.user = users.token").
				Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}

		default: {
			transaction = db.Table("searches").
				Joins("left join users on searches.user = users.token").
				// Where("users.form IS NOT NULL").
				Where("searches.results IS NOT NULL").
				Order("timestamp DESC")
			break
		}
	}


	db.Table("users").Offset(0).Limit(currLimit).Find(&tmpUsers)

	for _, user := range tmpUsers { // Per user
		b, err := json.Marshal(user)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal JSON")
			continue
		}

		parsed := gjson.ParseBytes(b)
		parsedForm      := gjson.Parse(parsed.Get("form").String())
		parsedSocial    := gjson.Parse(parsedForm.Get("social").String())
		parsedBrowser   := gjson.Parse(parsedForm.Get("browser").String())
		parsedLanguage  := gjson.Parse(parsedForm.Get("language").String())
		parsedEngine    := gjson.Parse(parsedForm.Get("search_engine").String())

		user["age"]         = parsedForm.Get("age").String()
		user["sex"]         = parsedForm.Get("sex").String()
		user["income"]      = parsedForm.Get("income").String()
		user["resident"]    = parsedForm.Get("resident").String()
		user["education"]   = parsedForm.Get("education").String()
		user["political"]   = parsedForm.Get("political").String()
		user["employment"]  = parsedForm.Get("employment").String()
		user["postcode"]    = parsedForm.Get("postcode.value").String()

		user["social_tv"]             = parsedSocial.Get("tv").String()
		user["social_radio"]          = parsedSocial.Get("radio").String()
		user["social_anders"]         = parsedSocial.Get("anders").String()
		user["social_reddit"]         = parsedSocial.Get("reddit").String()
		user["social_twitter"]        = parsedSocial.Get("twitter").String()
		user["social_youtube"]        = parsedSocial.Get("youtube").String()
		user["social_de-krant"]       = parsedSocial.Get("de-krant").String()
		user["social_facebook"]       = parsedSocial.Get("facebook").String()
		user["social_linkedin"]       = parsedSocial.Get("linkedin").String()
		user["social_telegram"]       = parsedSocial.Get("telegram").String()
		user["social_whatsapp"]       = parsedSocial.Get("whatsapp").String()
		user["social_instagram"]      = parsedSocial.Get("instagram").String()
		user["social_unselected"]     = parsedSocial.Get("unselected").String()
		user["social_nieuwswebsites"] = parsedSocial.Get("nieuwswebsites").String()

		user["browser_brave"]          = parsedBrowser.Get("brave").Bool()
		user["browser_opera"]          = parsedBrowser.Get("opera").Bool()
		user["browser_chrome"]         = parsedBrowser.Get("chrome").Bool()
		user["browser_safari"]         = parsedBrowser.Get("safari").Bool()
		user["browser_firefox"]        = parsedBrowser.Get("firefox").Bool()
		user["browser_unselected"]     = parsedBrowser.Get("unselected").Bool()
		user["browser_microsoft-edge"] = parsedBrowser.Get("microsoft-edge").Bool()

		user["language_duits"]      = parsedLanguage.Get("duits").Bool()
		user["language_frans"]      = parsedLanguage.Get("frans").Bool()
		user["language_engels"]     = parsedLanguage.Get("engels").Bool()
		user["language_spaans"]     = parsedLanguage.Get("spaans").Bool()
		user["language_italiaans"]  = parsedLanguage.Get("italiaans").Bool()
		user["language_nederlands"] = parsedLanguage.Get("nederlands").Bool()
		user["language_unselected"] = parsedLanguage.Get("unselected").Bool()

		user["search_engine_bing"]       = parsedEngine.Get("bing").Bool()
		user["search_engine_yahoo"]      = parsedEngine.Get("yahoo").Bool()
		user["search_engine_anders"]     = parsedEngine.Get("anders").Bool()
		user["search_engine_ecosia"]     = parsedEngine.Get("ecosia").Bool()
		user["search_engine_google"]     = parsedEngine.Get("google").Bool()
		user["search_engine_startpage"]  = parsedEngine.Get("startpage").Bool()
		user["search_engine_duckduckgo"] = parsedEngine.Get("duckduckgo").Bool()
		user["search_engine_unselected"] = parsedEngine.Get("unselected").Bool()

		hmUsers[user["token"].(string)] = user
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteJSON(w, http.StatusInternalServerError, JSON{"error": "Failed to flush"})
		return
	}

	for { // Per 10k
		transaction.Offset(currOffset).Limit(currLimit).Find(&tmpSearches)
		logger.Debug().Msg(strconv.Itoa(len(tmpSearches)))

		if len(tmpSearches) == 0 { break }

		for _, search := range tmpSearches { // Per search
			userToken, ok := search["user"].(string)
			if !ok { continue }

			user, exists := hmUsers[userToken]
			if !exists { continue }

			search["user"] = user

			// search["user"], ok = hmUsers[search["user"].(string)]
			// if !ok { continue }

			// hasConsented  := search["user"].(JSON)["form"].(string) != "null"

			currTotal++

			// if !hasConsented {
			// 	currSkipped++
			// 	continue // Skip user without consent
			// } 

			b, err := json.Marshal(search)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal JSON")
				continue
			}


			w.Write(b)
			w.Write([]byte("\n"))

			flusher.Flush()
		}

		tmpSearches = tmpSearches[:0]
		currOffset += currLimit

		logger.Info().
		Int("offset" , currOffset).
		Int("total"  , currTotal).
		Int("skipped", currSkipped).
		Msg("Download Searches Progress")
	}

	logger.Info().
	Int("total", currTotal).
	Int("skipped", currSkipped).
	Msg("Download Searches Complete")
}
