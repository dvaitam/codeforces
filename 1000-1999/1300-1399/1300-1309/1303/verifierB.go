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

func solveCase(n, g, b int64) string {
	need := (n + 1) / 2
	periods := (need + g - 1) / g
	full := periods - 1
	rem := need - full*g
	daysHigh := full*(g+b) + rem
	ans := daysHigh
	if n > ans {
		ans = n
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(100) + 1
	g := rng.Int63n(10) + 1
	b := rng.Int63n(10) + 1
	in := fmt.Sprintf("1\n%d %d %d\n", n, g, b)
	out := solveCase(n, g, b)
	return in, out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
