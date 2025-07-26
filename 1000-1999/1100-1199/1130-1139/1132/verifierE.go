package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testCaseE struct {
	W   int64
	cnt [9]int64
}

func genTestsE() []testCaseE {
	rand.Seed(113205)
	tests := make([]testCaseE, 100)
	for i := range tests {
		W := int64(rand.Intn(50) + 1)
		var cnt [9]int64
		for j := 1; j <= 8; j++ {
			cnt[j] = int64(rand.Intn(5))
		}
		tests[i] = testCaseE{W: W, cnt: cnt}
	}
	return tests
}

func solveE(tc testCaseE) int64 {
	targetW := tc.W
	cnt := tc.cnt
	var sum int64
	for i := 1; i <= 8; i++ {
		sum += cnt[i] * int64(i)
	}
	if sum <= targetW {
		return sum
	}
	var now int64
	var cnt2 [9]int64
	for i := 1; i <= 8; i++ {
		w := int64(i)
		if now+w*cnt[i] <= targetW {
			now += w * cnt[i]
			cnt2[i] = cnt[i]
			cnt[i] = 0
		} else {
			v := (targetW - now) / w
			if v < 0 {
				v = 0
			}
			now += w * v
			cnt[i] -= v
			cnt2[i] += v
		}
	}
	const AA = 1000
	var dp [2001]bool
	dp[AA] = true
	for v := 1; v <= 8; v++ {
		times := cnt[v]
		if times > 100 {
			times = 100
		}
		for j := int64(0); j < times; j++ {
			for k := 900; k >= -900; k-- {
				if dp[k+AA] {
					dp[k+v+AA] = true
				}
			}
		}
		times2 := cnt2[v]
		if times2 > 100 {
			times2 = 100
		}
		for j := int64(0); j < times2; j++ {
			for k := -900; k <= 900; k++ {
				if dp[k+AA] {
					dp[k-int(v)+AA] = true
				}
			}
		}
	}
	var ans int64
	maxR := targetW - now
	if maxR < 0 {
		maxR = 0
	}
	if maxR > 2000 {
		maxR = 2000
	}
	for i := int64(0); i <= maxR; i++ {
		if dp[int(i)+AA] {
			ans = now + i
		}
	}
	return ans
}

func run(bin string, in []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, tc.W)
		for i := 1; i <= 8; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.cnt[i])
		}
		input.WriteByte('\n')
		out, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "no output on test %d\n", idx+1)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
			os.Exit(1)
		}
		expected := solveE(tc)
		if val != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", idx+1, expected, val)
			os.Exit(1)
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "extra output on test %d\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
