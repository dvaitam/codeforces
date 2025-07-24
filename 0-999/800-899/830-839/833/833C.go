package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var L, R int64
	fmt.Fscan(reader, &L, &R)

	digitsL := make([]int, 19)
	digitsR := make([]int, 19)
	for i := 18; i >= 0; i-- {
		digitsL[i] = int(L % 10)
		digitsR[i] = int(R % 10)
		L /= 10
		R /= 10
	}

	pow20 := [9]uint64{1}
	for i := 1; i < 9; i++ {
		pow20[i] = pow20[i-1] * 20
	}

	// curr[tLow][tHigh][started] -> set of encoded counts
	curr := [2][2][2]map[uint64]struct{}{}
	for i := range curr {
		for j := range curr[i] {
			for k := range curr[i][j] {
				curr[i][j][k] = make(map[uint64]struct{})
			}
		}
	}
	curr[1][1][0][0] = struct{}{}

	for pos := 0; pos < 19; pos++ {
		next := [2][2][2]map[uint64]struct{}{}
		for i := range next {
			for j := range next[i] {
				for k := range next[i][j] {
					next[i][j][k] = make(map[uint64]struct{})
				}
			}
		}
		for tl := 0; tl <= 1; tl++ {
			for tr := 0; tr <= 1; tr++ {
				for st := 0; st <= 1; st++ {
					set := curr[tl][tr][st]
					if len(set) == 0 {
						continue
					}
					low := 0
					high := 9
					if tl == 1 {
						low = digitsL[pos]
					}
					if tr == 1 {
						high = digitsR[pos]
					}
					for code := range set {
						for d := low; d <= high; d++ {
							ntl := 0
							if tl == 1 && d == low {
								ntl = 1
							}
							ntr := 0
							if tr == 1 && d == high {
								ntr = 1
							}
							nst := st
							ncode := code
							if d > 0 {
								nst = 1
								ncode += pow20[d-1]
							}
							nextSet := next[ntl][ntr][nst]
							nextSet[ncode] = struct{}{}
						}
					}
				}
			}
		}
		curr = next
	}

	result := make(map[uint64]struct{})
	for tl := 0; tl <= 1; tl++ {
		for tr := 0; tr <= 1; tr++ {
			for code := range curr[tl][tr][1] {
				result[code] = struct{}{}
			}
		}
	}
	fmt.Println(len(result))
}
