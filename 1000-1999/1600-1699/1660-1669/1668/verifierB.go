package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveB implements the official solution for 1668B.
func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		var sum, mn, mx int64
		mn = math.MaxInt64
		for i := int64(0); i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
			if a[i] < mn {
				mn = a[i]
			}
			if a[i] > mx {
				mx = a[i]
			}
		}
		if n > m {
			sb.WriteString("NO\n")
			continue
		}
		need := n + sum + mx - mn
		if need <= m {
			sb.WriteString("YES\n")
		} else {
			sb.WriteString("NO\n")
		}
	}
	return sb.String()
}

func buildCaseB(n, m int64, arr []int64) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateRandomCaseB(rng *rand.Rand) string {
	n := rng.Intn(8) + 2 // 2..9
	m := int64(rng.Intn(100) + 1)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(100) + 1)
	}
	return buildCaseB(int64(n), m, arr)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []string
	// predetermined edge cases
	cases = append(cases, buildCaseB(3, 2, []int64{1, 1, 1}))
	cases = append(cases, buildCaseB(2, 4, []int64{1, 1}))
	cases = append(cases, buildCaseB(2, 5, []int64{1, 2}))
	cases = append(cases, buildCaseB(3, 10, []int64{2, 2, 2}))
	cases = append(cases, buildCaseB(3, 3, []int64{1, 1, 1}))
	for len(cases) < 100 {
		cases = append(cases, generateRandomCaseB(rng))
	}
	for i, tc := range cases {
		expect := strings.TrimSpace(solveB(tc))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
