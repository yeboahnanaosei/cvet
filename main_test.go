package main

import (
	"testing"
	"reflect"
	"io"
	"strings"
)

func TestParse(t *testing.T) {
	var table = []struct {
		r       io.Reader
		valid   [][]string
		invalid []invalidRecord
		err     error
	}{
		{
			r: strings.NewReader("name,phone,class\nnana,1111111111,golang\nnana,1111111111,golang\n,4545,accra"),
			valid: [][]string{
				{"nana", "1111111111", "golang"},
				{"nana", "1111111111", "golang"},
			},
			invalid: []invalidRecord{
				{RowNumber: 4, Columns: []string{"name"}},
			},
			err: nil,
		},
		{
			r: strings.NewReader("name,phone,class,address\nnana,1111111111,golang,\nnana,,golang,golang\n,4545,accra,"),
			valid: nil,
			invalid: []invalidRecord{
				{RowNumber: 2, Columns: []string{"address"}},
				{RowNumber: 3, Columns: []string{"phone"}},
				{RowNumber: 4, Columns: []string{"name","address"}},
			},
			err: nil,
		},
		
	}

	for _, tc := range table {
		v, inv, er := parse(tc.r)
		if !reflect.DeepEqual(tc.valid, v) {
			t.Errorf("Expected %v but got %v", tc.valid, v)
		}

		if !reflect.DeepEqual(tc.invalid, inv) {
			t.Errorf("Expected %v but got %v", tc.invalid, inv)
		}

		if er != nil {
			t.Errorf("Expected %v but got %v", tc.err, er)
		}
	}
}
