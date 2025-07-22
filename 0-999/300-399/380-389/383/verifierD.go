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

const MOD = 1000000007

type testCaseD struct {
	n        int
	arr      []int
	expected string
}

func solveCase(tc testCaseD) string {
	n := tc.n
	a := tc.arr
	total := 0
	for _, v := range a {
		total += v
	}
	offset := total
	size := 2*total + 1
	dpPrev := make([]int, size)
	dpCurr := make([]int, size)
	ans := 0
	for i := 0; i < n; i++ {
		ai := a[i]
		for j := 0; j < size; j++ {
			dpCurr[j] = 0
		}
		dpCurr[offset+ai] = (dpCurr[offset+ai] + 1) % MOD
		dpCurr[offset-ai] = (dpCurr[offset-ai] + 1) % MOD
		for j := 0; j < size; j++ {
			v := dpPrev[j]
			if v != 0 {
				jp := j + ai
				if jp < size {
					dpCurr[jp] = (dpCurr[jp] + v) % MOD
				}
				jm := j - ai
				if jm >= 0 {
					dpCurr[jm] = (dpCurr[jm] + v) % MOD
				}
			}
		}
		ans = (ans + dpCurr[offset]) % MOD
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	return strconv.Itoa(ans)
}

func run(bin string, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(5) + 1
		}
		tc := testCaseD{n: n, arr: arr}
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[j]))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
