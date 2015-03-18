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

type RewardsParser struct {
    /**
     * Parses out remaining number of rewards
     * Example input: Limited (5093 left of 26000)
     * Example match: (5093
     */
    remainingRegex *regexp.Regexp

    /**
     * Parses out total number of rewards
     * Example input: Limited (5093 left of 26000)
     * Example match:  26000)
     */
    totalRegex *regexp.Regexp

    /**
     * Determines if remaining number of rewards is 0
     */
    soldOutRegex *regexp.Regexp
}

/**
 */
func (p *RewardsParser) compileRegexsIfNecessary() {
    if p.remainingRegex == nil {
        p.remainingRegex = regexp.MustCompile("\\((.*?) ")
        p.totalRegex = regexp.MustCompile("f (.*?)\\)")
        p.soldOutRegex = regexp.MustCompile("All gone!")
    }
}

/**
 * Looks for remaining and total number of rewards
 * @param  string  s   Input string to be searched
 * @return integer r   Remaining rewards
 * @return integer t   Total rewards
 * @return error   err Error if failure
 */
func (p *RewardsParser) parseRewardAmounts(s string) (int64, int64, error) {
    p.compileRegexsIfNecessary()

    if p.soldOutRegex.FindString(s) != "" {
        return 0, 0, nil
    }

    ur := p.remainingRegex.FindString(s)
    if ur == "" {
        return 0, 0, fmt.Errorf("error parsing string for reward amounts: %s", s)
    }
    ut := p.totalRegex.FindString(s)
    if ut == "" {
        return 0, 0, fmt.Errorf("error parsing string for total amount: %s", s)
    }

    var r, t int64
    var err error

    r, err = strconv.ParseInt(strings.Trim(ur, "( "), 10, 32)
    if err != nil {
        return 0, 0, err
    }

    t, err = strconv.ParseInt(strings.Trim(ut, "f )"), 10, 32)
    if err != nil {
        return 0, 0, err
    }

    return r, t, err
}

var rewardsParser *RewardsParser

/*
 * Gets remaining and total amounts for a Kickstarter project reward
 * @param  string proj      Kickstarter project URL suffix like: "597507018/pebble-time-awesome-smartwatch-no-compromises"
 * @param  uint   idx       Zero-based index of reward
 * @return int64  remaining Amount remaining of the reward
 * @return int64  total     Total amount of reward (will be 0 if amount remaining is 0)
 * @return error  err       Error if one has occurred
 */
func scrape(proj string, idx uint) (int64, int64, error) {
    // TODO:
    // Make this method more testable.
    // It's tightly coupled to goquery at the moment.
    doc, err := goquery.NewDocument("https://www.kickstarter.com/projects/" + proj)
    if err != nil {
        return 0, 0, err
    }

    var remaining, total int64
    doc.Find(".NS-projects-reward").EachWithBreak(func(i int, s *goquery.Selection) bool {
        if uint(i) != idx {
            return true
        }

        span := s.Find(".backers-wrap").Text()
        remaining, total, err = rewardsParser.parseRewardAmounts(span)
        return false
    })

    if err != nil {
        return 0, 0, err
    }

    return int64(remaining), int64(total), nil
}

// Given the last sent amount and the current remaining amount, determines if a new text should be sent
func sendMessage(lastSentAmount int64, currentAmount int64) bool {
    if lastSentAmount < 0 {
        // initial
        return true
    }

    lastSentAmount-- // Adjust values since we want to send a new text at the diff multiple
    currentAmount--
    var factor int64 = 10 // Represents the multiple of 10 that is higher than the current amount
    for factor < currentAmount {
        factor *= 10
    }
    diff := factor / 10 // Threshold for sending a new text
    lastTier := (lastSentAmount - lastSentAmount%diff) / diff
    curTier := (currentAmount - currentAmount%diff) / diff
    changed := curTier < lastTier || (currentAmount != lastSentAmount && lastSentAmount < 0)
    return changed
}

var twilioNumber, sendToNumber string
var twilioClient *gotwilio.Twilio
var last int64

func handler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(fmt.Sprintf("<h1>Pebble availability</h1>\n<p>%d Pebble Time Steels of %d are remaining.</p>", last, 20000)))
}

func scrapeAndText() {
    var timeSteelIdx uint = 3
    fmt.Println("Starting scrape")
    remaining, total, err := scrape("597507018/pebble-time-awesome-smartwatch-no-compromises", timeSteelIdx)

    if err != nil {
        fmt.Println("Failed scrape. Error: ", err)
        return
    }

    if remaining < 0 {
        fmt.Println("Failed scrape")
        return
    }

    fmt.Println(remaining, "remaining")

    if !sendMessage(last, remaining) {
        last = remaining
        return
    }

    last = remaining
    message := fmt.Sprintf("%d Pebble Time Steels of %d are remaining.", remaining, total)
    fmt.Printf("Sending message: %s\n", message)
    twilioClient.SendSMS(twilioNumber, sendToNumber, message, "", "")
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
    scrapeAndText()
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("scraped"))
}

func main() {
    twilioSid := os.Getenv("TWILIO_SID")
    twilioToken := os.Getenv("TWILIO_TOKEN")
    twilioClient = gotwilio.NewTwilioClient(twilioSid, twilioToken)
    twilioNumber = os.Getenv("TWILIO_NUMBER")
    sendToNumber = os.Getenv("SEND_TO_NUMBER")

    rewardsParser = &RewardsParser{}

    last = -1

    http.HandleFunc("/", handler)
    http.HandleFunc("/ping", pingHandler)
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
        panic(err)
    }
}
