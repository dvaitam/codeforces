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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func exgcd(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	x = y1
	y = x1 - (a/b)*y1
	return
}

func solveCase(n, m, dx, dy int64, pts [][2]int64) string {
	g, invDx, _ := exgcd(dx, n)
	if g != 1 {
		return ""
	}
	invDx = (invDx%n + n) % n
	yOffset := (n - invDx) % n
	counts := make([]int, n)
	for _, p := range pts {
		t := (p[0] * yOffset) % n * dy % n
		k := (p[1] + t) % n
		counts[k]++
	}
	best := 0
	for i := 1; i < int(n); i++ {
		if counts[i] > counts[best] {
			best = i
		}
	}
	return fmt.Sprintf("0 %d", best)
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(20) + 2)
	var dx, dy int64
	for {
		dx = int64(rng.Intn(int(n-1)) + 1)
		if gcd(n, dx) == 1 {
			break
		}
	}
	for {
		dy = int64(rng.Intn(int(n-1)) + 1)
		if gcd(n, dy) == 1 {
			break
		}
	}
	m := int64(rng.Intn(20) + 1)
	pts := make([][2]int64, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, dx, dy)
	for i := int64(0); i < m; i++ {
		x := int64(rng.Intn(int(n)))
		y := int64(rng.Intn(int(n)))
		pts[i] = [2]int64{x, y}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	expect := solveCase(n, m, dx, dy, pts)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
