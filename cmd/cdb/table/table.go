package table

import (
	"fmt"
	"strings"
)

type Table struct {
	Headings []string
	Rows     [][]string
}

func (t Table) String() string {
	paddingSizes := make([]int, len(t.Headings))
	totalWidth := 0
	for idx := range t.Headings {
		paddingSizes[idx] = t.calculateColumnSize(idx)
		totalWidth += paddingSizes[idx]
	}
	totalWidth += 3

	output := strings.Builder{}
	separator := strings.Repeat("-", totalWidth)

	output.WriteString(separator)
	output.WriteRune('\n')

	printCell := func(idx int, data string) {
		if idx == 0 {
			output.WriteString("|")
		}
		output.WriteString(data)
		output.WriteString("|")
	}

	for idx, heading := range t.Headings {
		printCell(idx, padString(heading, paddingSizes[idx]))
	}

	output.WriteRune('\n')
	output.WriteString(separator)
	output.WriteRune('\n')

	for _, row := range t.Rows {
		for idx, data := range row {
			printCell(idx, padString(data, paddingSizes[idx]))
		}

		output.WriteRune('\n')
	}

	output.WriteString(separator)

	return output.String()
}

func (t Table) calculateColumnSize(headerIndex int) int {
	heading := t.Headings[headerIndex]
	maxLength := len(heading)

	for _, data := range t.Rows {
		columnData := data[headerIndex]
		dataLength := len(columnData)
		if dataLength > maxLength {
			maxLength = dataLength
		}
	}

	return maxLength + 2
}

func padString(data string, size int) string {
	// Subtract the prefix padding (i.e. the preceding space we will add)
	paddingSize := size - 1

	// Subtract the length of the data
	paddingSize = paddingSize - len(data)
	padding := strings.Repeat(" ", paddingSize)
	return fmt.Sprintf(" %s%s", data, padding)
}
