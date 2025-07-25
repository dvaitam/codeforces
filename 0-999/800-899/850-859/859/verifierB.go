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

type CaseB struct {
	input    string
	expected int
}

func expectedB(n int) int {
	best := math.MaxInt32
	limit := int(math.Sqrt(float64(n)))
	for w := 1; w <= limit; w++ {
		h := (n + w - 1) / w
		p := 2 * (w + h)
		if p < best {
			best = p
		}
	}
	return best
}

func generateCaseB(rng *rand.Rand) CaseB {
	n := rng.Intn(1_000_000) + 1
	return CaseB{fmt.Sprintf("%d\n", n), expectedB(n)}
}

func runCase(exe string, input string, expected int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := []CaseB{
		{"1\n", 4},
		{"6\n", 10},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseB(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
