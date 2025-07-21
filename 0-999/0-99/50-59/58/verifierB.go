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

func expected(n int) string {
	cur := n
	tmp := n
	res := []int{n}
	for p := 2; p*p <= tmp; p++ {
		for tmp%p == 0 {
			tmp /= p
			cur /= p
			res = append(res, cur)
		}
	}
	if tmp > 1 {
		cur /= tmp
		res = append(res, cur)
	}
	parts := make([]string, len(res))
	for i, v := range res {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ")
}

func runCase(bin string, n int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(n)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := []int{1, 2, 3, 4, 5, 6, 10, 12, 100, 99991}
	for len(nums) < 100 {
		nums = append(nums, rng.Intn(1000000)+1)
	}
	for i, n := range nums {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d\n", i+1, err, n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
