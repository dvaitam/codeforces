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

type testCaseA struct {
	n int
	k int64
	f []int64
	t []int64
}

func generateCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(10) + 1
	k := rng.Int63n(20)
	f := make([]int64, n)
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		f[i] = rng.Int63n(50) + 1
		t[i] = rng.Int63n(40) + 1
	}
	return testCaseA{n: n, k: k, f: f, t: t}
}

func buildInputA(tc testCaseA) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.f[i], tc.t[i])
	}
	return sb.String()
}

func expectedA(tc testCaseA) int64 {
	maxJoy := int64(-1 << 60)
	for i := 0; i < tc.n; i++ {
		joy := tc.f[i]
		if tc.t[i] > tc.k {
			joy -= tc.t[i] - tc.k
		}
		if joy > maxJoy {
			maxJoy = joy
		}
	}
	return maxJoy
}

func runCaseA(bin string, tc testCaseA) error {
	input := buildInputA(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	want := expectedA(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInputA(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
