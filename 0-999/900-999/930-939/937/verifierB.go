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

type testCaseB struct {
	p int
	y int
}

func generateB(rng *rand.Rand) testCaseB {
	p := rng.Intn(40) + 2     // 2..41
	y := p + rng.Intn(60) + 1 // p+1 .. p+60
	return testCaseB{p: p, y: y}
}

func expectedB(p, y int) int {
	for cand := y; cand > p; cand-- {
		limit := int(math.Sqrt(float64(cand)))
		ok := true
		for d := 2; d <= p && d <= limit; d++ {
			if cand%d == 0 {
				ok = false
				break
			}
		}
		if ok {
			return cand
		}
	}
	return -1
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCaseB) error {
	input := fmt.Sprintf("%d %d\n", tc.p, tc.y)
	expected := fmt.Sprintf("%d", expectedB(tc.p, tc.y))
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		tc := generateB(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
