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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(arrs [][]int64, w int) []int64 {
	ans := make([]int64, w)
	for col := 0; col < w; col++ {
		var sum int64
		for _, arr := range arrs {
			l := len(arr)
			best := int64(0)
			for shift := 0; shift <= w-l; shift++ {
				idx := col - shift
				var v int64
				if idx >= 0 && idx < l {
					v = arr[idx]
				} else {
					v = 0
				}
				if v > best {
					best = v
				}
			}
			sum += best
		}
		ans[col] = sum
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	w := rng.Intn(5) + 1
	arrs := make([][]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, w))
	for i := 0; i < n; i++ {
		l := rng.Intn(w) + 1
		arr := make([]int64, l)
		for j := 0; j < l; j++ {
			arr[j] = int64(rng.Intn(11) - 5)
		}
		arrs[i] = arr
		sb.WriteString(fmt.Sprintf("%d", l))
		for j := 0; j < l; j++ {
			sb.WriteString(fmt.Sprintf(" %d", arr[j]))
		}
		sb.WriteByte('\n')
	}
	res := solveCase(arrs, w)
	var exp strings.Builder
	for i, v := range res {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
