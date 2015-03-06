package main

import (
    "fmt"
    "log"
    "regexp"
    "strings"
    "strconv"
    "os"
    "time"
    "net/http"

    "github.com/PuerkitoBio/goquery"
    "github.com/sfreiberg/gotwilio"
)

func Scrape(rgx *regexp.Regexp) uint64 {
    doc, err := goquery.NewDocument("https://www.kickstarter.com/projects/597507018/pebble-time-awesome-smartwatch-no-compromises")
    if err != nil {
        log.Fatal(err)
    }

    var remaining uint64
    doc.Find(".limited-number").Each(func(i int, s *goquery.Selection) {
        timeSteel := 3
        if i != timeSteel {
            return
        }
        span := s.Text()
        raw := rgx.FindString(span)
        remaining, _ = strconv.ParseUint(strings.Trim(raw, "( "), 10, 32)
    })

    return remaining
}

func doEvery(d time.Duration, f func()) {
    for {
        time.Sleep(d)
        f()
    }
}

func handler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("<h1>Pebble availability</h1>"))
}

func main() {
    accountSid := os.Getenv("twilioSid")
    authToken := os.Getenv("twilioToken")
    twilio := gotwilio.NewTwilioClient(accountSid, authToken)
    from := os.Getenv("fromNum")
    to := os.Getenv("toNum")

    rgx := regexp.MustCompile("\\((.*?) ")
    var last uint64
    initial := true

    go doEvery(60*time.Second, func() {
        r := Scrape(rgx)
        var factor uint64 = 10
        for ; r > factor; factor *= 10 {
        }
        diff := factor / 10
        lastTier := (last - last % diff) / diff
        curTier := (r - r % diff) % diff
        changed := lastTier != curTier
        if !initial && !changed {
            return
        }
        initial = false
        last = r
        message := fmt.Sprintf("%d Pebble Time Steels of %d are remaining.", r, 20000)
        fmt.Printf("Sending message: %s\n", message)
        twilio.SendSMS(from, to, message, "", "")
    })

    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
    if err != nil {
        panic(err)
    }
}