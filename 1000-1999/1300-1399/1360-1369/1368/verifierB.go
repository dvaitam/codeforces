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

const target = "codeforces"

func runCmd(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	k := rand.Int63n(1e12) + 1
	return []byte(fmt.Sprintf("%d\n", k))
}

func mulCap(a, b, cap int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	if a > cap/b {
		return cap
	}
	v := a * b
	if v > cap {
		return cap
	}
	return v
}

func maxSubseqForLen(n int, k int64) int64 {
	q := n / len(target)
	r := n % len(target)
	prod := int64(1)
	for i := 0; i < len(target)-r; i++ {
		prod = mulCap(prod, int64(q), k)
	}
	for i := 0; i < r; i++ {
		prod = mulCap(prod, int64(q+1), k)
	}
	return prod
}

func minLenForK(k int64) int {
	n := len(target)
	for maxSubseqForLen(n, k) < k {
		n++
	}
	return n
}

func countSubseq(s string, k int64) int64 {
	dp := make([]int64, len(target)+1)
	dp[0] = 1
	for i := 0; i < len(s); i++ {
		ch := s[i]
		for j := len(target) - 1; j >= 0; j-- {
			if target[j] == ch {
				dp[j+1] += dp[j]
				if dp[j+1] > k {
					dp[j+1] = k
				}
			}
		}
	}
	return dp[len(target)]
}

func isValidAnswer(k int64, out string) error {
	s := strings.TrimSpace(out)
	if s == "" {
		return fmt.Errorf("empty output")
	}
	if countSubseq(s, k) < k {
		return fmt.Errorf("contains fewer than k codeforces subsequences")
	}
	minLen := minLenForK(k)
	if len(s) != minLen {
		return fmt.Errorf("output is not shortest: got len %d, expected %d", len(s), minLen)
	}
	return nil
}

func main() {
	var cand string
	if len(os.Args) == 2 {
		cand = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		cand = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		var k int64
		if _, err := fmt.Sscan(string(input), &k); err != nil {
			fmt.Println("internal error parsing generated input:", err)
			os.Exit(1)
		}

		got, err := runCmd(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if err := isValidAnswer(k, got); err != nil {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("got:\n", got)
			fmt.Println("reason:", err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
