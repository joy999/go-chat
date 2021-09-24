package service

import "testing"

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

	for k, v := range result {
		if msg := FilterService.Filter(k); msg != v {
			t.Error(k, " filter failed!Now is:", msg)
		}
	}
}
