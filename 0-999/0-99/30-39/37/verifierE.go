package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type caseE struct {
	n, m int
	grid [][]byte
}

func generateCase(rng *rand.Rand) caseE {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = 'B'
			} else {
				row[j] = 'W'
			}
		}
		grid[i] = row
	}
	return caseE{n, m, grid}
}

// solveCase returns the minimal number of painting operations needed
// to obtain the target grid from an all white board. It mimics the
// official solution by running a 0-1 BFS from every starting cell and
// taking the minimum over the maximum distance.
func solveCase(tc caseE) int {
	n, m := tc.n, tc.m
	grid := tc.grid
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	const inf = int(1e9)

	dis := make([][]int, n)
	vis := make([][]bool, n)
	for i := range dis {
		dis[i] = make([]int, m)
		vis[i] = make([]bool, m)
	}

	type pair struct{ x, y int }
	q := make([]pair, 0, n*m)

	ans := inf
	for sx := 0; sx < n; sx++ {
		for sy := 0; sy < m; sy++ {
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					dis[i][j] = inf
					vis[i][j] = false
				}
			}
			q = q[:0]
			dis[sx][sy] = 1
			q = append(q, pair{sx, sy})
			vis[sx][sy] = true

			for h := 0; h < len(q); h++ {
				x, y := q[h].x, q[h].y
				vis[x][y] = false
				for _, d := range dirs {
					nx, ny := x+d[0], y+d[1]
					if nx < 0 || nx >= n || ny < 0 || ny >= m {
						continue
					}
					w := 0
					if grid[nx][ny] != grid[x][y] {
						w = 1
					}
					if dis[nx][ny] > dis[x][y]+w {
						dis[nx][ny] = dis[x][y] + w
						if !vis[nx][ny] {
							q = append(q, pair{nx, ny})
							vis[nx][ny] = true
						}
					}
				}
			}

			f := 0
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					cur := dis[i][j]
					if grid[i][j] == 'W' {
						cur--
					}
					if cur > f {
						f = cur
					}
				}
			}
			if f < ans {
				ans = f
			}
		}
	}
	return ans
}

func runCase(bin string, tc caseE) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			input.WriteByte(tc.grid[i][j])
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprintf("%d", solveCase(tc))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
