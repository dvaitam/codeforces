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

type testCase struct{ input string }

type state struct{ i, j, mask int }

func solveCase(in string) string {
	rdr := strings.NewReader(in)
	var n, m int
	fmt.Fscan(rdr, &n, &m)
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &grid[i])
	}
	var si, sj int
	maxDigit := '0'
	bombCount := 0
	for i := 0; i < n; i++ {
		for j, c := range grid[i] {
			if c == 'S' {
				si = i
				sj = j
			}
			if c >= '1' && c <= '8' {
				if c > maxDigit {
					maxDigit = c
				}
			}
			if c == 'B' {
				bombCount++
			}
		}
	}
	t := int(maxDigit - '0')
	treasureVals := make([]int, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(rdr, &treasureVals[i])
	}
	o := t + bombCount
	xs := make([]int, o)
	ys := make([]int, o)
	vals := make([]int, o)
	bombMask := 0
	for i := 0; i < n; i++ {
		for j, c := range grid[i] {
			if c >= '1' && c <= '8' {
				tid := int(c - '1')
				xs[tid] = j
				ys[tid] = i
				vals[tid] = treasureVals[tid]
			}
		}
	}
	bi := t
	for i := 0; i < n; i++ {
		for j, c := range grid[i] {
			if c == 'B' {
				xs[bi] = j
				ys[bi] = i
				vals[bi] = 0
				bombMask |= 1 << bi
				bi++
			}
		}
	}
	maxMask := 1 << o
	const INF = -1
	dist := make([][][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([][]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = make([]int, maxMask)
			for k := 0; k < maxMask; k++ {
				dist[i][j][k] = INF
			}
		}
	}
	q := []state{{si, sj, 0}}
	dist[si][sj][0] = 0
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		d := dist[cur.i][cur.j][cur.mask]
		for _, dir := range dirs {
			ni := cur.i + dir[0]
			nj := cur.j + dir[1]
			if ni < 0 || ni >= n || nj < 0 || nj >= m {
				continue
			}
			c := grid[ni][nj]
			if c == '#' || c == 'B' || (c >= '1' && c <= '8') {
				continue
			}
			newMask := cur.mask
			if cur.j == nj {
				x := cur.j
				y1, y2 := cur.i, ni
				if y1 > y2 {
					y1, y2 = y2, y1
				}
				for k := 0; k < o; k++ {
					if xs[k] < x && ys[k] >= y1 && ys[k] < y2 {
						newMask ^= 1 << k
					}
				}
			}
			if dist[ni][nj][newMask] == INF {
				dist[ni][nj][newMask] = d + 1
				q = append(q, state{ni, nj, newMask})
			}
		}
	}
	ans := 0
	for mask := 0; mask < maxMask; mask++ {
		k := dist[si][sj][mask]
		if k <= 0 {
			continue
		}
		if mask&bombMask != 0 {
			continue
		}
		sum := 0
		for i := 0; i < o; i++ {
			if mask&(1<<i) != 0 {
				sum += vals[i]
			}
		}
		profit := sum - k
		if profit > ans {
			ans = profit
		}
	}
	return fmt.Sprint(ans)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.input)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}
	si := rng.Intn(n)
	sj := rng.Intn(m)
	grid[si][sj] = 'S'
	objCount := rng.Intn(3) + 1
	digit := byte('1')
	bombs := 0
	for i := 0; i < objCount; i++ {
		if rng.Intn(2) == 0 {
			// treasure
			for {
				x := rng.Intn(n)
				y := rng.Intn(m)
				if grid[x][y] == '.' {
					grid[x][y] = digit
					digit++
					break
				}
			}
		} else {
			// bomb
			for {
				x := rng.Intn(n)
				y := rng.Intn(m)
				if grid[x][y] == '.' {
					grid[x][y] = 'B'
					bombs++
					break
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' && rng.Intn(5) == 0 {
				grid[i][j] = '#'
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	treasures := int(digit - '1')
	for i := 0; i < treasures; i++ {
		val := rng.Intn(401) - 200
		sb.WriteString(fmt.Sprintf("%d\n", val))
	}
	return testCase{input: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 105)
	cases = append(cases, randomCase(rng))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
