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

const modC = 1000000007

type testCaseC struct {
	n int
	s string
}

func generateCase(rng *rand.Rand) (string, testCaseC) {
	n := rng.Intn(100) + 1
	letters := []byte{'A', 'C', 'G', 'T'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(4)]
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, testCaseC{n: n, s: s}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= modC
	for e > 0 {
		if e&1 == 1 {
			res = res * a % modC
		}
		a = a * a % modC
		e >>= 1
	}
	return res
}

func expected(tc testCaseC) int64 {
	cnt := make(map[byte]int)
	for i := 0; i < tc.n; i++ {
		cnt[tc.s[i]]++
	}
	maxCnt := 0
	for _, c := range []byte{'A', 'C', 'G', 'T'} {
		if cnt[c] > maxCnt {
			maxCnt = cnt[c]
		}
	}
	letters := 0
	for _, c := range []byte{'A', 'C', 'G', 'T'} {
		if cnt[c] == maxCnt {
			letters++
		}
	}
	return modPow(int64(letters), int64(tc.n))
}

func runCase(bin string, input string, tc testCaseC) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
