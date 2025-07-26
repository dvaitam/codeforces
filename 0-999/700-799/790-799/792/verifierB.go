package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testB struct {
	n, k int
	a    []int
}

func solveCaseB(tc testB) []int {
	children := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		children[i] = i + 1
	}
	leader := 0
	res := make([]int, tc.k)
	for i := 0; i < tc.k; i++ {
		idx := (leader + tc.a[i]) % len(children)
		res[i] = children[idx]
		children = append(children[:idx], children[idx+1:]...)
		if len(children) > 0 {
			leader = idx % len(children)
		} else {
			leader = 0
		}
	}
	return res
}

func genCaseB(rng *rand.Rand) (string, []int) {
	n := rng.Intn(98) + 2 // 2..99
	k := rng.Intn(n-1) + 1
	a := make([]int, k)
	for i := range a {
		a[i] = rng.Intn(100) + 1
	}
	tc := testB{n, k, a}
	res := solveCaseB(tc)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), res
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCaseB(bin string, in string, expected []int) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) != len(expected) {
		return fmt.Errorf("expected %d numbers, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("invalid integer %q", tok)
		}
		if val != expected[i] {
			return fmt.Errorf("at position %d expected %d got %d", i+1, expected[i], val)
		}
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
		in, exp := genCaseB(rng)
		if err := runCaseB(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
