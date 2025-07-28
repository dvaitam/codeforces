package main

import (
	"bytes"
	"fmt"
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

type testCaseF2 struct {
	n   int
	arr []int
}

func genCase(rng *rand.Rand) testCaseF2 {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) // small numbers
	}
	return testCaseF2{n: n, arr: arr}
}

func solveF2(tc testCaseF2) string {
	const MAXV = 5000
	const LIMIT = 13
	counts := make([]int, MAXV+1)
	vals := make([]int, 0)
	for _, x := range tc.arr {
		if counts[x] < LIMIT {
			counts[x]++
			vals = append(vals, x)
		}
	}
	const MAXX = 1 << 13
	const INF = MAXV + 1
	dp := make([]int, MAXX)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for _, v := range vals {
		for x := 0; x < MAXX; x++ {
			if dp[x] < v {
				if v < dp[x^v] {
					dp[x^v] = v
				}
			}
		}
	}
	res := []int{}
	for x := 0; x < MAXX; x++ {
		if dp[x] < INF {
			res = append(res, x)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candF2")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.arr[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveF2(tc)
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
