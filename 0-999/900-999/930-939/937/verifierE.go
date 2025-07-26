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

type testCaseE struct {
	n           int
	s           string
	t           string
	hasSolution bool
}

func isPerm(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	cnt := make([]int, 26)
	for _, c := range a {
		cnt[c-'a']++
	}
	for _, c := range b {
		cnt[c-'a']--
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	return true
}

func generateE(rng *rand.Rand) testCaseE {
	n := rng.Intn(10) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	sb := make([]byte, n)
	tb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(5)]
	}
	if rng.Intn(2) == 0 {
		// permutation
		tb = append([]byte(nil), sb...)
		rng.Shuffle(n, func(i, j int) { tb[i], tb[j] = tb[j], tb[i] })
	} else {
		for i := 0; i < n; i++ {
			tb[i] = letters[rng.Intn(5)]
		}
		for isPerm(string(sb), string(tb)) { // ensure not permutation
			for i := 0; i < n; i++ {
				tb[i] = letters[rng.Intn(5)]
			}
		}
	}
	hs := isPerm(string(sb), string(tb))
	return testCaseE{n: n, s: string(sb), t: string(tb), hasSolution: hs}
}

func applyShift(p string, x int) string {
	n := len(p)
	if x < 0 || x > n {
		return p
	}
	alpha := p[:n-x]
	beta := p[n-x:]
	// reverse beta
	rb := []byte(beta)
	for i, j := 0, len(rb)-1; i < j; i, j = i+1, j-1 {
		rb[i], rb[j] = rb[j], rb[i]
	}
	return string(rb) + alpha
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

func runCase(bin string, tc testCaseE) error {
	input := fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.t)
	gotStr, err := run(bin, input)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(gotStr))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	kStr := scanner.Text()
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k == -1 {
		if tc.hasSolution {
			return fmt.Errorf("should be solvable")
		}
		if scanner.Scan() {
			return fmt.Errorf("extra output after -1")
		}
		return nil
	}
	if !tc.hasSolution {
		return fmt.Errorf("expected -1, got %d", k)
	}
	if k < 0 || k > 6100 {
		return fmt.Errorf("invalid k value")
	}
	ops := make([]int, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing op %d", i)
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid op")
		}
		if v < 0 || v > tc.n {
			return fmt.Errorf("op out of range")
		}
		ops[i] = v
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	cur := tc.s
	for _, x := range ops {
		cur = applyShift(cur, x)
	}
	if cur != tc.t {
		return fmt.Errorf("result mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateE(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
