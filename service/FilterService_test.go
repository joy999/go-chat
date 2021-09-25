package service

import (
	"testing"
	"time"
)

func TestFilter(t *testing.T) {

	result := map[string]string{
		"4r5e":        "****",
		"4r5e11":      "****11",
		"a4r5e":       "a****",
		"a4r5e11":     "a****11",
		"ass-fucker":  "**********",
		"ass-fEcker":  "***-fEcker",
		"Eass-fucker": "E**********",
		"Eass-fEcker": "E***-fEcker",
	}

	if msg1 := FilterService.Filter("4r5e"); msg1 != "****" {
		t.Error("4r5e filter failed!")
	}
	if msg1 := FilterService.Filter("4r5e11"); msg1 != "****11" {
		t.Error("4r5e11 filter failed!Now is:", msg1)
	}

	if msg1 := FilterService.Filter("4r5e4r5e"); msg1 != "********" {
		t.Error("4r5e filter failed!")
	}

	for k, v := range result {
		if msg := FilterService.Filter(k); msg != v {
			t.Error(k, " filter failed!Now is:", msg)
		}
	}
	time.Sleep(time.Second * 10)
	FilterService.Filter("abc!hello,world!abc,abc,abc")
	time.Sleep(time.Second * 1)
	FilterService.Filter("haha~~~How are you? hello!")
	time.Sleep(time.Second * 2)
	FilterService.Filter("Hello, everyone! hello?")
	time.Sleep(time.Second * 2)

	s1 := FilterService.PopularWords(1)
	if s1 != "" {
		t.Error("popular test s1 failed!", s1)
	}
	s2 := FilterService.PopularWords(3)
	if s2 != "Hello" {
		t.Error("popular test s2 failed!", s2)
	}
	s3 := FilterService.PopularWords(10)
	if s3 != "abc" {
		t.Error("popular test s3 failed!", s3)
	}
}
