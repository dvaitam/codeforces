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

type point struct{ x, y int }

func solve(n int, pts []point) string {
	needU, needR, needD, needL := false, false, false, false
	for _, p := range pts {
		if p.y > 0 {
			needU = true
		}
		if p.y < 0 {
			needD = true
		}
		if p.x > 0 {
			needR = true
		}
		if p.x < 0 {
			needL = true
		}
	}
	cnt := 0
	if needU {
		cnt++
	}
	if needR {
		cnt++
	}
	if needD {
		cnt++
	}
	if needL {
		cnt++
	}
	if cnt <= 3 {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	pts := make([]point, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		pts[i] = point{x, y}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	expect := solve(n, pts)
	return sb.String(), expect
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
