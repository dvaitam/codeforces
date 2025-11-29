package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func checkAndCount(d, L int, chars []byte, cnt []int) int {
	k := L / d

	// Check if chars is periodic with period d
	for i := 0; i < L-d; i++ {
		if chars[i] != chars[i+d] {
			return 0
		}
	}

	// Check consistency of internal gaps
	for r := 1; r < d; r++ {
		target := cnt[r]
		for j := 1; j < k; j++ {
			if cnt[j*d+r] != target {
				return 0
			}
		}
	}

	X := cnt[0]
	Y := cnt[L]

	if k == 1 {
		return (X + 1) * (Y + 1)
	}

	// Determine K: the minimum gap available between consecutive instances of t
	minGap := cnt[d]
	for j := 2; j < k; j++ {
		if cnt[j*d] < minGap {
			minGap = cnt[j*d]
		}
	}
	K := minGap

	M := X
	if K < M {
		M = K
	}
	if M < 0 {
		return 0
	}

	res := 0
	x_cut := K - Y

	// Part 1: x <= x_cut
	up1 := M
	if x_cut < up1 {
		up1 = x_cut
	}
	if up1 >= 0 {
		res += (up1 + 1) * (Y + 1)
	}

	// Part 2: x > x_cut
	low2 := 0
	if x_cut+1 > low2 {
		low2 = x_cut + 1
	}
	if low2 <= M {
		count := M - low2 + 1
		sumX := (low2 + M) * count / 2
		res += count*(K+1) - sumX
	}

	return res
}

func solveReference(s string) int {
	n := len(s)
	var chars []byte
	var cnt []int
	curr := 0
	for i := 0; i < n; i++ {
		if s[i] == 'a' {
			curr++
		} else {
			cnt = append(cnt, curr)
			chars = append(chars, s[i])
			curr = 0
		}
	}
	cnt = append(cnt, curr)

	L := len(chars)
	if L == 0 {
		return n - 1
	}

	ans := 0
	for d := 1; d*d <= L; d++ {
		if L%d == 0 {
			ans += checkAndCount(d, L, chars, cnt)
			if d*d != L {
				ans += checkAndCount(L/d, L, chars, cnt)
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 100; i++ {
		l := rng.Intn(50) + 5
		s := randomString(rng, l)
		input := fmt.Sprintf("1\n%s\n", s)

		exp := solveReference(s)

		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}

		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %s\n", i, gotStr)
			os.Exit(1)
		}

		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %d\ngot: %d\n", i, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
