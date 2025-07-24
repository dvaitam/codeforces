package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

func distances(x, y int) [10][10]int {
	var dist [10][10]int
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			dist[i][j] = inf
		}
	}
	for start := 0; start < 10; start++ {
		q := []int{start}
		dist[start][start] = 0
		for head := 0; head < len(q); head++ {
			cur := q[head]
			step := dist[start][cur]
			nxt1 := (cur + x) % 10
			if dist[start][nxt1] > step+1 {
				dist[start][nxt1] = step + 1
				q = append(q, nxt1)
			}
			nxt2 := (cur + y) % 10
			if dist[start][nxt2] > step+1 {
				dist[start][nxt2] = step + 1
				q = append(q, nxt2)
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	digits := []byte(s)
	var cnt [10][10]int64
	for i := 1; i < len(digits); i++ {
		a := digits[i-1] - '0'
		b := digits[i] - '0'
		cnt[a][b]++
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			dist := distances(x, y)
			total := int64(0)
			ok := true
			for a := 0; a < 10 && ok; a++ {
				for b := 0; b < 10; b++ {
					c := cnt[a][b]
					if c == 0 {
						continue
					}
					d := dist[a][b]
					if d == inf {
						ok = false
						break
					}
					total += int64(d-1) * c
				}
			}
			if !ok {
				fmt.Fprint(out, -1)
			} else {
				fmt.Fprint(out, total)
			}
			if y < 9 {
				fmt.Fprint(out, " ")
			} else {
				fmt.Fprintln(out)
			}
		}
	}
}
