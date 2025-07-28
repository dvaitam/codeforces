package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", tag, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type testCaseE struct {
	n   int
	arr []int64
}

func genCase(rng *rand.Rand) testCaseE {
	n := rng.Intn(8) + 2
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(20) + 1
	}
	return testCaseE{n: n, arr: arr}
}

const negInf int64 = -1 << 60

func maxK(n int) int {
	return int((math.Sqrt(float64(8*n+1)) - 1) / 2)
}

func solveE(n int, a []int64) int {
	maxLen := maxK(n)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}

	dp := make([][]int64, maxLen+1)
	suf := make([][]int64, maxLen+1)

	dp[1] = make([]int64, n)
	suf[1] = make([]int64, n+1)
	for i := 0; i <= n; i++ {
		suf[1][i] = negInf
	}
	for i := n - 1; i >= 0; i-- {
		dp[1][i] = a[i]
		if dp[1][i] > suf[1][i+1] {
			suf[1][i] = dp[1][i]
		} else {
			suf[1][i] = suf[1][i+1]
		}
	}

	for l := 2; l <= maxLen; l++ {
		dp[l] = make([]int64, n)
		suf[l] = make([]int64, n+1)
		for i := 0; i <= n; i++ {
			suf[l][i] = negInf
		}
		for i := n - 1; i >= 0; i-- {
			if i+l <= n {
				sum := prefix[i+l] - prefix[i]
				if suf[l-1][i+l] > sum {
					dp[l][i] = sum
				} else {
					dp[l][i] = negInf
				}
			} else {
				dp[l][i] = negInf
			}
			if dp[l][i] > suf[l][i+1] {
				suf[l][i] = dp[l][i]
			} else {
				suf[l][i] = suf[l][i+1]
			}
		}
	}

	for k := maxLen; k >= 1; k-- {
		if suf[k][0] > negInf {
			return k
		}
	}
	return 1
}

func solveCase(tc testCaseE) string {
	return fmt.Sprintf("%d\n", solveE(tc.n, tc.arr))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.arr[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCase(tc)
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
