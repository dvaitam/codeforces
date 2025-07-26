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
	n int
	v []int64
	t []int64
}

func solveB(tc testB) []int64 {
	piles := make([]int64, 0)
	res := make([]int64, tc.n)
	for day := 0; day < tc.n; day++ {
		piles = append(piles, tc.v[day])
		melt := int64(0)
		newPiles := make([]int64, 0, len(piles))
		for _, val := range piles {
			if val > tc.t[day] {
				melt += tc.t[day]
				newPiles = append(newPiles, val-tc.t[day])
			} else {
				melt += val
			}
		}
		piles = newPiles
		res[day] = melt
	}
	return res
}

func genTest(rng *rand.Rand) testB {
	n := rng.Intn(8) + 1
	v := make([]int64, n)
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		v[i] = int64(rng.Intn(20))
	}
	for i := 0; i < n; i++ {
		t[i] = int64(rng.Intn(20) + 1)
	}
	return testB{n, v, t}
}

func formatInput(tc testB) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, x := range tc.v {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", x)
	}
	sb.WriteByte('\n')
	for i, x := range tc.t {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", x)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		input := formatInput(tc)
		expVals := solveB(tc)
		var exp strings.Builder
		for j, v := range expVals {
			if j > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprintf(&exp, "%d", v)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp.String() {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp.String(), got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
