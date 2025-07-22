package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func binom(n, k int) int64 {
	if k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := 1; i <= k; i++ {
		res = res * int64(n-k+i) / int64(i)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(30)
		k := rand.Intn(n + 1)
		expected := binom(n, k)

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d %d\n", n, k))
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		gotStr := string(bytes.TrimSpace(out))
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: n=%d k=%d expected %d got %s\n", i+1, n, k, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
