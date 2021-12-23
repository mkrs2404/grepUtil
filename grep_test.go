package main

import (
	"reflect"
	"testing"
)

type grepTest struct {
	testName  string
	searchKey string
	inputs    []string
	want      []string
}

var grepTests = []grepTest{
	{"TestFindMatchesNoSpaceSameCase", "grep", []string{"grep is a great tool", "ghrep does not exist", "This is an attempt to understand grep"}, []string{"grep is a great tool", "This is an attempt to understand grep"}},
	{"TestFindMatchesNoSpaceDiffCase", "grep", []string{"Grep is a great tool", "ghrep does not exist", "This is an attempt to understand gREp"}, []string{"Grep is a great tool", "This is an attempt to understand gREp"}},
	{"TestFindMatchesWithSpace", "gr ep", []string{"Gr ep is a great tool", "grep does not exist", "This is an attempt to understand gre p"}, []string{"Gr ep is a great tool"}}}

func TestFindMatches(t *testing.T) {

	for _, test := range grepTests {
		t.Run(test.testName, func(t *testing.T) {
			got := findMatches(test.searchKey, test.inputs)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FindMatches(%s) got %v, want %v", test.testName, got, test.want)
			}
		})
	}
}
