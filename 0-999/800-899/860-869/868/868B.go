package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var h, m, s, t1, t2 int
	if _, err := fmt.Fscan(reader, &h, &m, &s, &t1, &t2); err != nil {
		return
	}

	// convert positions to [0,12) in units of hours
	hPos := float64(h%12) + float64(m)/60.0 + float64(s)/3600.0
	mPos := float64(m)/5.0 + float64(s)/300.0
	sPos := float64(s) / 5.0
	t1Pos := float64(t1 % 12)
	t2Pos := float64(t2 % 12)

	type item struct {
		pos float64
		id  int
	}
	arr := []item{
		{hPos, 0},
		{mPos, 1},
		{sPos, 2},
		{t1Pos, 3},
		{t2Pos, 4},
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].pos < arr[j].pos })

	// check consecutive intervals for t1 and t2
	for i := 0; i < 5; i++ {
		j := (i + 1) % 5
		if arr[i].id >= 3 && arr[j].id >= 3 {
			fmt.Fprintln(writer, "YES")
			return
		}
	}
	fmt.Fprintln(writer, "NO")
}
