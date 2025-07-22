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

func expectedMoves(x []int, a, b int) int {
	if a == b {
		return 0
	}
	L := a - b
	const INF = int(1e9)
	dist := make([]int, L+1)
	for i := range dist {
		dist[i] = INF
	}
	queue := []int{L}
	dist[L] = 0
	for len(queue) > 0 {
		d := queue[0]
		queue = queue[1:]
		cur := dist[d]
		if d == 0 {
			return cur
		}
		if d > 0 && dist[d-1] > cur+1 {
			dist[d-1] = cur + 1
			queue = append(queue, d-1)
		}
		val := b + d
		for _, xi := range x {
			m := val - val%xi
			if m < b {
				continue
			}
			nd := m - b
			if dist[nd] > cur+1 {
				dist[nd] = cur + 1
				queue = append(queue, nd)
			}
		}
	}
	return dist[0]
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = rng.Intn(9) + 2 // 2..10
	}
	a := rng.Intn(50) + 50 // 50..99
	b := rng.Intn(a + 1)
	if a-b > 30 {
		b = a - 30
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range x {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	fmt.Fprintf(&sb, "\n%d %d\n", a, b)
	return sb.String()
}

func parseCase(input string) ([]int, int, int) {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	x := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &x[i])
	}
	var a, b int
	fmt.Fscan(r, &a, &b)
	return x, a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		x, a, b := parseCase(tc)
		expect := expectedMoves(x, a, b)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", idx+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
