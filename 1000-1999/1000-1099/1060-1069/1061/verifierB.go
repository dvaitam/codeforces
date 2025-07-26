package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func solveCase(n int, m int64, arr []int64) int64 {
	orig := make([]int64, n)
	copy(orig, arr)
	sum := int64(0)
	for _, v := range orig {
		sum += v
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var maxx int64
	for i := n - 1; i > 0; i-- {
		if arr[i-1] >= arr[i]-1 {
			if arr[i] <= 1 {
				arr[i-1] = 1
			} else {
				arr[i-1] = arr[i] - 1
			}
			maxx++
		} else {
			maxx += arr[i] - arr[i-1]
		}
	}
	maxx += arr[0]
	return sum - maxx
}

func genCase(rng *rand.Rand) (int, int64, []int64) {
	n := rng.Intn(20) + 1
	m := int64(rng.Intn(50) + 1)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(int(m)) + 1)
	}
	return n, m, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, arr := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		expect := solveCase(n, m, append([]int64(nil), arr...))
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil || got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
