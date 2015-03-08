package main

import (
    "testing"
)

func TestParseRewardAmountsNormal(te *testing.T) {
    r, t, err := (&RewardsParser{}).parseRewardAmounts("Limited (5093 left of 20000)")
    if err != nil {
        te.Fail()
    }
    if r != 5093 || t != 20000 {
        te.Fail()
    }
}

func TestParseRewardAmountsZero(te *testing.T) {
    r, _, err := (&RewardsParser{}).parseRewardAmounts("All gone!")
    if err != nil {
        te.Fail()
    }
    if r != 0 {
        te.Fail()
    }
}

func TestParseRewardAmountsFailure(te *testing.T) {
    _, _, err := (&RewardsParser{}).parseRewardAmounts("Aseidjn2871(0")
    if err == nil {
        te.Fail()
    }
}

func TestSendMessage1(t *testing.T) {
    if sendMessage(550, 545) {
        t.Fail()
    }
}

func TestSendMessage2(t *testing.T) {
    if !sendMessage(550, 500) {
        t.Fail()
    }
}

func TestSendMessage3(t *testing.T) {
    if sendMessage(1000, 999) {
        t.Fail()
    }
}

func TestSendMessage4(t *testing.T) {
    if !sendMessage(9, 8) {
        t.Fail()
    }
}

func TestSendMessage5(t *testing.T) {
    if !sendMessage(0, 1) {
        t.Fail()
    }
}

func TestSendMessage6(t *testing.T) {
    if sendMessage(0, 0) {
        t.Fail()
    }
}
