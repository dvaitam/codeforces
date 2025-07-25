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

func runCandidate(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

const INF = int64(1 << 62)

func expected(n int, x0, y0 int64, types []byte, xs, ys []int64) string {
	dist := [8]int64{INF, INF, INF, INF, INF, INF, INF, INF}
	piece := [8]byte{}
	abs := func(v int64) int64 {
		if v < 0 {
			return -v
		}
		return v
	}
	for i := 0; i < n; i++ {
		dx := xs[i] - x0
		dy := ys[i] - y0
		dir := -1
		var d int64
		if dx == 0 {
			if dy > 0 {
				dir, d = 0, dy
			} else if dy < 0 {
				dir, d = 4, -dy
			}
		} else if dy == 0 {
			if dx > 0 {
				dir, d = 2, dx
			} else if dx < 0 {
				dir, d = 6, -dx
			}
		} else if abs(dx) == abs(dy) {
			switch {
			case dx > 0 && dy > 0:
				dir, d = 1, dx
			case dx > 0 && dy < 0:
				dir, d = 3, dx
			case dx < 0 && dy < 0:
				dir, d = 5, -dx
			case dx < 0 && dy > 0:
				dir, d = 7, -dx
			}
		}
		if dir >= 0 && d < dist[dir] {
			dist[dir] = d
			piece[dir] = types[i]
		}
	}
	for dir, t := range piece {
		if t == 0 {
			continue
		}
		if dir%2 == 0 {
			if t == 'R' || t == 'Q' {
				return "YES"
			}
		} else {
			if t == 'B' || t == 'Q' {
				return "YES"
			}
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	x0 := int64(rng.Intn(11) - 5)
	y0 := int64(rng.Intn(11) - 5)
	types := make([]byte, n)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		types[i] = []byte{'B', 'R', 'Q'}[rng.Intn(3)]
		for {
			x := int64(rng.Intn(11) - 5)
			y := int64(rng.Intn(11) - 5)
			if x != x0 || y != y0 {
				xs[i] = x
				ys[i] = y
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(fmt.Sprintf("%d %d\n", x0, y0))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%c %d %d\n", types[i], xs[i], ys[i]))
	}
	input := sb.String()
	exp := expected(n, x0, y0, types, xs, ys)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
