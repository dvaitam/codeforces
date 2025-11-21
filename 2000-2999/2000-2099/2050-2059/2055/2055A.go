package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	a, b, turn int
}

const (
	win  = 1
	lose = 2
)

func solve(n int) [][][]int {
	status := make([][][]int, n+1)
	outdeg := make([][][]int, n+1)
	for i := 0; i <= n; i++ {
		status[i] = make([][]int, n+1)
		outdeg[i] = make([][]int, n+1)
		for j := 0; j <= n; j++ {
			status[i][j] = make([]int, 2)
			outdeg[i][j] = make([]int, 2)
		}
	}

	queue := make([]state, 0)
	dirs := []int{-1, 1}

	for a := 1; a <= n; a++ {
		for b := 1; b <= n; b++ {
			if a == b {
				continue
			}
			for turn := 0; turn < 2; turn++ {
				cnt := 0
				if turn == 0 {
					for _, d := range dirs {
						na := a + d
						if na >= 1 && na <= n && na != b {
							cnt++
						}
					}
				} else {
					for _, d := range dirs {
						nb := b + d
						if nb >= 1 && nb <= n && nb != a {
							cnt++
						}
					}
				}
				outdeg[a][b][turn] = cnt
				if cnt == 0 {
					status[a][b][turn] = lose
					queue = append(queue, state{a, b, turn})
				}
			}
		}
	}

	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		curStatus := status[cur.a][cur.b][cur.turn]
		if cur.turn == 0 {
			for _, d := range dirs {
				prevB := cur.b + d
				if prevB < 1 || prevB > n || prevB == cur.a {
					continue
				}
				if status[cur.a][prevB][1] != 0 {
					continue
				}
				if curStatus == lose {
					status[cur.a][prevB][1] = win
					queue = append(queue, state{cur.a, prevB, 1})
				} else {
					outdeg[cur.a][prevB][1]--
					if outdeg[cur.a][prevB][1] == 0 {
						status[cur.a][prevB][1] = lose
						queue = append(queue, state{cur.a, prevB, 1})
					}
				}
			}
		} else {
			for _, d := range dirs {
				prevA := cur.a + d
				if prevA < 1 || prevA > n || prevA == cur.b {
					continue
				}
				if status[prevA][cur.b][0] != 0 {
					continue
				}
				if curStatus == lose {
					status[prevA][cur.b][0] = win
					queue = append(queue, state{prevA, cur.b, 0})
				} else {
					outdeg[prevA][cur.b][0]--
					if outdeg[prevA][cur.b][0] == 0 {
						status[prevA][cur.b][0] = lose
						queue = append(queue, state{prevA, cur.b, 0})
					}
				}
			}
		}
	}

	return status
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(in, &n, &a, &b)
		status := solve(n)
		if status[a][b][0] == win {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
