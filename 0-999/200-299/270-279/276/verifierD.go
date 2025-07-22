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

type testCaseD struct {
	l uint64
	r uint64
}

func generateCaseD(rng *rand.Rand) testCaseD {
	l := rng.Uint64()%1000 + 1
	r := l + rng.Uint64()%1000
	return testCaseD{l: l, r: r}
}

func expectedD(tc testCaseD) uint64 {
	x := tc.l ^ tc.r
	var ans uint64
	for x > 0 {
		ans = (ans << 1) | 1
		x >>= 1
	}
	return ans
}

func runCaseD(bin string, tc testCaseD) error {
	input := fmt.Sprintf("%d %d\n", tc.l, tc.r)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got uint64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	want := expectedD(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
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
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%d %d\n", i+1, err, tc.l, tc.r)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
