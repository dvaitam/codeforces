package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func dist(a, b []int) int {
	d := 0
	for i := range a {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	// read n, m
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	pos := -1
	for i := 1; i < n; i++ {
		d := dist(a[0], a[i])
		if d >= 5 {
			fmt.Println("No")
			return
		}
		if d >= 3 {
			pos = i
		}
	}
	// no far row, a[0] is OK
	if pos == -1 {
		fmt.Println("Yes")
		out := make([]string, m)
		for i := 0; i < m; i++ {
			out[i] = strconv.Itoa(a[0][i])
		}
		fmt.Println(strings.Join(out, " "))
		return
	}
	// get differing positions
	breakers := make([]int, 0, m)
	for i := 0; i < m; i++ {
		if a[0][i] != a[pos][i] {
			breakers = append(breakers, i)
		}
	}
	// try when 4 diffs
	if len(breakers) == 4 {
		for mask := 0; mask < (1 << 4); mask++ {
			if bits.OnesCount(uint(mask)) != 2 {
				continue
			}
			x := make([]int, m)
			copy(x, a[0])
			for j := 0; j < 4; j++ {
				if mask&(1<<j) != 0 {
					x[breakers[j]] = a[pos][breakers[j]]
				}
			}
			ok := true
			for i := 0; i < n; i++ {
				if dist(a[i], x) > 2 {
					ok = false
					break
				}
			}
			if ok {
				fmt.Println("Yes")
				out := make([]string, m)
				for i := 0; i < m; i++ {
					out[i] = strconv.Itoa(x[i])
				}
				fmt.Println(strings.Join(out, " "))
				return
			}
		}
		fmt.Println("No")
		return
	}
	// len(breakers) == 3
	// simple masks
	if len(breakers) == 3 {
		for mask := 0; mask < (1 << 3); mask++ {
			pc := bits.OnesCount(uint(mask))
			if pc != 1 && pc != 2 {
				continue
			}
			x := make([]int, m)
			copy(x, a[0])
			for j := 0; j < 3; j++ {
				if mask&(1<<j) != 0 {
					x[breakers[j]] = a[pos][breakers[j]]
				}
			}
			ok := true
			for i := 0; i < n; i++ {
				if dist(a[i], x) > 2 {
					ok = false
					break
				}
			}
			if ok {
				fmt.Println("Yes")
				out := make([]string, m)
				for i := 0; i < m; i++ {
					out[i] = strconv.Itoa(x[i])
				}
				fmt.Println(strings.Join(out, " "))
				return
			}
		}
		// advanced scenario
		for d := 0; d < 3; d++ {
			for e := 0; e < 3; e++ {
				if d == e {
					continue
				}
				x := make([]int, m)
				copy(x, a[0])
				// apply one
				x[breakers[e]] = a[pos][breakers[e]]
				sec := breakers[3-d-e]
				// initialize sec
				x[sec] = 1
				fl := true
				for i := 0; i < n; i++ {
					di := 0
					for j := 0; j < m; j++ {
						if j == sec {
							continue
						}
						if a[i][j] != x[j] {
							di++
						}
					}
					if di > 2 {
						fl = false
						break
					}
					if di == 2 {
						x[sec] = a[i][sec]
					}
				}
				if !fl {
					continue
				}
				for i := 0; i < n; i++ {
					if dist(a[i], x) > 2 {
						fl = false
						break
					}
				}
				if fl {
					fmt.Println("Yes")
					out := make([]string, m)
					for i := 0; i < m; i++ {
						out[i] = strconv.Itoa(x[i])
					}
					fmt.Println(strings.Join(out, " "))
					return
				}
			}
		}
	}
	fmt.Println("No")
}
