package main

// Generate event code from the events.json swagger

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var typeMappings map[string]string

func init() {
	typeMappings = make(map[string]string)
	typeMappings["boolean"] = "bool"
	typeMappings["Channel"] = "ChannelData"
	typeMappings["List[string]"] = "[]string"
	typeMappings["Bridge"] = "BridgeData"
	typeMappings["Playback"] = "PlaybackData"
	typeMappings["LiveRecording"] = "LiveRecordingData"
	typeMappings["StoredRecording"] = "StoredRecordingData"
	typeMappings["Endpoint"] = "EndpointData"
	typeMappings["DeviceState"] = "DeviceStateData"
	typeMappings["TextMessage"] = "TextMessageData"
	typeMappings["object"] = "interface{}"
}

type event struct {
	Name        string
	Event       string
	Description string
	Properties  propList
}

type prop struct {
	Name        string
	JSONName    string
	Mapping     string
	Type        string
	Description string
	Required    bool
}

type propList []prop

func (pl propList) Len() int {
	return len(pl)
}

func (pl propList) Less(l int, r int) bool {
	return pl[l].Name < pl[r].Name
}

func (pl propList) Swap(l int, r int) {
	tmp := pl[r]
	pl[r] = pl[l]
	pl[l] = tmp
}

type eventList []event

func (el eventList) Len() int {
	return len(el)
}

func (el eventList) Less(l int, r int) bool {
	return el[l].Name < el[r].Name
}

func (el eventList) Swap(l int, r int) {
	tmp := el[r]
	el[r] = el[l]
	el[l] = tmp
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <template> <specFile.json>\n", os.Args[0])
		return
	}

	templateFile := os.Args[1]
	specFile := os.Args[2]

	// load template
	tmpl, err := template.New("eventsTemplate").ParseFiles(templateFile)
	if err != nil {
		log.Fatalln("failed to parse template", err)
	}

	// load file
	input, err := os.Open(specFile)
	if err != nil {
		log.Fatalln("failed to open event definition file", err)
	}

	defer input.Close()

	// parse data
	data := make(map[string]interface{})
	dec := json.NewDecoder(input)

	if err := dec.Decode(&data); err != nil {
		log.Fatalln("failed to decode event definition file", err)
	}

	// convert data

	var events eventList

	models, ok := data["models"].(map[string]interface{})
	if !ok {
		log.Fatalln("failed to get models")
	}

	if len(models) < 1 {
		log.Fatalln("no models found")
	}

	for mkey, m := range models {
		model := m.(map[string]interface{})
		name := strings.Replace(mkey, "Id", "ID", -1)

		if name == "Message" || name == "Event" {
			continue
		}

		var pl propList

		props := model["properties"].(map[string]interface{})
		for pkey, p := range props {
			propm := p.(map[string]interface{})
			desc, _ := propm["description"].(string)

			desc = strings.Replace(desc, "\n", "", -1)
			desc = strings.Replace(desc, "\r", "", -1)

			if desc != "" {
				desc = "// " + desc
			}

			t, ok := typeMappings[propm["type"].(string)]
			if !ok {
				t = propm["type"].(string)
			}

			var name string

			items := strings.Split(pkey, "_")
			for _, x := range items {
				name += cases.Title(language.English).String(x)
			}

			required := true
			if req, ok := propm["required"].(bool); ok {
				required = req
			}

			mapping := "`json:\"" + pkey + "\"` "
			if !required {
				mapping = "`json:\"" + pkey + ",omitempty\"`"
			}

			pl = append(pl, prop{
				Name:        name,
				Mapping:     mapping,
				JSONName:    pkey,
				Type:        t,
				Description: desc,
			})
		}

		sort.Sort(pl)

		desc, _ := model["description"].(string)
		desc = strings.Replace(desc, "\n", "", -1)
		desc = strings.Replace(desc, "\r", "", -1)

		events = append(events, event{
			Name:        name,
			Event:       mkey,
			Description: desc,
			Properties:  pl,
		})
	}

	sort.Sort(events)

	if err := tmpl.ExecuteTemplate(os.Stdout, "template.tmpl", events); err != nil {
		log.Fatalln("failed to execute template:", err)
	}
}
