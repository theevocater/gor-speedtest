package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dt/go-metrics-reporting"
)

const (
	RequestFlag          = '1'
	ResponseFlag         = '2'
	ReplayedResponseFlag = '3'
)

type SettingDefs struct {
	src      string
	dst      string
	graphite string
	prefix   string
}

var Settings SettingDefs

func main() {
	flag.StringVar(&Settings.src, "src", "src", "name for traffic src")
	flag.StringVar(&Settings.dst, "dst", "dst", "name for traffic dst")
	flag.StringVar(&Settings.graphite, "graphite", "", "name for traffic dst")
	flag.StringVar(&Settings.prefix, "prefix", "", "prefix for reported timings")
	flag.Parse()

	if Settings.graphite == "" || Settings.prefix == "" {
		log.Fatal("must supply graphite server and prefix")
	}

	report.NewRecorder().
		ReportTo(Settings.graphite, Settings.prefix).
		LogToConsole().
		SetAsDefault()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		encoded := scanner.Bytes()

		buf := make([]byte, hex.DecodedLen(len(encoded)))
		hex.Decode(buf, encoded)

		kind := buf[0]
		if kind == RequestFlag {
			os.Stdout.Write(encoded)
		}
		go handle(kind, buf)
	}
}

func handle(kind byte, buf []byte) {
	headerSize := bytes.IndexByte(buf, '\n') + 1
	header := buf[:headerSize-1]
	meta := bytes.Split(header, []byte(" "))
	log.Println("meta", meta)

	//id := string(meta[1])
	ts, _ := strconv.ParseInt(string(meta[2]), 10, 64)

	if kind == ResponseFlag {
		report.Time(Settings.src, time.Duration(ts))
	} else if kind == ReplayedResponseFlag {
		report.Time(Settings.src, time.Duration(ts))
	}
}
