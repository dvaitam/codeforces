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

type pair struct {
	val int64
	idx int
}

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

func expected(arr []int64) string {
	n := len(arr)
	ps := make([]pair, n)
	for i := 0; i < n; i++ {
		ps[i] = pair{arr[i], i + 1}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val > ps[j].val })
	ans := make([]int64, n+1)
	var sum int64
	for i, p := range ps {
		k := int64(i/2 + 1)
		if i%2 == 0 {
			ans[p.idx] = k
		} else {
			ans[p.idx] = -k
		}
		sum += 2 * p.val * k
	}
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(sum, 10))
	sb.WriteByte('\n')
	sb.WriteString("0")
	for i := 1; i <= n; i++ {
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(ans[i], 10))
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1_000_000 + 1)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	inp := sb.String()
	exp := expected(arr)
	return inp, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
