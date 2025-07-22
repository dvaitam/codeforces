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

type testCase struct {
	arr []int
}

func expected(a []int) int {
	n := len(a)
	if n == 0 {
		return 0
	}
	dp1 := make([]int, n)
	dp2 := make([]int, n)
	dp1[0] = 1
	for i := 1; i < n; i++ {
		if a[i] > a[i-1] {
			dp1[i] = dp1[i-1] + 1
		} else {
			dp1[i] = 1
		}
	}
	dp2[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		if a[i] < a[i+1] {
			dp2[i] = dp2[i+1] + 1
		} else {
			dp2[i] = 1
		}
	}
	maxLen := 1
	for i := 0; i < n; i++ {
		if dp1[i] > maxLen {
			maxLen = dp1[i]
		}
	}
	if n > 1 {
		if dp2[1]+1 > maxLen {
			maxLen = dp2[1] + 1
		}
		if dp1[n-2]+1 > maxLen {
			maxLen = dp1[n-2] + 1
		}
	}
	for i := 1; i < n-1; i++ {
		if a[i+1]-a[i-1] >= 2 {
			l := dp1[i-1] + dp2[i+1] + 1
			if l > maxLen {
				maxLen = l
			}
		}
	}
	return maxLen
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(arr)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
