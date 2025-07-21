package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(s string) string {
	n := len(s)
	dpPrev := make([]*big.Int, 10)
	dpCur := make([]*big.Int, 10)
	for i := 0; i < 10; i++ {
		dpPrev[i] = big.NewInt(1)
		dpCur[i] = big.NewInt(0)
	}
	for i := 1; i < n; i++ {
		for d := 0; d < 10; d++ {
			dpCur[d].SetInt64(0)
		}
		ai := int(s[i] - '0')
		for prev := 0; prev < 10; prev++ {
			count := dpPrev[prev]
			if count.Sign() == 0 {
				continue
			}
			sum := ai + prev
			if sum%2 == 0 {
				next := sum / 2
				dpCur[next].Add(dpCur[next], count)
			} else {
				low := sum / 2
				high := low + 1
				dpCur[low].Add(dpCur[low], count)
				dpCur[high].Add(dpCur[high], count)
			}
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	res := big.NewInt(0)
	for d := 0; d < 10; d++ {
		res.Add(res, dpPrev[d])
	}
	return res.String()
}

func randDigits(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('0' + rng.Intn(10))
	}
	return string(b)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	s := randDigits(rng, n)
	return testCase{input: s + "\n", expected: solveCase(s)}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "0\n", expected: "1"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
