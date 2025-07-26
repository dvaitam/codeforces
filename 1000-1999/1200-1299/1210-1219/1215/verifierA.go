package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	a1, a2, k1, k2, n int
}

func compute(a1, a2, k1, k2, n int) (int, int) {
	safe := a1*(k1-1) + a2*(k2-1)
	minPlayers := 0
	if n > safe {
		minPlayers = n - safe
	}

	if k1 > k2 {
		k1, k2 = k2, k1
		a1, a2 = a2, a1
	}
	maxPlayers := 0
	use := n / k1
	if use > a1 {
		use = a1
	}
	maxPlayers += use
	n -= use * k1

	use = n / k2
	if use > a2 {
		use = a2
	}
	maxPlayers += use

	return minPlayers, maxPlayers
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(42)
	const T = 100
	for i := 0; i < T; i++ {
		tc := testCase{
			a1: rand.Intn(10) + 1,
			a2: rand.Intn(10) + 1,
			k1: rand.Intn(4) + 1,
			k2: rand.Intn(4) + 1,
		}
		tc.n = rand.Intn(tc.a1*tc.k1+tc.a2*tc.k2) + 1

		input := fmt.Sprintf("%d %d %d %d %d\n", tc.a1, tc.a2, tc.k1, tc.k2, tc.n)
		expMin, expMax := compute(tc.a1, tc.a2, tc.k1, tc.k2, tc.n)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput: %s\n", i+1, err, input)
			os.Exit(1)
		}
		var gotMin, gotMax int
		if _, err := fmt.Sscanf(out, "%d %d", &gotMin, &gotMax); err != nil {
			fmt.Printf("test %d: failed to parse output '%s'\n", i+1, out)
			os.Exit(1)
		}
		if gotMin != expMin || gotMax != expMax {
			fmt.Printf("test %d failed: input %s expected %d %d got %d %d\n", i+1, input, expMin, expMax, gotMin, gotMax)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
