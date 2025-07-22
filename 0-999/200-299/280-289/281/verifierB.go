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

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func solveB(x, y, n int64) string {
	const INF = int64(1 << 62)
	bestNum := INF
	var bestA, bestDen int64 = 0, 1
	for b := int64(1); b <= n; b++ {
		a := x * b / y
		for _, ai := range []int64{a, a + 1} {
			if ai < 0 {
				continue
			}
			diff := abs64(x*b - ai*y)
			if diff*bestDen < bestNum*b ||
				(diff*bestDen == bestNum*b && (b < bestDen || (b == bestDen && ai < bestA))) {
				bestNum = diff
				bestDen = b
				bestA = ai
			}
		}
	}
	return fmt.Sprintf("%d/%d", bestA, bestDen)
}

func generateCase(rng *rand.Rand) (string, string) {
	x := int64(rng.Intn(100000) + 1)
	y := int64(rng.Intn(100000) + 1)
	n := int64(rng.Intn(100000) + 1)
	input := fmt.Sprintf("%d %d %d\n", x, y, n)
	return input, solveB(x, y, n)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
