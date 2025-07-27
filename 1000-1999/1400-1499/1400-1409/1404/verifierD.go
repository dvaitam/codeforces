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

func solve(reader *bufio.Reader) string {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	if n%2 == 0 {
		var sb strings.Builder
		sb.WriteString("First\n")
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", i, i+n)
		}
		return strings.TrimSpace(sb.String())
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i], &b[i])
	}
	classA := make([]int, n)
	classB := make([]int, n)
	graph := make([][]struct{ to, idx int }, n)
	for i := 0; i < n; i++ {
		u := (a[i] - 1) % n
		v := (b[i] - 1) % n
		classA[i] = u
		classB[i] = v
		graph[u] = append(graph[u], struct{ to, idx int }{v, i})
		graph[v] = append(graph[v], struct{ to, idx int }{u, i})
	}
	dir := make([]bool, n)
	used := make([]bool, n)
	cycles := make([][]int, 0)
	for start := 0; start < n; start++ {
		if len(graph[start]) == 0 {
			continue
		}
		for _, e0 := range graph[start] {
			if used[e0.idx] {
				continue
			}
			var cycle []int
			curr := start
			prev := -1
			for {
				var e struct{ to, idx int }
				for _, ee := range graph[curr] {
					if ee.idx != prev {
						e = ee
						break
					}
				}
				if used[e.idx] {
					break
				}
				used[e.idx] = true
				cycle = append(cycle, e.idx)
				p := e.idx
				if classA[p] == curr && classB[p] == e.to {
					dir[p] = true
				} else if classB[p] == curr && classA[p] == e.to {
					dir[p] = false
				}
				prev = e.idx
				curr = e.to
				if curr == start {
					break
				}
			}
			if len(cycle) > 0 {
				cycles = append(cycles, cycle)
			}
		}
	}
	t := 0
	for i := 0; i < n; i++ {
		if dir[i] {
			if b[i] > n {
				t++
			}
		} else {
			if a[i] > n {
				t++
			}
		}
	}
	kp := ((n + 1) / 2) & 1
	if (t & 1) != kp {
		for _, cycle := range cycles {
			if len(cycle)%2 == 1 {
				for _, p := range cycle {
					dir[p] = !dir[p]
				}
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("Second\n")
	for i := 0; i < n; i++ {
		if dir[i] {
			fmt.Fprintf(&sb, "%d", b[i])
		} else {
			fmt.Fprintf(&sb, "%d", a[i])
		}
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	even := rng.Intn(2) == 0
	if even {
		n := 2 * (rng.Intn(4) + 1)
		return fmt.Sprintf("%d\n", n)
	}
	n := 2*(rng.Intn(3)+1) + 1
	vals := make([]int, 2*n)
	for i := 1; i <= 2*n; i++ {
		vals[i-1] = i
	}
	rng.Shuffle(2*n, func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", vals[2*i], vals[2*i+1])
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
