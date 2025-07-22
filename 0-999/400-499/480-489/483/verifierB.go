package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	cnt1, cnt2, x, y int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func good(v, cnt1, cnt2, x, y int64) bool {
	s1 := v - v/x
	s2 := v - v/y
	s12 := v - v/lcm(x, y)
	return s1 >= cnt1 && s2 >= cnt2 && s12 >= cnt1+cnt2
}

func solveB(cnt1, cnt2, x, y int64) string {
	lo, hi := int64(0), int64(1)
	for !good(hi, cnt1, cnt2, x, y) {
		hi <<= 1
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if good(mid, cnt1, cnt2, x, y) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return fmt.Sprintf("%d\n", lo)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func generateTests() []testCaseB {
	tests := make([]testCaseB, 0, 100)
	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for i := 0; i < 90; i++ {
		cnt1 := int64(i%20 + 1)
		cnt2 := int64((i*3)%20 + 1)
		x := primes[i%len(primes)]
		y := primes[(i+3)%len(primes)]
		if x == y {
			y = primes[(i+4)%len(primes)]
		}
		tests = append(tests, testCaseB{cnt1: cnt1, cnt2: cnt2, x: x, y: y})
	}
	base := int64(1_000_000_000)
	for i := 0; i < 10; i++ {
		tests = append(tests, testCaseB{cnt1: base + int64(i), cnt2: base + int64(2*i), x: 29989, y: 29989 + 2})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d %d\n", t.cnt1, t.cnt2, t.x, t.y)
		expected := solveB(t.cnt1, t.cnt2, t.x, t.y)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed: cnt1=%d cnt2=%d x=%d y=%d expected %q got %q\n", i+1, t.cnt1, t.cnt2, t.x, t.y, strings.TrimSpace(expected), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
