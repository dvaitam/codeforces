package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseE struct {
	n     int
	times [][3]int
}

func solveE(tc testCaseE) int {
	n := tc.n
	times := tc.times
	k0 := [3]int{}
	H := make([][]int, 3)
	for i := 0; i < n; i++ {
		for j := 0; j < 3; j++ {
			t := times[i][j]
			if t != 0 {
				k0[j]++
			}
			if t < 0 {
				H[j] = append(H[j], i)
			}
		}
	}
	for j := 0; j < 3; j++ {
		sort.Slice(H[j], func(a, b int) bool {
			ta := -times[H[j][a]][j]
			tb := -times[H[j][b]][j]
			return ta > tb
		})
	}
	scores := []int{500, 1000, 1500, 2000, 2500, 3000}
	bestRank := n + 1
	for _, s0 := range scores {
		for _, s1 := range scores {
			for _, s2 := range scores {
				s := [3]int{s0, s1, s2}
				valid := true
				minx := [3]int{}
				maxx := [3]int{}
				for j := 0; j < 3; j++ {
					L, R := 0, n
					switch s[j] {
					case 500:
						L = n/2 + 1
						R = n
					case 1000:
						L = n/4 + 1
						R = n / 2
					case 1500:
						L = n/8 + 1
						R = n / 4
					case 2000:
						L = n/16 + 1
						R = n / 8
					case 2500:
						L = n/32 + 1
						R = n / 16
					case 3000:
						L = 0
						R = n / 32
					}
					lo := k0[j] - R
					hi := k0[j] - L
					if lo < 0 {
						lo = 0
					}
					if hi > len(H[j]) {
						hi = len(H[j])
					}
					if lo > hi {
						valid = false
						break
					}
					minx[j] = lo
					maxx[j] = hi
				}
				if !valid {
					continue
				}
				for mask := 0; mask < 8; mask++ {
					x := [3]int{}
					for j := 0; j < 3; j++ {
						if (mask>>j)&1 == 0 {
							x[j] = minx[j]
						} else {
							x[j] = maxx[j]
						}
					}
					hacked := make([][3]bool, n)
					for j := 0; j < 3; j++ {
						for t := 0; t < x[j]; t++ {
							i := H[j][t]
							hacked[i][j] = true
						}
					}
					Sc := 0.0
					for j := 0; j < 3; j++ {
						t := times[0][j]
						if t > 0 {
							Sc += float64(s[j]) * float64(250-t) / 250.0
						}
					}
					Sc += float64(100 * (x[0] + x[1] + x[2]))
					rank := 1
					for i := 1; i < n; i++ {
						Pi := 0.0
						for j := 0; j < 3; j++ {
							t := times[i][j]
							if t > 0 && !hacked[i][j] {
								Pi += float64(s[j]) * float64(250-t) / 250.0
							}
						}
						if Pi > Sc {
							rank++
							if rank >= bestRank {
								break
							}
						}
					}
					if rank < bestRank {
						bestRank = rank
					}
				}
			}
		}
	}
	return bestRank
}

func runCaseE(bin string, tc testCaseE) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.times[i][0], tc.times[i][1], tc.times[i][2]))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveE(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(5) + 2
	times := make([][3]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < 3; j++ {
			if i == 0 {
				// you (participant 1) non-negative
				if rng.Intn(2) == 0 {
					times[i][j] = 0
				} else {
					times[i][j] = rng.Intn(120) + 1
				}
			} else {
				t := rng.Intn(241) - 120 // [-120,120]
				if t == 0 {
					times[i][j] = 0
				} else {
					times[i][j] = t
				}
			}
		}
	}
	return testCaseE{n, times}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
