package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(n, k int64) string {
	pow := int64(1)
	for {
		count := n/pow - n/(pow*2)
		if k > count {
			k -= count
			pow <<= 1
		} else {
			ans := pow * (2*k - 1)
			return fmt.Sprint(ans)
		}
	}
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	var tests [][2]int64
	for i := 0; i < 100; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		k := rng.Int63n(n) + 1
		tests = append(tests, [2]int64{n, k})
	}
	for idx, t := range tests {
		in := fmt.Sprintf("1\n%d %d\n", t[0], t[1])
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(t[0], t[1])
		if got != exp {
			fmt.Printf("test %d failed: n=%d k=%d expected=%s got=%s\n", idx+1, t[0], t[1], exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
