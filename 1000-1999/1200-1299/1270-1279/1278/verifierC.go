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

func runCandidate(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(arr []int) int {
	n := len(arr) / 2
	for i := range arr {
		if arr[i] == 2 {
			arr[i] = -1
		} else {
			arr[i] = 1
		}
	}
	total := 0
	for _, v := range arr {
		total += v
	}
	right := map[int]int{0: 0}
	diff := 0
	for i := 0; i < n; i++ {
		diff += arr[n+i]
		if _, ok := right[diff]; !ok {
			right[diff] = i + 1
		}
	}
	ans := 2 * n
	if val, ok := right[total]; ok && val < ans {
		ans = val
	}
	diffLeft := 0
	for j := 1; j <= n; j++ {
		diffLeft += arr[n-j]
		need := total - diffLeft
		if k, ok := right[need]; ok {
			if j+k < ans {
				ans = j + k
			}
		}
	}
	return ans
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(50) + 1
	arr := make([]int, 2*n)
	for i := range arr {
		if r.Intn(2) == 0 {
			arr[i] = 1
		} else {
			arr[i] = 2
		}
	}
	expect := fmt.Sprintf("%d", solveC(append([]int(nil), arr...)))
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
