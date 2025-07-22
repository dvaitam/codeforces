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

func dfs(res map[int64]struct{}, lastB int64, remOps int, prod int64, l, r int64) {
	for t := 0; t <= remOps-1; t++ {
		if lastB == 0 && t < 1 {
			continue
		}
		newB := lastB + int64(t)
		if newB <= 0 {
			continue
		}
		nextRem := remOps - t - 1
		if prod > r/newB {
			continue
		}
		newProd := prod * newB
		if newProd > r {
			continue
		}
		if newProd >= l {
			res[newProd] = struct{}{}
		}
		if nextRem > 0 {
			dfs(res, newB, nextRem, newProd, l, r)
		}
	}
}

func expectedE(l, r int64, p int) int {
	res := make(map[int64]struct{})
	dfs(res, 0, p, 1, l, r)
	cnt := 0
	for x := range res {
		if x >= l && x <= r {
			cnt++
		}
	}
	return cnt
}

func genCaseE(rng *rand.Rand) (string, string) {
	l := int64(rng.Intn(20) + 1)
	r := l + int64(rng.Intn(20))
	p := rng.Intn(6) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", l, r, p)
	exp := fmt.Sprintf("%d\n", expectedE(l, r, p))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(expected), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
