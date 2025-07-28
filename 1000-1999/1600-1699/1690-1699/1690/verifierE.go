package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solveCase(n int, k int64, arr []int64) int64 {
	var res int64
	rems := make([]int64, n)
	for i, v := range arr {
		res += v / k
		rems[i] = v % k
	}
	sort.Slice(rems, func(i, j int) bool { return rems[i] < rems[j] })
	i, j := 0, n-1
	for i < j {
		if rems[i]+rems[j] >= k {
			res++
			i++
			j--
		} else {
			i++
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(46))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 2
		k := int64(rng.Intn(10) + 1)
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = int64(rng.Intn(50))
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCase(n, k, arr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscanf(strings.TrimSpace(out), "%d", &got)
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", tc+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
