package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// isPossible reports whether a valid placement of platforms exists.
func isPossible(n, m, d int, c []int) bool {
	sum := 0
	for _, v := range c {
		sum += v
	}
	rem := n - sum
	return rem <= (m+1)*(d-1)
}

// validateOutput checks the candidate output for correctness.
func validateOutput(out string, n, m, d int, c []int, ok bool) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	if tokens[0] == "NO" {
		if ok {
			return fmt.Errorf("expected YES but got NO")
		}
		if len(tokens) > 1 {
			return fmt.Errorf("unexpected extra tokens after NO")
		}
		return nil
	}
	if tokens[0] != "YES" {
		return fmt.Errorf("first token must be YES or NO")
	}
	if !ok {
		return fmt.Errorf("expected NO but got YES")
	}
	if len(tokens) != n+1 {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens)-1)
	}
	river := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("invalid integer %q", tokens[i+1])
		}
		river[i] = v
	}
	pos := 0
	for idx := 0; idx < m; idx++ {
		cnt := 0
		for pos < n && river[pos] == 0 {
			cnt++
			pos++
		}
		if cnt >= d {
			return fmt.Errorf("gap before platform %d too large", idx+1)
		}
		for j := 0; j < c[idx]; j++ {
			if pos >= n || river[pos] != idx+1 {
				return fmt.Errorf("platform %d incorrect", idx+1)
			}
			pos++
		}
	}
	cnt := 0
	for pos < n {
		if river[pos] != 0 {
			return fmt.Errorf("unexpected value %d after platforms", river[pos])
		}
		cnt++
		pos++
	}
	if cnt >= d {
		return fmt.Errorf("too many zeros after last platform")
	}
	return nil
}

func generateCase() (string, int, int, int, []int, bool) {
	m := rand.Intn(5) + 1
	c := make([]int, m)
	sum := 0
	for i := 0; i < m; i++ {
		c[i] = rand.Intn(3) + 1
		sum += c[i]
	}
	d := rand.Intn(3) + 1
	var n int
	if rand.Intn(2) == 0 {
		n = sum + rand.Intn((m+1)*(d-1)+1)
	} else {
		n = sum + (m+1)*(d-1) + rand.Intn(5) + 1
	}
	ok := isPossible(n, m, d, c)
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d %d\n", n, m, d)
	for i := 0; i < m; i++ {
		if i+1 == m {
			fmt.Fprintf(&in, "%d\n", c[i])
		} else {
			fmt.Fprintf(&in, "%d ", c[i])
		}
	}
	return in.String(), n, m, d, c, ok
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for i := 0; i < 100; i++ {
		in, n, m, d, c, ok := generateCase()
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validateOutput(strings.TrimSpace(got), n, m, d, append([]int(nil), c...), ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nerror: %v\noutput:\n%s\n", i+1, in, err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
