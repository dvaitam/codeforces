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

type testCaseB struct {
	n int
	k int
	s string
}

func solveCaseB(tc testCaseB) string {
	cnt := make([]int, 26)
	for i := 0; i < tc.n; i++ {
		cnt[tc.s[i]-'a']++
	}
	pairs := 0
	for _, c := range cnt {
		pairs += c / 2
	}
	remaining := tc.n - tc.k
	if pairs >= remaining/2 {
		return "YES"
	}
	return "NO"
}

func runCaseB(bin string, tc testCaseB) error {
	input := fmt.Sprintf("1\n%d %d\n%s\n", tc.n, tc.k, tc.s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCaseB(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return testCaseB{n: n, k: k, s: string(b)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseB{{n: 4, k: 1, s: "aaaa"}, {n: 3, k: 2, s: "abc"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseB(rng))
	}
	for idx, tc := range cases {
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", idx+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
