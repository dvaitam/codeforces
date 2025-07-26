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

func maxSum(len, H int64) int64 {
	if len <= H {
		return len * (len + 1) / 2
	}
	k := (len - H + 2) / 2
	return k*H + k*(k-1)/2 + (len-k)*(len-k+1)/2
}

func expectedD(n, H int64) int64 {
	lo, hi := int64(1), int64(2000000000)
	for lo < hi {
		mid := (lo + hi) / 2
		if maxSum(mid, H) >= n {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func generateCaseD(rng *rand.Rand) (int64, int64) {
	n := rng.Int63n(1000000) + 1
	H := rng.Int63n(1000000) + 1
	return n, H
}

func runCaseD(bin string, n, H int64) error {
	input := fmt.Sprintf("%d %d\n", n, H)
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
	expected := expectedD(n, H)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
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
		n, H := generateCaseD(rng)
		if err := runCaseD(bin, n, H); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%d %d\n", i+1, err, n, H)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
