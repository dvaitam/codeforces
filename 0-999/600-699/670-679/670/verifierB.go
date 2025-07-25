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

func expected(n int, k int64, ids []int) string {
	kk := k
	for i := 1; i <= n; i++ {
		if kk <= int64(i) {
			return fmt.Sprintf("%d", ids[kk-1])
		}
		kk -= int64(i)
	}
	return ""
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	ids := make([]int, n)
	used := make(map[int]struct{})
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(1000000000) + 1
			if _, ok := used[v]; !ok {
				ids[i] = v
				used[v] = struct{}{}
				break
			}
		}
	}
	total := int64(n * (n + 1) / 2)
	var k int64
	if total > 2000000000 {
		total = 2000000000
	}
	k = rng.Int63n(total) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range ids {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	inp := sb.String()
	exp := expected(n, k, ids)
	return inp, exp
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
