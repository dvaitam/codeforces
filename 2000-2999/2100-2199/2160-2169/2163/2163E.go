package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type fastReader struct {
	r *bufio.Reader
}

func newFastReader() *fastReader {
	return &fastReader{r: bufio.NewReader(os.Stdin)}
}

func (fr *fastReader) readInt() int {
	sign := 1
	val := 0
	c := fr.readByteSkippingSpaces()
	if c == '-' {
		sign = -1
		c = fr.readByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c = fr.readByte()
	}
	return sign * val
}

func (fr *fastReader) readString() string {
	var buf []byte
	for {
		c := fr.readByte()
		if c == 0 {
			break
		}
		if c > ' ' {
			buf = append(buf, c)
			break
		}
	}
	for {
		c := fr.readByte()
		if c == 0 || c <= ' ' {
			break
		}
		buf = append(buf, c)
	}
	return string(buf)
}

func (fr *fastReader) readByteSkippingSpaces() byte {
	for {
		c := fr.readByte()
		if c == 0 || c > ' ' {
			return c
		}
	}
}

func (fr *fastReader) readByte() byte {
	c, err := fr.r.ReadByte()
	if err != nil {
		return 0
	}
	return c
}

type entry struct {
	s   string
	idx int
}

func findPair(rows, cols []string, wantLE bool) (int, int) {
	n := len(rows)
	rowEntries := make([]entry, n)
	colEntries := make([]entry, n)
	for i := 0; i < n; i++ {
		rowEntries[i] = entry{s: rows[i], idx: i}
		colEntries[i] = entry{s: cols[i], idx: i}
	}
	sort.Slice(rowEntries, func(i, j int) bool {
		return strings.Compare(rowEntries[i].s, rowEntries[j].s) < 0
	})
	sort.Slice(colEntries, func(i, j int) bool {
		return strings.Compare(colEntries[i].s, colEntries[j].s) < 0
	})
	colStrings := make([]string, n)
	for i := 0; i < n; i++ {
		colStrings[i] = colEntries[i].s
	}
	if wantLE {
		for _, row := range rowEntries {
			idx := sort.Search(n, func(k int) bool {
				return strings.Compare(colStrings[k], row.s) >= 0
			})
			if idx < n {
				return row.idx, colEntries[idx].idx
			}
		}
	} else {
		for _, row := range rowEntries {
			idx := sort.Search(n, func(k int) bool {
				return strings.Compare(colStrings[k], row.s) >= 0
			})
			if idx > 0 {
				return row.idx, colEntries[idx-1].idx
			}
		}
	}
	// Fallback to brute force (should not happen theoretically)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cmp := strings.Compare(rows[i], cols[j])
			if wantLE {
				if cmp <= 0 {
					return i, j
				}
			} else {
				if cmp > 0 {
					return i, j
				}
			}
		}
	}
	return 0, 0
}

func buildColumns(grid []string) []string {
	n := len(grid)
	cols := make([]string, n)
	for j := 0; j < n; j++ {
		col := make([]byte, n)
		for i := 0; i < n; i++ {
			col[i] = grid[i][j]
		}
		cols[j] = string(col)
	}
	return cols
}

func runFirst(fr *fastReader, out *bufio.Writer) {
	t := fr.readInt()
	for ; t > 0; t-- {
		n := fr.readInt()
		cVal := fr.readInt()
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			grid[i] = fr.readString()
		}
		cols := buildColumns(grid)
		var r, c int
		if cVal == 1 {
			r, c = findPair(grid, cols, true)
		} else {
			r, c = findPair(grid, cols, false)
		}
		fmt.Fprintf(out, "%d %d\n", r+1, c+1)
	}
}

func runSecond(fr *fastReader, out *bufio.Writer) {
	t := fr.readInt()
	for ; t > 0; t-- {
		_ = fr.readInt() // n
		row := fr.readString()
		col := fr.readString()
		if strings.Compare(row, col) <= 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}

func main() {
	fr := newFastReader()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	mode := fr.readString()
	if mode == "first" {
		runFirst(fr, out)
	} else {
		runSecond(fr, out)
	}
}
