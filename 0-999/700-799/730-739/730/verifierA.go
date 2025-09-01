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

type Case struct{ input string }

func runBinary(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 2 // 2..6
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(6))
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func verifyOne(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return fmt.Errorf("bad n: %v", err)
	}
	r := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &r[i]); err != nil {
			return fmt.Errorf("bad r[%d]", i)
		}
	}

	scan := bufio.NewScanner(strings.NewReader(output))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing R")
	}
	var R int
	if _, err := fmt.Sscan(scan.Text(), &R); err != nil {
		return fmt.Errorf("bad R")
	}
	if !scan.Scan() {
		return fmt.Errorf("missing K")
	}
	var K int
	if _, err := fmt.Sscan(scan.Text(), &K); err != nil {
		return fmt.Errorf("bad K")
	}
	if K < 0 {
		return fmt.Errorf("negative K")
	}

	needs := make([]int, n)
	sumNeeds := 0
	for i := 0; i < n; i++ {
		needs[i] = r[i] - R
		if needs[i] < 0 {
			return fmt.Errorf("R too large: r[%d]=%d < R=%d", i, r[i], R)
		}
		sumNeeds += needs[i]
	}
	triples := 0
	pairs := 0
	for step := 0; step < K; step++ {
		if !scan.Scan() {
			return fmt.Errorf("step %d: missing bitstring", step+1)
		}
		line := scan.Text()
		if len(line) != n {
			return fmt.Errorf("step %d: bitstring length %d != n=%d", step+1, len(line), n)
		}
		ones := 0
		for i := 0; i < n; i++ {
			c := line[i]
			if c != '0' && c != '1' {
				return fmt.Errorf("step %d: invalid char", step+1)
			}
			if c == '1' {
				needs[i]--
				if needs[i] < 0 {
					return fmt.Errorf("step %d: index %d decremented too many times", step+1, i+1)
				}
				ones++
			}
		}
		if ones == 3 {
			triples++
		} else if ones == 2 {
			pairs++
		} else {
			return fmt.Errorf("step %d: must select 2 or 3 indices, got %d", step+1, ones)
		}
	}
	for i := 0; i < n; i++ {
		if needs[i] != 0 {
			return fmt.Errorf("index %d remaining %d (not equalized)", i+1, needs[i])
		}
	}
	if 2*pairs+3*triples != sumNeeds {
		return fmt.Errorf("mismatch in total decrements: need %d got %d", sumNeeds, 2*pairs+3*triples)
	}
	if scan.Scan() {
		return fmt.Errorf("extra output: %s", scan.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := verifyOne(c.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
