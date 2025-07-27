package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	n, k     int
	initial  string
	possible bool
}

// computeMoves replicates the reference algorithm to determine minimal and total swaps.
func computeMoves(n int, s string) ([][]int, int) {
	b := []byte(s)
	moves := make([][]int, 0)
	total := 0
	for {
		cur := 0
		mn := n + 1
		v := make([]int, 0)
		for i := 1; i < n; i++ {
			if b[i] == 'L' {
				cur++
			} else {
				cur--
			}
			if b[i] == 'L' && b[i-1] == 'R' {
				if cur < mn {
					mn = cur
					v = v[:0]
					v = append(v, i)
				} else if cur == mn {
					v = append(v, i)
				}
			}
		}
		if len(v) == 0 {
			break
		}
		for _, x := range v {
			b[x-1], b[x] = b[x], b[x-1]
		}
		total += len(v)
		cp := append([]int(nil), v...)
		moves = append(moves, cp)
	}
	return moves, total
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2 // 2..6
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = 'L'
		} else {
			b[i] = 'R'
		}
	}
	s := string(b)
	moves, tot := computeMoves(n, s)
	kmin := len(moves)
	kmax := tot
	if kmax == 0 {
		kmin = 0
	}
	var k int
	possible := true
	if kmax == 0 {
		if rng.Intn(2) == 0 {
			k = 0
		} else {
			k = rng.Intn(3) + 1
			possible = false
		}
	} else {
		if rng.Intn(2) == 0 {
			k = kmin + rng.Intn(kmax-kmin+1)
			possible = true
		} else {
			if rng.Intn(2) == 0 {
				k = rng.Intn(kmin)
			} else {
				k = kmax + 1 + rng.Intn(3)
			}
			possible = false
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n%s\n", n, k, s))
	return testCase{input: sb.String(), n: n, k: k, initial: s, possible: possible}
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validate(tc testCase, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	tok := scanner.Text()
	if tok == "-1" {
		if tc.possible {
			if scanner.Scan() {
				return fmt.Errorf("unexpected extra output after -1")
			}
			return fmt.Errorf("solution exists but candidate printed -1")
		}
		if scanner.Scan() {
			return fmt.Errorf("unexpected extra output after -1")
		}
		return nil
	}
	if !tc.possible {
		return fmt.Errorf("expected -1 but got answer")
	}
	// tok is first n_i
	b := []byte(tc.initial)
	step := 0
	for {
		if step == tc.k {
			return fmt.Errorf("too many moves")
		}
		n_i, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("bad integer %q", tok)
		}
		if n_i <= 0 || n_i > tc.n/2 {
			return fmt.Errorf("invalid count at step %d", step+1)
		}
		moves := make([]int, n_i)
		for j := 0; j < n_i; j++ {
			if !scanner.Scan() {
				return fmt.Errorf("not enough numbers for step %d", step+1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return fmt.Errorf("bad number at step %d", step+1)
			}
			moves[j] = val
		}
		used := make(map[int]bool)
		for _, x := range moves {
			if x <= 0 || x >= tc.n {
				return fmt.Errorf("bad index %d at step %d", x, step+1)
			}
			if used[x] {
				return fmt.Errorf("duplicate index %d at step %d", x, step+1)
			}
			used[x] = true
			if b[x-1] != 'R' || b[x] != 'L' {
				return fmt.Errorf("invalid pair at step %d index %d", step+1, x)
			}
		}
		for _, x := range moves {
			b[x-1], b[x] = b[x], b[x-1]
		}
		step++
		if step == tc.k {
			break
		}
		if !scanner.Scan() {
			return fmt.Errorf("not enough steps")
		}
		tok = scanner.Text()
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	if strings.Contains(string(b), "RL") {
		return fmt.Errorf("process not finished")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
	// simple deterministic cases
	cases = append(cases, testCase{input: "2 1\nRL\n", n: 2, k: 1, initial: "RL", possible: true})
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if err := validate(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
