package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expectedA(s, x int64) int64 {
	if s < x || (s-x)%2 != 0 {
		return 0
	}
	t := (s - x) / 2
	if t&x != 0 {
		return 0
	}
	cnt := 0
	for tmp := x; tmp > 0; tmp >>= 1 {
		if tmp&1 == 1 {
			cnt++
		}
	}
	res := int64(1) << cnt
	if t == 0 {
		res -= 2
	}
	return res
}

func runCase(exe string, s, x int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d %d\n", s, x))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprint(expectedA(s, x))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		s := rng.Int63n(1_000_000_000) + 2
		x := rng.Int63n(1_000_000_000)
		if err := runCase(exe, s, x); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
