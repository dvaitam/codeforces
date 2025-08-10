package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSubseq(sub []int, b []int) bool {
	j := 0
	for _, x := range b {
		if j < len(sub) && x == sub[j] {
			j++
			if j == len(sub) {
				return true
			}
		}
	}
	return j == len(sub)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var la, lb int
	if _, err := fmt.Fscan(reader, &la, &lb); err != nil {
		return
	}
	a := make([]int, la)
	for i := 0; i < la; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, lb)
	for i := 0; i < lb; i++ {
		fmt.Fscan(reader, &b[i])
	}
	best := 0
	for sa := 0; sa < la; sa++ {
		aa := append(append([]int{}, a[sa:]...), a[:sa]...)
		for sb := 0; sb < lb; sb++ {
			bb := append(append([]int{}, b[sb:]...), b[:sb]...)
			for l := 0; l < la; l++ {
				for r := l; r < la; r++ {
					sub := aa[l : r+1]
					if isSubseq(sub, bb) {
						if len(sub) > best {
							best = len(sub)
						}
					}
				}
			}
		}
	}
	if best > la {
		best = la
	}
	if best > lb {
		best = lb
	}
	fmt.Println(best)
}
