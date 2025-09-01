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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d", val))
	}
	sb.WriteByte('\n')
	return sb.String()
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isPerm(arr []int, l int) bool {
	if l <= 0 {
		return false
	}
	seen := make([]bool, l+1)
	cnt := 0
	mx := 0
	for i := 0; i < l; i++ {
		v := arr[i]
		if v < 1 || v > l || seen[v] {
			return false
		}
		seen[v] = true
		cnt++
		if v > mx {
			mx = v
		}
	}
	return cnt == l && mx == l
}

func validPairs(a []int) map[[2]int]bool {
	n := len(a)
	res := make(map[[2]int]bool)
	for k := 1; k <= n-1; k++ {
		if !isPerm(a[:k], k) {
			continue
		}
		l2 := n - k
		if !isPerm(a[k:], l2) {
			continue
		}
		res[[2]int{k, l2}] = true
	}
	return res
}

func verifyOne(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t, n int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("bad t: %v", err)
	}
	if t != 1 {
		return fmt.Errorf("expected single test in generated case")
	}
	if _, err := fmt.Fscan(in, &n); err != nil {
		return fmt.Errorf("bad n: %v", err)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("missing count")
	}
	var m int
	if _, err := fmt.Sscan(scanner.Text(), &m); err != nil {
		return fmt.Errorf("bad count")
	}
	pairs := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("pair %d: missing first", i+1)
		}
		var x int
		if _, err := fmt.Sscan(scanner.Text(), &x); err != nil {
			return fmt.Errorf("pair %d: bad first", i+1)
		}
		if !scanner.Scan() {
			return fmt.Errorf("pair %d: missing second", i+1)
		}
		var y int
		if _, err := fmt.Sscan(scanner.Text(), &y); err != nil {
			return fmt.Errorf("pair %d: bad second", i+1)
		}
		pairs[[2]int{x, y}] = true
	}
	// no extra tokens allowed
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	vp := validPairs(a)
	// Check all printed pairs are valid
	for p := range pairs {
		if p[0]+p[1] != n || p[0] <= 0 || p[1] <= 0 {
			return fmt.Errorf("invalid pair %d %d", p[0], p[1])
		}
		if !vp[p] {
			return fmt.Errorf("pair %d %d not valid for input", p[0], p[1])
		}
	}
	// Must print exactly all valid pairs (no missing ones)
	if len(pairs) != len(vp) {
		return fmt.Errorf("expected %d pairs, got %d", len(vp), len(pairs))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		out, err := run(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := verifyOne(in, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
