package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (previously in testcasesB.txt) to keep verifier self contained.
const rawTestcasesB = `
5
19
3
9
4
16
15
16
13
7
4
16
1
13
14
20
1
15
9
8
19
4
11
1
1
1
18
1
13
7
14
1
17
8
15
16
18
8
12
8
8
15
10
1
14
18
4
6
10
4
11
17
14
17
7
10
10
19
16
17
13
19
2
16
8
13
14
6
12
18
12
3
15
17
4
6
17
13
12
16
1
16
2
10
20
19
19
13
6
6
17
8
1
7
18
18
8
13
17
12
`

// modPow computes a^b mod m.
func modPow(a, b, m int) int {
	res := 1 % m
	a %= m
	for b > 0 {
		if b&1 == 1 {
			res = (res * a) % m
		}
		a = (a * a) % m
		b >>= 1
	}
	return res
}

// solve248B mirrors 248B.go to compute expected output for a given n.
func solve248B(n int) string {
	if n < 3 {
		return "-1"
	}
	rem := modPow(10, n-1, 210)
	remNeeded := (210 - rem) % 210
	digits := make([]int, n)
	digits[0] = 1
	t := remNeeded
	pos := n - 1
	for t > 0 && pos >= 0 {
		digits[pos] += t % 10
		t /= 10
		pos--
	}
	for i := n - 1; i > 0; i-- {
		if digits[i] >= 10 {
			carry := digits[i] / 10
			digits[i] %= 10
			digits[i-1] += carry
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('0' + digits[i]))
	}
	return sb.String()
}

func loadTestcases() ([]int, error) {
	lines := strings.Fields(rawTestcasesB)
	nums := make([]int, 0, len(lines))
	for idx, s := range lines {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("token %d (%q): %w", idx+1, s, err)
		}
		nums = append(nums, v)
	}
	return nums, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, n := range testcases {
		expect := solve248B(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
