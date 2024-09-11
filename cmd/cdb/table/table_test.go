package table

import (
	"testing"
)

func TestCalculateColumnSize(t *testing.T) {
	testCases := []struct {
		table       Table
		expectation int
	}{
		{
			Table{
				Headings: []string{
					"col1",
					"col2",
				},
				Rows: [][]string{
					{
						"a",
						"b",
						"c",
					},
				},
			},
			6,
		},
		{
			Table{
				Headings: []string{
					"col1",
					"col2",
				},
				Rows: [][]string{
					{
						"wichita kansas",
						"b",
						"c",
					},
				},
			},
			16,
		},
	}

	for _, testCase := range testCases {
		maxColSize := testCase.table.calculateColumnSize(0)
		if maxColSize != testCase.expectation {
			t.Errorf("Expected column size %d got: %d", testCase.expectation, maxColSize)
		}
	}
}

func TestPadString(t *testing.T) {
	testCases := []struct {
		padding     int
		data        string
		expectation string
	}{
		{
			5,
			"mat",
			" mat ",
		},
		{
			10,
			"mat",
			" mat      ",
		},
	}

	for _, testCase := range testCases {
		padded := padString(testCase.data, testCase.padding)
		if padded != testCase.expectation {
			t.Errorf("Expected padded string '%s' got: '%s'", testCase.expectation, padded)
		}
	}

}
