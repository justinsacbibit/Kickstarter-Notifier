package main

import (
    "testing"
)

func TestSendMessage1(t *testing.T) {
    if send := sendMessage(550, 545); send {
        t.Fail()
    }
}

func TestSendMessage2(t *testing.T) {
    if send := sendMessage(550, 500); !send {
        t.Fail()
    }
}

func TestSendMessage3(t *testing.T) {
    if send := sendMessage(1000, 999); send {
        t.Fail()
    }
}

func TestSendMessage4(t *testing.T) {
    if send := sendMessage(9, 8); !send {
        t.Fail()
    }
}

func TestSendMessage5(t *testing.T) {
    if send := sendMessage(0, 1); !send {
        t.Fail()
    }
}

func TestSendMessage6(t *testing.T) {
    if send := sendMessage(0, 0); send {
        t.Fail()
    }
}
