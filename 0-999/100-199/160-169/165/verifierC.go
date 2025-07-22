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

func expected(k int, s string) int64 {
	n := len(s)
	if k == 0 {
		var res int64
		var zeroLen int64
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zeroLen++
			} else {
				res += zeroLen * (zeroLen + 1) / 2
				zeroLen = 0
			}
		}
		res += zeroLen * (zeroLen + 1) / 2
		return res
	}
	positions := []int{-1}
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			positions = append(positions, i)
		}
	}
	positions = append(positions, n)
	cnt := len(positions) - 2
	if k > cnt {
		return 0
	}
	var res int64
	for i := 1; i <= cnt-k+1; i++ {
		leftZeros := int64(positions[i] - positions[i-1])
		rightZeros := int64(positions[i+k] - positions[i+k-1])
		res += leftZeros * rightZeros
	}
	return res
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(50) + 1
	if rng.Float64() < 0.2 {
		n = rng.Intn(200) + 1
	}
	bs := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bs[i] = '0'
		} else {
			bs[i] = '1'
		}
	}
	k := rng.Intn(n + 1)
	input := fmt.Sprintf("%d\n%s\n", k, string(bs))
	return input, expected(k, string(bs))
}

func runCase(bin, input string, exp int64) error {
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
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
