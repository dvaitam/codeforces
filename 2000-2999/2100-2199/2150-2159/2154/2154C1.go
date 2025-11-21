package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	maxA = 200000 + 2
	inf  = int(1e9)
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner(reader io.Reader) *fastScanner {
	return &fastScanner{r: bufio.NewReader(reader)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

type primeCost struct {
	prime int
	cost  int
}

type bestPair struct {
	first  int
	second int
}

func buildSPF(limit int) []int {
	spf := make([]int, limit+1)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= limit/i {
				for j := i * i; j <= limit; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	spf[1] = 1
	return spf
}

func main() {
	spf := buildSPF(maxA)
	fs := newFastScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	counts := make([]int, maxA+1)
	used := make([]int, 0, 512)

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
		}
		for i := 0; i < n; i++ {
			_ = fs.nextInt()
		}

		used = used[:0]
		found := false
		for i := 0; i < n && !found; i++ {
			temp := a[i]
			for temp > 1 {
				p := spf[temp]
				if counts[p] == 0 {
					used = append(used, p)
				}
				counts[p]++
				if counts[p] >= 2 {
					found = true
					break
				}
				for temp%p == 0 {
					temp /= p
				}
			}
		}

		if found {
			fmt.Fprintln(out, 0)
			for _, p := range used {
				counts[p] = 0
			}
			continue
		}

		min1, min2 := inf, inf
		for _, val := range a {
			c := val & 1
			if c < min1 {
				min2 = min1
				min1 = c
			} else if c < min2 {
				min2 = c
			}
		}
		best := min1 + min2

		bestPairs := make(map[int]bestPair)
		for _, val := range a {
			local := make([]primeCost, 0, 8)
			for delta := 0; delta <= 2; delta++ {
				v := val + delta
				if v <= 1 {
					continue
				}
				temp := v
				for temp > 1 {
					p := spf[temp]
					for temp%p == 0 {
						temp /= p
					}
					idx := -1
					for i := range local {
						if local[i].prime == p {
							if delta < local[i].cost {
								local[i].cost = delta
							}
							idx = i
							break
						}
					}
					if idx == -1 {
						local = append(local, primeCost{prime: p, cost: delta})
					}
				}
			}
			for _, pc := range local {
				bp, ok := bestPairs[pc.prime]
				if !ok {
					bp = bestPair{first: inf, second: inf}
				}
				if pc.cost < bp.first {
					bp.second = bp.first
					bp.first = pc.cost
				} else if pc.cost < bp.second {
					bp.second = pc.cost
				}
				bestPairs[pc.prime] = bp
			}
		}

		for _, bp := range bestPairs {
			if bp.second < inf {
				sum := bp.first + bp.second
				if sum < best {
					best = sum
				}
			}
		}

		fmt.Fprintln(out, best)

		for _, p := range used {
			counts[p] = 0
		}
	}
}
