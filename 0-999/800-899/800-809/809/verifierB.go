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

type testB struct {
	n      int
	dishes []int
}

func generateCaseB(rng *rand.Rand) testB {
	n := rng.Intn(50) + 2
	k := rng.Intn(n-1) + 2 // at least 2
	perm := rng.Perm(n)
	dishes := make([]int, k)
	for i := 0; i < k; i++ {
		dishes[i] = perm[i] + 1
	}
	return testB{n, dishes}
}

func runCase(bin string, tc testB) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.dishes)))
	for i, v := range tc.dishes {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var t int
	var x, y int
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &t, &x, &y); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if t != 2 || x < 1 || x > tc.n || y < 1 || y > tc.n || x == y {
		return fmt.Errorf("bad output values")
	}
	okx, oky := false, false
	for _, d := range tc.dishes {
		if d == x {
			okx = true
		}
		if d == y {
			oky = true
		}
	}
	if !okx || !oky {
		return fmt.Errorf("indices not from ordered set")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d k=%d dishes=%v\n", i+1, err, tc.n, len(tc.dishes), tc.dishes)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
