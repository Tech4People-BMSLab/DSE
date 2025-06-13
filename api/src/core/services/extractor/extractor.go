package extractor

import (
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/event"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/xmlpath.v2"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type Packet = models.Packet
type State  = models.State
type User   = models.User
type Task   = models.Task

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
	
	chan_files = make(chan string, 100)
)
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func Slugify(s string) string {
	s   = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s   = re.ReplaceAllString(s, "-")
	s   = strings.Trim(s, "-")
	return s
}

func Unescape(s string) string {
	json.Unmarshal([]byte(s), &s)
	strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "\u0026", "&")
	s = strings.ReplaceAll(s, "\u003c", "<")
	s = strings.ReplaceAll(s, "\u003e", ">")
	s = strings.ReplaceAll(s, "\u0022", "\"")
	s = strings.ReplaceAll(s, "\u0027", "'")
	s = strings.ReplaceAll(s, "\u002f", "/")
	s = strings.ReplaceAll(s, "\u003d", "=")
	s = strings.ReplaceAll(s, "\u0025", "%")
	s = strings.ReplaceAll(s, "\u002b", "+")
	s = strings.ReplaceAll(s, "\u0023", "#")
	s = strings.ReplaceAll(s, "\u003f", "?")
	
	return s
}
// ------------------------------------------------------------
// : Parser
// ------------------------------------------------------------
func ParseGoogle(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Google")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//*[@class='g Ww4FFb vt6azd tF2Cxc asEBEc']`)
		
		//// 1
		path_title := xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  := xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  := xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  := xmlpath.MustCompile(`.//div[@class='VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac lEBKkf']`)


		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)
			
			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

		//// 2
		path_nodes = xmlpath.MustCompile(`//*[@class='g Ww4FFb vt6azd tF2Cxc asEBEc']`)

		path_title = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  = xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  = xmlpath.MustCompile(`.//div[@class='VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac']`)
		

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

		//// 3
		path_nodes = xmlpath.MustCompile(`//div[@class="eKjLze"]`)

		path_title = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  = xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  = xmlpath.MustCompile(`.//div[@class="VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac lEBKkf"]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result3.%d", counter), value)
			counter += 1
		}
	}

	{ //// featured_snippets
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@class="eKjLze"]`)

		path_title  := xmlpath.MustCompile(`.//a`)
		path_link   := xmlpath.MustCompile(`.//a/@href`)
		path_desc   := xmlpath.MustCompile(`.//div[@class="zz3gNc"]`)
		

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	{ //// sidebar_result
		var title       string
		var link        string
		var description string

		path_node := xmlpath.MustCompile(`//div[@class="kp-wholepage ss6qqb u7yw9 zLsiYe mnr-c UBoxCb kp-wholepage-osrp Jb0Zif EyBRub"]`)
		node_iter := path_node.Iter(root)

		if node_iter.Next() {
			path_title  := xmlpath.MustCompile(`.//span[@class="yKMVIe"]`)
			path_link   := xmlpath.MustCompile(`.//a[@class="ruhjFe NJLBac fl"]`)

			title, _ = path_title.String(node_iter.Node())
			title    = Unescape(title)
			if title == "" { return output, nil }

			link, _ = path_link.String(node_iter.Node())
			link    = Unescape(link)
			if link == "" { return output, nil }
		}

		output, _ = sjson.Set(output, "sidebar_result", map[string]string{
			"title"      : title,
			"link"       : link,
			"description": description,
		})

	}

	{ //// people_also_ask
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@class="Wt5Tfe"]//div[@class="dnXCYb"]`)

		path_question := xmlpath.MustCompile(`.//div[@class="L3Ezfd"]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var question    string

			question, _ = path_question.String(it.Node())
			question    = Unescape(question)
			if question == "" { continue }

			value := map[string]string{
				"question": question,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("people_also_ask.%d", counter), value)
			counter += 1
		}
	}

	{ //// related_searches
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//a[@class="k8XOCe R0xfCb VCOFK s8bAkb"]`)

		path_title := xmlpath.MustCompile(`.//div[@class="s75CSd u60jwe r2fjmd AB4Wff"]`)
		path_link  := xmlpath.MustCompile(`.//@href`)


		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			value := map[string]string{
				"title": title,
				"link" : link,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("related_searches.%d", counter), value)
			counter += 1
		}
	}

	{ //// videos
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@jsname="pKB8Bc"]`)

		path_title   := xmlpath.MustCompile(`.//span[@class="cHaqb"]`)
		path_link    := xmlpath.MustCompile(`.//a[@class="X5OiLe"]`)
		path_channel := xmlpath.MustCompile(`.//span[@class='pcJO7e']/span`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title	   string
			var link	   string
			var channel	   string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			channel, _ = path_channel.String(it.Node())
			channel    = Unescape(channel)

			value := map[string]string{
				"title"  : title,
				"link"   : link,
				"channel": channel,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("videos.%d", counter), value)
			counter += 1
		}

		//// 2
		path_nodes = xmlpath.MustCompile(`//div[@jscontroller="rTuANe"]`)

		path_title   = xmlpath.MustCompile(`.//h3[@class="LC20lb MBeuO DKV0Md"]`)
		path_link    = xmlpath.MustCompile(`.//a[@jsname="UWckNb"]`)
		path_channel = xmlpath.MustCompile(`.//div[@class="gqF9jc"]/span[2]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title	   string
			var link	   string
			var channel	   string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			channel, _ = path_channel.String(it.Node())
			channel    = Unescape(channel)

			value := map[string]string{
				"title"  : title,
				"link"   : link,
				"channel": channel,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("videos.%d", counter), value)
			counter += 1
		}
	}
	return output, nil
}

func ParseBing(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Bing")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//li[contains(@class, "b_algo")]`)

		path_title := xmlpath.MustCompile(`.//h2//text()`)
		path_link  := xmlpath.MustCompile(`.//h2//a/@href`)
		path_desc  := xmlpath.MustCompile(`.//p[contains(@class, "b_paractl")]//text()`)

		//// 1
		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string
			var desc  string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			desc_iter := path_desc.Iter(it.Node())
			for desc_iter.Next() {
				desc += desc_iter.Node().String()
			}
			desc    = Unescape(desc)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": desc,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	return output, nil
}

func ParseDuckDuckGo(html string) (string, error) {
	var output string
	
	logger.Info().Msg("Parsing DuckDuckGo")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		path_nodes := xmlpath.MustCompile(`//article`)
		counter := 0
		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			path_title  := xmlpath.MustCompile(`.//div//h2/a/span/text()`)
			path_link   := xmlpath.MustCompile(`.//div/h2/a/@href`)
			path_desc   := xmlpath.MustCompile(`.//article//div[@class='OgdwYG6KE2qthn9XQWFC']//span[@class='kY2IgmnCmOGjharHErah']`)
			
			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

	}

	{ //// recent_news
		path_nodes := xmlpath.MustCompile(`//div[@class="module--carousel__item has-image"]`)
		counter := 0
		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			path_title := xmlpath.MustCompile(`.//a`)
			path_link  := xmlpath.MustCompile(`.//a/@href`)
			path_desc  := xmlpath.MustCompile(`.//span[@class="module--carousel__source result__url"]`)

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)
			if title == "" { continue }

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			if link == "" { continue }

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("recent_news.%d", counter), value)
			counter += 1
		}
	}

	obj := map[string]interface{}{}

	json.Unmarshal([]byte(output), &obj)

	return output, nil
}


// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func OnUpload(user *User, data []byte) {
	parsed := gjson.ParseBytes(data)

	token   := user.Token
	website := Slugify(parsed.Get("website").String())
	keyword := Slugify(parsed.Get("keyword").String())

	path := fmt.Sprintf("./data/extractor/%s.%s.%s.json", token, website, keyword)

	f, err := os.Create(path)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create file")
		return
	}

	f.WriteString(parsed.String())
	chan_files <- path
}

func OnFile() {
	for {
		select {
		case path := <-chan_files:
			event.Emit(event.ExtractorItemStarted)

			var err    error
			var result string

			b, err := os.ReadFile(path)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to read file")
				continue
			}

			logger.Info().Str("path", path).Msg("Reading")

			parsed       := gjson.ParseBytes(b)
			token        := parsed.Get("token").String()
			url          := parsed.Get("url").String()
			browser      := parsed.Get("browser").Value()
			website      := parsed.Get("website").String()
			keyword      := parsed.Get("keyword").String()
			timestamp    := parsed.Get("timestamp").String()
			localization := parsed.Get("localization").String()

			html := parsed.Get("html").String()

			switch website {
				case "Google"    : result, err = ParseGoogle(html)
				case "Bing"      : result, err = ParseBing(html)
				case "DuckDuckGo": result, err = ParseDuckDuckGo(html)

				default: logger.Warn().Str("website", website).Msg("Unknown website"); continue
			}

			if err != nil {
				logger.Error().Err(err).Msg("Failed to parse")
				continue
			}

			var metadata = map[string]interface{}{
				"url"         : url,
				"browser"     : browser,
				"website"     : website,
				"keyword"     : keyword,
				"localization": localization,
				"results"     : gjson.Parse(result).Value(),
			}

			db.CreateSearch(&models.Search{
				Token    : token,
				Timestamp: timestamp,
				Metadata : metadata,
			})

			os.Remove(path)

			event.Emit(event.ExtractorItemDone)
			break
		}
	}
}

// ------------------------------------------------------------
// : Monitor
// ------------------------------------------------------------
func Monitor() {
	for {
		entries, err := os.ReadDir("./data/extractor/")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read directory")
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				chan_files <- fmt.Sprintf("./data/extractor/%s", entry.Name())
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func Init() {
	go Monitor()
	go OnFile()
}
