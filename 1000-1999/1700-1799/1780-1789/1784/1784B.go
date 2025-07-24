package main

import (
	"bufio"
	"fmt"
	"os"
)

type op struct {
	a  int
	c1 byte
	b  int
	c2 byte
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	letters := []byte{'w', 'i', 'n'}
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(in, &m)
		var pair [3][3][]int
		for i := 1; i <= m; i++ {
			var s string
			fmt.Fscan(in, &s)
			cnt := [3]int{}
			for j := 0; j < 3; j++ {
				switch s[j] {
				case 'w':
					cnt[0]++
				case 'i':
					cnt[1]++
				case 'n':
					cnt[2]++
				}
			}
			extra := [3]int{}
			need := [3]int{}
			for j := 0; j < 3; j++ {
				if cnt[j] > 1 {
					extra[j] = cnt[j] - 1
				}
				if cnt[j] < 1 {
					need[j] = 1 - cnt[j]
				}
			}
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					for extra[x] > 0 && need[y] > 0 {
						pair[x][y] = append(pair[x][y], i)
						extra[x]--
						need[y]--
					}
				}
			}
		}
		var ops []op
		for x := 0; x < 3; x++ {
			for y := x + 1; y < 3; y++ {
				for len(pair[x][y]) > 0 && len(pair[y][x]) > 0 {
					a := pair[x][y][len(pair[x][y])-1]
					pair[x][y] = pair[x][y][:len(pair[x][y])-1]
					b := pair[y][x][len(pair[y][x])-1]
					pair[y][x] = pair[y][x][:len(pair[y][x])-1]
					ops = append(ops, op{a, letters[x], b, letters[y]})
				}
			}
		}
		cycles := [][]int{{0, 1, 2}, {0, 2, 1}}
		for _, cyc := range cycles {
			x, y, z := cyc[0], cyc[1], cyc[2]
			for len(pair[x][y]) > 0 && len(pair[y][z]) > 0 && len(pair[z][x]) > 0 {
				a := pair[x][y][len(pair[x][y])-1]
				pair[x][y] = pair[x][y][:len(pair[x][y])-1]
				b := pair[y][z][len(pair[y][z])-1]
				pair[y][z] = pair[y][z][:len(pair[y][z])-1]
				c := pair[z][x][len(pair[z][x])-1]
				pair[z][x] = pair[z][x][:len(pair[z][x])-1]
				ops = append(ops, op{a, letters[x], b, letters[y]})
				ops = append(ops, op{b, letters[x], c, letters[z]})
			}
		}
		fmt.Fprintln(out, len(ops))
		for _, v := range ops {
			fmt.Fprintf(out, "%d %c %d %c\n", v.a, v.c1, v.b, v.c2)
		}
	}
}
