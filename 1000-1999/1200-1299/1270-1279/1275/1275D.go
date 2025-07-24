package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	if k < 0 || k > 28 {
		return
	}

	// Precompute helpers
	pairs := [][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 3}, {2, 3}}
	var subsetBits [16]int
	var inc [16][6]int
	for m := 0; m < 16; m++ {
		subsetBits[m] = bits.OnesCount(uint(m))
		for i, p := range pairs {
			if (m>>p[0])&1 == 1 && (m>>p[1])&1 == 1 {
				inc[m][i] = 1
			}
		}
	}

	// decode states
	decode := make([][6]int, 4096)
	for s := 0; s < 4096; s++ {
		x := s
		for i := 5; i >= 0; i-- {
			decode[s][i] = x & 3
			x >>= 2
		}
	}
	// transitions and fail flags
	var trans [4096][16]int
	var failState [4096]bool
	for s := 0; s < 4096; s++ {
		for _, c := range decode[s] {
			if c >= 3 {
				failState[s] = true
				break
			}
		}
		for m := 0; m < 16; m++ {
			var nc [6]int
			for i := 0; i < 6; i++ {
				v := decode[s][i] + inc[m][i]
				if v > 3 {
					v = 3
				}
				nc[i] = v
			}
			idx := 0
			for i := 0; i < 6; i++ {
				idx = idx*4 + nc[i]
			}
			trans[s][m] = idx
		}
	}

	size := 4096 * (k + 1) * 16
	dpPrev := make([]int64, size)
	idx := func(state, t, mask int) int { return ((state*(k+1) + t) << 4) | mask }
	dpPrev[idx(0, 0, 0)] = 1

	for col := 0; col < 7; col++ {
		dpNext := make([]int64, size)
		for state := 0; state < 4096; state++ {
			base := state * (k + 1) * 16
			for t := 0; t <= k; t++ {
				off := base + t*16
				for mask := 0; mask < 16; mask++ {
					val := dpPrev[off+mask]
					if val == 0 {
						continue
					}
					for subset := 0; subset < 16; subset++ {
						nt := t + subsetBits[subset]
						if nt > k {
							continue
						}
						ns := trans[state][subset]
						nm := mask
						if col == 0 {
							nm = subset
						}
						dpNext[((ns*(k+1)+nt)<<4)+nm] += val
					}
				}
			}
		}
		dpPrev = dpNext
	}

	total := int64(0)
	fail := int64(0)
	var countSubset [8]int64
	for state := 0; state < 4096; state++ {
		base := state*(k+1)*16 + k*16
		for mask := 0; mask < 16; mask++ {
			val := dpPrev[base+mask]
			if val == 0 {
				continue
			}
			total += val
			if failState[state] {
				fail += val
			}
			if mask&1 != 0 { // cell (row0,col0) lost
				for sub := 1; sub < 8; sub++ {
					ok := true
					for r := 1; r <= 3 && ok; r++ {
						if (sub>>(r-1))&1 == 1 {
							if mask&(1<<r) == 0 {
								ok = false
								break
							}
							pairIdx := r - 1
							if decode[state][pairIdx] < 3 {
								ok = false
								break
							}
						}
					}
					if ok {
						countSubset[sub] += val
					}
				}
			}
		}
	}

	countUnrec := int64(0)
	for sub := 1; sub < 8; sub++ {
		bits := bits.OnesCount(uint(sub))
		if bits%2 == 1 {
			countUnrec += countSubset[sub]
		} else {
			countUnrec -= countSubset[sub]
		}
	}

	probFail := float64(fail) / float64(total)
	probUnrecCell := float64(countUnrec) / float64(total)
	expected := 15.0 * probUnrecCell
	fmt.Printf("%.10f %.10f\n", probFail, expected)
}
