package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestFloattostr(t *testing.T) {

	var expected string = "5.50"
	result := floattostr(5.5)
	if result != expected || result == "5.5" {
		t.Errorf("floattostr was incorrect, float was not converted to string " + result)
	}
	result = floattostr(6.823947982)
	if result != "6.82" {
		t.Errorf("floattostr was incorrect, the decimal is greather than two")
	}
}

func TestReplaceStopWords(t *testing.T) {

	csvStopFile, _ := os.Open("stopwords.csv")
	readerStop := csv.NewReader(bufio.NewReader(csvStopFile))
	var stopWords []StopWord
	for {
		line, error := readerStop.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		stopWords = append(stopWords, StopWord{replace_to: line[0], replace_with: line[1]})
		fmt.Println(line)
	}
	/*
		// ".TEST. ,"
		var test string = "test"
		result := replaceStopWords(test, stopWords)
		if result != "test" {
			t.Errorf("replaceStopWords is not working result: " + result)
		}*/
}
