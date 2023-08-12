package parser_test

import (
	"strings"
	"testing"

	ggx "bitbucket.org/digi-sense/gg-core-x"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_utils"
)

func TestParseContent(t *testing.T) {
	text, err := qbc.IO.ReadTextFromFile("./ultrabike2.csv")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	options := &qb_utils.CsvOptions{
		Comma:          ";",
		Comment:        "",
		FirstRowHeader: true,
	}
	data, err := qbc.CSV.ReadAll(text, options)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	if len(data) > 0 {
		result := make([]map[string]interface{}, 0)
		for _, row := range data {
			lon := row["Longitudine"]
			lat := row["Latitudine"]
			name := row["Name"]
			html := row["description"]

			//fmt.Println(html)
			parser, _ := ggx.HTML.NewParser(html)
			txt := parser.TextAll()
			lines := strings.Split(txt, "\n")
			userTokens := strings.Split(lines[0], " ")
			userId := qbc.Arrays.GetAt(userTokens, 0, "").(string)
			userNum := qbc.Arrays.GetAt(userTokens, 1, "").(string)
			userName := qbc.Arrays.GetAt(userTokens, 3, "").(string)
			statusTokens := strings.Split(lines[1], " ")
			status := strings.Join(statusTokens[1:len(statusTokens)], " ")

			// parse Data
			tag := "<span><b>Data: </b></span>"
			i := strings.Index(html, tag)
			ii := strings.Index(html[i:], "<br /><span><b>")
			date := html[i+len(tag) : i+ii]
			timestamp, _ := qbc.Dates.ParseDate(date, "dd/MM/yyyy HH:mm:ss") // 14/10/2022 06:56:21
			// parse Velocità
			tag = "<span><b>Velocità (Km/h): </b></span>"
			i = strings.Index(html, tag)
			ii = strings.Index(html[i:], "<br /><span><b>")
			velocity := html[i+len(tag) : i+ii]
			// parse Altitudine
			tag = "<span><b>Altitudine (m): </b></span>"
			i = strings.Index(html, tag)
			ii = strings.Index(html[i:], "<br /><span><b>")
			altitude := html[i+len(tag) : i+ii]
			// parse Batteria
			tag = "<span><b>Batteria (%): </b></span>"
			i = strings.Index(html, tag)
			ii = strings.Index(html[i:], "</div>")
			battery := html[i+len(tag) : i+ii]

			// normalize lat, lon
			lat = strings.Replace(lat, ".", ",", 1)
			lat = strings.Replace(lat, ".", "", -1)
			lon = strings.Replace(lon, ".", ",", 1)
			lon = strings.Replace(lon, ".", "", -1)
			lat = strings.Replace(lat, ",", ".", 1)
			lon = strings.Replace(lon, ",", ".", 1)

			if len(lat) < 9 {
				// 404.300
				flat := qbc.Convert.ToFloat64(lat)
				if flat > 100 {
					flat = flat * 0.1
				} else {
					flat = flat * 10
				}
				lat = qbc.Convert.ToString(flat)
			}

			if len(lon) < 8 {
				// 404.300
				flon := qbc.Convert.ToFloat64(lon)
				if flon > 100 {
					flon = flon * 0.01
				} else {
					flon = flon * 100
				}
				lon = qbc.Convert.ToString(flon)
			}

			item := map[string]interface{}{
				"lat":       lat,
				"lon":       lon,
				"name":      name,
				"user_id":   userId,
				"user_nr":   userNum,
				"alias":     userName,
				"status":    status,
				"time":      date,
				"timestamp": qbc.Convert.ToString(timestamp.Unix()),
				"velocity":  velocity,
				"alt":       altitude,
				"battery":   battery,
			}
			result = append(result, item)
		}
		err = qbc.CSV.WriteFile(result, options, "./ultrabike_full.csv")
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
	}
}
