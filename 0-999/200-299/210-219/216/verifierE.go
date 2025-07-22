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

type testCaseE struct {
	k int64
	b int64
	a []int64
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func generateCase() testCaseE {
	k := int64(rand.Intn(9) + 2) // 2..10
	b := int64(rand.Intn(int(k)))
	n := rand.Intn(20) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rand.Intn(int(k)))
	}
	return testCaseE{k: k, b: b, a: arr}
}

func compute(tc testCaseE) int64 {
	k := tc.k
	b := tc.b
	n := len(tc.a)
	P := make([]int64, n+1)
	for i := 0; i < n; i++ {
		P[i+1] = P[i] + tc.a[i]
	}
	M := k - 1
	if b == 0 {
		cnt := make(map[int64]int64)
		for _, v := range P {
			cnt[v]++
		}
		var result int64
		for _, c := range cnt {
			if c > 1 {
				result += c * (c - 1) / 2
			}
		}
		return result
	}
	if b == M {
		cntMod := make(map[int64]int64)
		cntExact := make(map[int64]int64)
		cntMod[P[0]%M] = 1
		cntExact[P[0]] = 1
		var totalMod, totalExact int64
		for j := 1; j <= n; j++ {
			r := P[j] % M
			totalMod += cntMod[r]
			totalExact += cntExact[P[j]]
			cntMod[r]++
			cntExact[P[j]]++
		}
		return totalMod - totalExact
	}
	cntMod := make(map[int64]int64)
	cntMod[P[0]%M] = 1
	var result int64
	for j := 1; j <= n; j++ {
		r := P[j] % M
		need := r - b
		need %= M
		if need < 0 {
			need += M
		}
		result += cntMod[need]
		cntMod[r]++
	}
	return result
}

func buildInput(tc testCaseE) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.k, tc.b, len(tc.a)))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCase()
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", i)
			os.Exit(1)
		}
		exp := compute(tc)
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
