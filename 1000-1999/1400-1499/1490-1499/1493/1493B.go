package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	// valid digit mapping when mirrored
	mirror := map[int]int{0: 0, 1: 1, 2: 5, 5: 2, 8: 8}

	for ; T > 0; T-- {
		var h, m int
		fmt.Fscan(in, &h, &m)
		var s string
		fmt.Fscan(in, &s)

		// parse time HH:MM
		hour := int((s[0]-'0')*10 + (s[1] - '0'))
		minute := int((s[3]-'0')*10 + (s[4] - '0'))

		for i := 0; i < h*m; i++ {
			h1 := hour / 10
			h2 := hour % 10
			m1 := minute / 10
			m2 := minute % 10
			d1, ok1 := mirror[m2]
			d2, ok2 := mirror[m1]
			d3, ok3 := mirror[h2]
			d4, ok4 := mirror[h1]
			if ok1 && ok2 && ok3 && ok4 {
				rh := d1*10 + d2
				rm := d3*10 + d4
				if rh < h && rm < m {
					fmt.Fprintf(out, "%02d:%02d\n", hour, minute)
					break
				}
			}
			// increment time by one minute
			minute++
			if minute == m {
				minute = 0
				hour++
				if hour == h {
					hour = 0
				}
			}
		}
	}
}
