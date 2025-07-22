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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func computeExpected(n int, xs, ys []int64) int64 {
	type lineKey struct{ A, B, C int64 }
	lines := make(map[lineKey]map[int]struct{})
	m := len(xs)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			x1, y1 := xs[i], ys[i]
			x2, y2 := xs[j], ys[j]
			A := y2 - y1
			B := x1 - x2
			C := -(A*x1 + B*y1)
			if A == 0 && B == 0 {
				continue
			}
			g := gcd(gcd(abs64(A), abs64(B)), abs64(C))
			A /= g
			B /= g
			C /= g
			if A < 0 || (A == 0 && B < 0) {
				A, B, C = -A, -B, -C
			}
			key := lineKey{A, B, C}
			set, ok := lines[key]
			if !ok {
				set = make(map[int]struct{})
				lines[key] = set
			}
			set[i] = struct{}{}
			set[j] = struct{}{}
		}
	}
	extra := make(map[int]int)
	for k, pts := range lines {
		count := len(pts)
		if count < 2 || k.A == 0 {
			continue
		}
		if (-k.C)%k.A != 0 {
			continue
		}
		x0 := int((-k.C) / k.A)
		if x0 < 1 || x0 > n {
			continue
		}
		if extra[x0] < count-1 {
			extra[x0] = count - 1
		}
	}
	ans := int64(n)
	for _, v := range extra {
		ans += int64(v)
	}
	return ans
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	xs := make([]int64, m)
	ys := make([]int64, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		xs[i] = int64(rng.Intn(n) + 1)
		ys[i] = int64(rng.Intn(20) + 1)
		fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
	}
	ans := computeExpected(n, xs, ys)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
