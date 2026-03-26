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

// ─── embedded correct solver ───

type pairF struct {
	idx int
	w   int
	c   int
}

func solveF(input string) string {
	data := []byte(input)
	pos := 0

	readInt := func() int {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		if pos >= len(data) {
			return 0
		}
		res := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			res = res*10 + int(data[pos]-'0')
			pos++
		}
		return res
	}

	n := readInt()
	k := readInt()
	l := readInt()

	for pos < len(data) && data[pos] <= ' ' {
		pos++
	}
	start := pos
	for pos < len(data) && data[pos] > ' ' {
		pos++
	}
	s := data[start:pos]

	if int64(k)*int64(l) >= int64(n) {
		return "0"
	}

	ones := make([]int, n+1)
	dp := make([]int, n+1)
	cnt := make([]int, n+1)
	q := make([]pairF, n+1)

	solve := func(targetLower bool) int {
		for i := 0; i < n; i++ {
			ones[i+1] = ones[i]
			isLower := s[i] >= 'a' && s[i] <= 'z'
			if isLower == targetLower {
				ones[i+1]++
			}
		}

		low, high := 0, l+1
		bestCov := 0

		for low <= high {
			mid := (low + high) / 2

			head, tail := 0, 0
			q[tail] = pairF{0, 0, 0}
			tail++

			for i := 1; i <= n; i++ {
				for head < tail && q[head].idx < i-l {
					head++
				}

				dp[i] = dp[i-1]
				cnt[i] = cnt[i-1]

				if head < tail {
					best := q[head]
					val := best.w + ones[i] - mid
					c := best.c + 1
					if val > dp[i] || (val == dp[i] && c < cnt[i]) {
						dp[i] = val
						cnt[i] = c
					}
				}

				w := dp[i] - ones[i]
				c := cnt[i]
				for head < tail {
					last := q[tail-1]
					if last.w < w || (last.w == w && last.c >= c) {
						tail--
					} else {
						break
					}
				}
				q[tail] = pairF{i, w, c}
				tail++
			}

			if cnt[n] <= k {
				bestCov = dp[n] + mid*k
				high = mid - 1
			} else {
				low = mid + 1
			}
		}

		return ones[n] - bestCov
	}

	ans1 := solve(true)
	ans2 := solve(false)

	if ans1 < ans2 {
		return strconv.Itoa(ans1)
	}
	return strconv.Itoa(ans2)
}

// ─── verifier ───

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func deterministicCases() []string {
	return []string{
		"1 1 1\na\n",
		"3 2 2\nAbC\n",
		"5 3 1\nabcDE\n",
	}
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = byte('a' + rng.Intn(26))
		} else {
			b[i] = byte('A' + rng.Intn(26))
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	l := rng.Intn(n) + 1
	s := randString(rng, n)
	return fmt.Sprintf("%d %d %d\n%s\n", n, k, l, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, in := range cases {
		want := solveF(in)
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\n got: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
