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

type TestCase struct {
	input  string
	output string
}

// compute the expected outcome for given parameters
func solve(n int, c int, p []int, t []int) string {
	timeL := 0
	scoreL := 0
	for i := 0; i < n; i++ {
		timeL += t[i]
		s := p[i] - c*timeL
		if s < 0 {
			s = 0
		}
		scoreL += s
	}

	timeR := 0
	scoreR := 0
	for i := n - 1; i >= 0; i-- {
		timeR += t[i]
		s := p[i] - c*timeR
		if s < 0 {
			s = 0
		}
		scoreR += s
	}

	if scoreL > scoreR {
		return "Limak"
	} else if scoreL < scoreR {
		return "Radewoosh"
	}
	return "Tie"
}

func generateCase(rng *rand.Rand) TestCase {
	n := rng.Intn(50) + 1
	c := rng.Intn(1000) + 1

	p := make([]int, n)
	t := make([]int, n)

	curP := rng.Intn(20) + 1
	curT := rng.Intn(20) + 1
	for i := 0; i < n; i++ {
		curP += rng.Intn(20) + 1
		curT += rng.Intn(20) + 1
		p[i] = curP
		t[i] = curT
	}

	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, c)
	for i := 0; i < n; i++ {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", p[i])
	}
	in.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", t[i])
	}
	in.WriteByte('\n')

	out := solve(n, c, p, t)
	return TestCase{input: in.String(), output: out}
}

func runCase(bin string, tc TestCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != tc.output {
		return fmt.Errorf("expected %q got %q", tc.output, got)
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

	const cases = 100
	for i := 0; i < cases; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
