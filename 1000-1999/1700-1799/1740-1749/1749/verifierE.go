package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const INF = 1000000000

type deque struct {
	data       []int
	head, tail int
}

func newDeque(n int) *deque {
	size := 2*n + 10
	d := &deque{data: make([]int, size)}
	d.head = size / 2
	d.tail = d.head
	return d
}

func (d *deque) empty() bool     { return d.head == d.tail }
func (d *deque) pushFront(x int) { d.head--; d.data[d.head] = x }
func (d *deque) pushBack(x int)  { d.data[d.tail] = x; d.tail++ }
func (d *deque) popFront() int   { x := d.data[d.head]; d.head++; return x }

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isok(r, c, n, m int, s [][]byte) bool {
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for _, dxy := range dirs {
		nr := r + dxy[0]
		nc := c + dxy[1]
		if nr >= 0 && nr < n && nc >= 0 && nc < m && s[nr][nc] == '#' {
			return false
		}
	}
	return true
}

func oracle(input string) string {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(r, &n, &m)
		s := make([][]byte, n)
		for i := 0; i < n; i++ {
			var row string
			fmt.Fscan(r, &row)
			s[i] = []byte(row)
		}
		nm := n * m
		dis := make([]int, nm)
		par := make([]int, nm)
		mark := make([]bool, nm)
		for i := 0; i < nm; i++ {
			dis[i] = INF
		}
		dq := newDeque(nm)
		for i := 0; i < n; i++ {
			idx := i * m
			dis[idx] = 0
			if isok(i, 0, n, m, s) {
				if s[i][0] != '#' {
					dis[idx] = 1
					dq.pushBack(idx)
				} else {
					dq.pushFront(idx)
				}
			}
		}
		for !dq.empty() {
			x := dq.popFront()
			if x%m == m-1 || mark[x] {
				if mark[x] {
					continue
				}
				continue
			}
			mark[x] = true
			r0 := x / m
			c0 := x % m
			for _, dxy := range [][2]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}} {
				nr := r0 + dxy[0]
				nc := c0 + dxy[1]
				if nr < 0 || nr >= n || nc < 0 || nc >= m {
					continue
				}
				ni := nr*m + nc
				if mark[ni] {
					continue
				}
				if !isok(nr, nc, n, m, s) {
					continue
				}
				w := 0
				if s[nr][nc] != '#' {
					w = 1
				}
				nd := dis[x] + w
				if nd < dis[ni] {
					dis[ni] = nd
					par[ni] = x
					if w == 0 {
						dq.pushFront(ni)
					} else {
						dq.pushBack(ni)
					}
				}
			}
		}
		ans := INF
		stp := -1
		for i := 0; i < n; i++ {
			idx := i*m + (m - 1)
			if dis[idx] < ans {
				ans = dis[idx]
				stp = idx
			}
		}
		if ans == INF {
			out.WriteString("NO\n")
			continue
		}
		out.WriteString("YES\n")
		x := stp
		for x%m != 0 {
			r1 := x / m
			c1 := x % m
			s[r1][c1] = '#'
			x = par[x]
		}
		s[x/m][0] = '#'
		for i := 0; i < n; i++ {
			out.Write(s[i])
			out.WriteByte('\n')
		}
	}
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if j > 0 && row[j-1] == '#' {
				row[j] = '.'
				continue
			}
			if i > 0 && grid[i-1][j] == '#' {
				row[j] = '.'
				continue
			}
			if rng.Intn(3) == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(46))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expect := oracle(fmt.Sprintf("1\n%s", tc))
		got, err := run(bin, fmt.Sprintf("1\n%s", tc))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
