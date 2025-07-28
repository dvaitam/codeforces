package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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

func expected(n int, l, r, k int64, arr []int64) string {
	filtered := make([]int64, 0, n)
	for _, v := range arr {
		if v >= l && v <= r {
			filtered = append(filtered, v)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i] < filtered[j] })
	count := 0
	budget := k
	for _, p := range filtered {
		if p <= budget {
			budget -= p
			count++
		} else {
			break
		}
	}
	return strconv.Itoa(count)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	l := rng.Int63n(1_000_000_000) + 1
	var r int64
	if l == 1_000_000_000 {
		r = 1_000_000_000
	} else {
		r = l + rng.Int63n(1_000_000_000-l+1)
	}
	k := rng.Int63n(1_000_000_000) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1_000_000_000) + 1
	}
	exp := expected(n, l, r, k, arr)

	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, l, r, k)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')

	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
