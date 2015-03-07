package main

import (
    "fmt"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/sfreiberg/gotwilio"
)

var rgx *regexp.Regexp

func Scrape(proj string, idx uint) int64 {
    doc, err := goquery.NewDocument("https://www.kickstarter.com/projects/" + proj)
    if err != nil {
        fmt.Println("Error getting doc")
        fmt.Println(err)
        return -1
    }

    var remaining int64 = -1
    parseErr := false
    doc.Find(".NS-projects-reward").EachWithBreak(func(i int, s *goquery.Selection) bool {
        if uint(i) != idx {
            return true
        }

        span := s.Find(".limited-number").Text() // Limited (5093 left of 26000)
        raw := rgx.FindString(span)              // (5093  (trailing space)

        if raw == "" {
            // check if none remaining
            span = s.Find(".sold-out").Text()
            if span == "" {
                fmt.Println("Failed to parse number of remaining")
                fmt.Println(s.Text())
                parseErr = true
            } else {
                remaining = 0
            }
            return false
        }

        remaining, err = strconv.ParseInt(strings.Trim(raw, "( "), 10, 32) // 5093
        if err != nil {
            fmt.Println("Error parsing int")
            fmt.Println(err)
            parseErr = true
        }
        return false
    })

    if parseErr {
        return -1
    }

    return remaining
}

// Given the last sent amount and the current remaining amount, determines if a new text should be sent
func sendMessage(lastSentAmount int64, currentAmount int64) bool {
    lastSentAmount-- // Adjust values since we want to send a new text at the diff multiple
    currentAmount--
    var factor int64 = 10 // Represents the multiple of 10 that is higher than the current amount
    for factor < currentAmount {
        factor *= 10
    }
    diff := factor / 10 // Threshold for sending a new text
    lastTier := (lastSentAmount - lastSentAmount%diff) / diff
    curTier := (currentAmount - currentAmount%diff) / diff
    isInitial := lastSentAmount < 0
    changed := curTier < lastTier || (currentAmount != lastSentAmount && lastSentAmount == 0)
    return isInitial || changed
}

var accountSid, authToken, from, to string
var twilio *gotwilio.Twilio
var last int64

func handler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(fmt.Sprintf("<h1>Pebble availability</h1>\n<p>%d Pebble Time Steels of %d are remaining.</p>", last, 20000)))
}

func scrapeAndText() {
    var timeSteelIdx uint = 3
    fmt.Println("Starting scrape")
    remaining := Scrape("597507018/pebble-time-awesome-smartwatch-no-compromises", timeSteelIdx)
    if remaining < 0 {
        fmt.Println("Failed scrape")
        return
    }
    fmt.Println(remaining, "remaining")
    if !sendMessage(last, remaining) {
        return
    }

    last = remaining
    message := fmt.Sprintf("%d Pebble Time Steels of %d are remaining.", remaining, 20000)
    fmt.Printf("Sending message: %s\n", message)
    twilio.SendSMS(from, to, message, "", "")
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
    scrapeAndText()
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("scraped"))
}

func main() {
    accountSid = os.Getenv("twilioSid")
    authToken = os.Getenv("twilioToken")
    twilio = gotwilio.NewTwilioClient(accountSid, authToken)
    from = os.Getenv("fromNum")
    to = os.Getenv("toNum")

    rgx = regexp.MustCompile("\\((.*?) ")
    last = -1

    http.HandleFunc("/", handler)
    http.HandleFunc("/ping", pingHandler)
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
        panic(err)
    }
}
