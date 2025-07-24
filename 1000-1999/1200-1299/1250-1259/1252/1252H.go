package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ a, b int64 }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	lands := make([]pair, n)
	var singleScaled int64
	for i := 0; i < n; i++ {
		var l, w int64
		fmt.Fscan(in, &l, &w)
		if l*w > singleScaled {
			singleScaled = l * w
		}
		var a, b int64
		if l < w {
			a = l
			b = w
		} else {
			a = w
			b = l
		}
		lands[i] = pair{a, b}
	}
	sort.Slice(lands, func(i, j int) bool { return lands[i].a > lands[j].a })
	var top1, top2 int64
	var twoScaled int64
	for i, p := range lands {
		b := p.b
		if b > top1 {
			top2 = top1
			top1 = b
		} else if b > top2 {
			top2 = b
		}
		if i >= 1 {
			if top2 >= p.a {
				area := p.a * top2 * 2 // scaled by 2
				if area > twoScaled {
					twoScaled = area
				}
			}
		}
	}
	ansScaled := singleScaled
	if twoScaled > ansScaled {
		ansScaled = twoScaled
	}
	if ansScaled%2 == 0 {
		fmt.Printf("%d.0\n", ansScaled/2)
	} else {
		fmt.Printf("%d.5\n", ansScaled/2)
	}
}
