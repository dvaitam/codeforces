package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	x1 := rng.Int63n(1_000_000)
	y1 := rng.Int63n(1_000_000)
	x2 := rng.Int63n(1_000_000)
	y2 := rng.Int63n(1_000_000)
	x3 := rng.Int63n(1_000_000)
	y3 := rng.Int63n(1_000_000)
	input := fmt.Sprintf("1\n%d %d\n%d %d\n%d %d\n", x1, y1, x2, y2, x3, y3)
	ans := 0.0
	if y1 == y2 && y1 > y3 {
		ans = math.Hypot(float64(x1-x2), float64(y1-y2))
	} else if y2 == y3 && y2 > y1 {
		ans = math.Hypot(float64(x2-x3), float64(y2-y3))
	} else if y1 == y3 && y1 > y2 {
		ans = math.Hypot(float64(x1-x3), float64(y1-y3))
	}
	expect := fmt.Sprintf("%.10f", ans)
	return input, expect
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
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s got: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
