package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const INF = int(^uint(0) >> 1)

func solve(n, m int, mat []int, specs [][2]int) int {
	K := len(specs)
	nm := n * m
	st := make([]int, nm)
	spec := make([]int, K)
	for i, p := range specs {
		idx := p[0]*m + p[1]
		st[idx] = 1 << i
		spec[i] = idx
	}
	maxMask := 1 << K
	dp := make([][]int, nm)
	inq := make([][]bool, nm)
	preIdx := make([][]int, nm)
	preMask := make([][]int, nm)
	for i := 0; i < nm; i++ {
		dp[i] = make([]int, maxMask)
		inq[i] = make([]bool, maxMask)
		preIdx[i] = make([]int, maxMask)
		preMask[i] = make([]int, maxMask)
		for j := 0; j < maxMask; j++ {
			dp[i][j] = INF
			preIdx[i][j] = -1
		}
	}
	for i := 0; i < K; i++ {
		idx := spec[i]
		mask := 1 << i
		dp[idx][mask] = mat[idx]
	}
	dx := []int{0, 1, 0, -1}
	dy := []int{1, 0, -1, 0}
	type node struct{ u, s int }
	full := maxMask - 1
	last := spec[K-1]
	for s := 1; s <= full; s++ {
		queue := make([]node, 0)
		for u := 0; u < nm; u++ {
			if st[u] != 0 && st[u]&s == 0 {
				continue
			}
			for p := (s - 1) & s; p > 0; p = (p - 1) & s {
				if dp[u][p] >= INF || dp[u][s-p] >= INF {
					continue
				}
				s1 := p | st[u]
				s2 := (s - p) | st[u]
				cost := dp[u][s1] + dp[u][s2] - mat[u]
				if cost < dp[u][s] {
					dp[u][s] = cost
					preIdx[u][s] = u
					preMask[u][s] = s1
				}
			}
			if dp[u][s] < INF {
				inq[u][s] = true
				queue = append(queue, node{u, s})
			}
		}
		for head := 0; head < len(queue); head++ {
			cur := queue[head]
			u, ms := cur.u, cur.s
			inq[u][ms] = false
			ux := u / m
			uy := u % m
			base := dp[u][ms]
			for dir := 0; dir < 4; dir++ {
				vx := ux + dx[dir]
				vy := uy + dy[dir]
				if vx < 0 || vx >= n || vy < 0 || vy >= m {
					continue
				}
				v := vx*m + vy
				ts := ms | st[v]
				nc := base + mat[v]
				if nc < dp[v][ts] {
					dp[v][ts] = nc
					preIdx[v][ts] = u
					preMask[v][ts] = ms
					if !inq[v][ts] && ts == ms {
						inq[v][ts] = true
						queue = append(queue, node{v, ts})
					}
				}
			}
		}
	}
	return dp[last][full]
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(3) + 2 // 2..4
	m := rng.Intn(3) + 2 // 2..4
	k := rng.Intn(3) + 1 // 1..3
	nm := n * m
	mat := make([]int, nm)
	for i := 0; i < nm; i++ {
		mat[i] = rng.Intn(5) + 1
	}
	specs := make([][2]int, k)
	posUsed := make(map[int]bool)
	for i := 0; i < k; i++ {
		for {
			x := rng.Intn(n)
			y := rng.Intn(m)
			idx := x*m + y
			if !posUsed[idx] {
				posUsed[idx] = true
				specs[i] = [2]int{x, y}
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < nm; i++ {
		fmt.Fprintf(&sb, "%d ", mat[i])
	}
	sb.WriteByte('\n')
	for _, p := range specs {
		fmt.Fprintf(&sb, "%d %d\n", p[0]+1, p[1]+1)
	}
	cost := solve(n, m, mat, specs)
	return sb.String(), cost
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(n, m int, specs [][2]int, mat []int, output string, expect int) error {
	r := bufio.NewReader(strings.NewReader(output))
	var cost int
	if _, err := fmt.Fscan(r, &cost); err != nil {
		return fmt.Errorf("cannot read cost")
	}
	if cost != expect {
		return fmt.Errorf("expected cost %d got %d", expect, cost)
	}
	layout := make([]string, n)
	for i := 0; i < n; i++ {
		line, err := r.ReadString('\n')
		if err != nil && i != n-1 {
			return fmt.Errorf("not enough lines")
		}
		line = strings.TrimSpace(line)
		if len(line) != m {
			return fmt.Errorf("invalid line length")
		}
		layout[i] = line
	}
	// check cost from layout
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if layout[i][j] == 'X' {
				sum += mat[i*m+j]
			}
		}
	}
	if sum != cost {
		return fmt.Errorf("cost mismatch")
	}
	// check all special cells covered and connected
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	var start [2]int
	for _, sp := range specs {
		if layout[sp[0]][sp[1]] != 'X' {
			return fmt.Errorf("special cell not covered")
		}
		start = sp
	}
	queue := [][2]int{start}
	visited[start[0]][start[1]] = true
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		for _, d := range dirs {
			nx := p[0] + d[0]
			ny := p[1] + d[1]
			if nx >= 0 && nx < n && ny >= 0 && ny < m {
				if !visited[nx][ny] && layout[nx][ny] == 'X' {
					visited[nx][ny] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}
	}
	for _, sp := range specs {
		if !visited[sp[0]][sp[1]] {
			return fmt.Errorf("special cells not connected")
		}
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
		tc, expect := generateCase(rng)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		// parse input again to obtain matrix and specs
		fields := strings.Fields(tc)
		idx := 0
		n := toInt(fields[idx])
		idx++
		m := toInt(fields[idx])
		idx++
		k := toInt(fields[idx])
		idx++
		nm := n * m
		mat := make([]int, nm)
		for i2 := 0; i2 < nm; i2++ {
			mat[i2] = toInt(fields[idx])
			idx++
		}
		specs := make([][2]int, k)
		for i2 := 0; i2 < k; i2++ {
			x := toInt(fields[idx])
			idx++
			y := toInt(fields[idx])
			idx++
			specs[i2] = [2]int{x - 1, y - 1}
		}
		if err := parseOutput(n, m, specs, mat, got, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func toInt(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}
