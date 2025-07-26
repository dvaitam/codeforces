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

type testCase struct{ h1, a1, c1, h2, a2 int }

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func solve(tc testCase) []string {
	h1 := tc.h1
	a1 := tc.a1
	c1 := tc.c1
	h2 := tc.h2
	a2 := tc.a2
	actions := []string{}
	for {
		needStrikes := (h2 + a1 - 1) / a1
		survive := (h1 + a2 - 1) / a2
		if needStrikes <= survive {
			for i := 0; i < needStrikes; i++ {
				actions = append(actions, "STRIKE")
			}
			break
		}
		actions = append(actions, "HEAL")
		h1 += c1
		h1 -= a2
	}
	return actions
}

func validate(tc testCase, actions []string) error {
	h1 := tc.h1
	a1 := tc.a1
	c1 := tc.c1
	h2 := tc.h2
	a2 := tc.a2
	for idx, act := range actions {
		switch act {
		case "HEAL":
			h1 += c1
		case "STRIKE":
			h2 -= a1
		default:
			return fmt.Errorf("invalid action %q", act)
		}
		if h2 <= 0 {
			if idx+1 != len(actions) {
				return fmt.Errorf("extra actions after kill")
			}
			return nil
		}
		h1 -= a2
		if h1 <= 0 {
			return fmt.Errorf("hero died")
		}
	}
	if h2 > 0 {
		return fmt.Errorf("monster still alive")
	}
	return nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d\n%d %d\n", tc.h1, tc.a1, tc.c1, tc.h2, tc.a2)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	reader := strings.NewReader(out)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	actions := make([]string, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &actions[i]); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
	}
	if err := validate(tc, actions); err != nil {
		return err
	}
	exp := len(solve(tc))
	if n != exp {
		return fmt.Errorf("expected %d actions got %d", exp, n)
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
		tc := testCase{
			h1: rng.Intn(100) + 1,
			a1: rng.Intn(100) + 1,
			c1: rng.Intn(99) + 2,
			h2: rng.Intn(100) + 1,
			a2: rng.Intn(99) + 1,
		}
		if tc.a2 >= tc.c1 {
			tc.a2 = tc.c1 - 1
		}
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
