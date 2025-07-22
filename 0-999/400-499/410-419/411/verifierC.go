package main

import (
	"bufio"
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
	return out.String(), nil
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var a [5]int
	var b [5]int
	for i := 1; i <= 4; i++ {
		fmt.Fscan(in, &a[i], &b[i])
	}
	var payoff [2][2]int
	for c1 := 0; c1 < 2; c1++ {
		var t1atk, t1def int
		if c1 == 0 {
			t1atk = b[1]
			t1def = a[2]
		} else {
			t1atk = b[2]
			t1def = a[1]
		}
		for c2 := 0; c2 < 2; c2++ {
			var t2atk, t2def int
			if c2 == 0 {
				t2atk = b[3]
				t2def = a[4]
			} else {
				t2atk = b[4]
				t2def = a[3]
			}
			if t1def > t2atk && t1atk > t2def {
				payoff[c1][c2] = 1
			} else if t2def > t1atk && t2atk > t1def {
				payoff[c1][c2] = -1
			} else {
				payoff[c1][c2] = 0
			}
		}
	}
	r0 := payoff[0][0]
	if payoff[0][1] < r0 {
		r0 = payoff[0][1]
	}
	r1 := payoff[1][0]
	if payoff[1][1] < r1 {
		r1 = payoff[1][1]
	}
	result := r0
	if r1 > result {
		result = r1
	}
	switch result {
	case 1:
		return "Team 1\n"
	case -1:
		return "Team 2\n"
	default:
		return "Draw\n"
	}
}

func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < 4; i++ {
		a := rng.Intn(100) + 1
		b := rng.Intn(100) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	return sb.String()
}

func runCase(bin, input string) error {
	expect := strings.TrimSpace(solve(input))
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expect {
		return fmt.Errorf("expected %q got %q", expect, strings.TrimSpace(out))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{}
	for i := 0; i < 100; i++ {
		tests = append(tests, generateCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
