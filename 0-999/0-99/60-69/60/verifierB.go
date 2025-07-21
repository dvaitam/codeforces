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

type node struct{ z, x, y int }

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(r *bufio.Reader) string {
	var k, n, m int
	fmt.Fscan(r, &k, &n, &m)
	grid := make([][][]byte, k)
	for z := 0; z < k; z++ {
		grid[z] = make([][]byte, n)
		for i := 0; i < n; i++ {
			var line string
			fmt.Fscan(r, &line)
			grid[z][i] = []byte(line)
		}
	}
	var x, y int
	fmt.Fscan(r, &x, &y)
	sx, sy := x-1, y-1
	visited := make([][][]bool, k)
	for z := 0; z < k; z++ {
		visited[z] = make([][]bool, n)
		for i := 0; i < n; i++ {
			visited[z][i] = make([]bool, m)
		}
	}
	q := []node{{0, sx, sy}}
	visited[0][sx][sy] = true
	dirs := []node{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}
	count := 0
	for head := 0; head < len(q); head++ {
		u := q[head]
		count++
		for _, d := range dirs {
			nz, nx, ny := u.z+d.z, u.x+d.x, u.y+d.y
			if nz >= 0 && nz < k && nx >= 0 && nx < n && ny >= 0 && ny < m {
				if !visited[nz][nx][ny] && grid[nz][nx][ny] == '.' {
					visited[nz][nx][ny] = true
					q = append(q, node{nz, nx, ny})
				}
			}
		}
	}
	return fmt.Sprintf("%d", count)
}

func generateCaseB(rng *rand.Rand) string {
	k := rng.Intn(3) + 1
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", k, n, m)
	for z := 0; z < k; z++ {
		for i := 0; i < n; i++ {
			line := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Float64() < 0.3 {
					line[j] = '#'
				} else {
					line[j] = '.'
				}
			}
			if z == 0 && i == 0 {
				line[0] = '.'
			}
			sb.WriteString(string(line) + "\n")
		}
	}
	x := 1
	y := 1
	fmt.Fprintf(&sb, "%d %d\n", x, y)
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
