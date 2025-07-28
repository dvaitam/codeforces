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

func expected(T []int64, queries [][]int64) string {
	lastans := int64(0)
	var sb strings.Builder
	for i := 0; i < len(T); i++ {
		for _, xp := range queries[i] {
			x := (xp + lastans) % (1_000_000_000 + 1)
			cur := x
			for d := 0; d <= i; d++ {
				if cur < T[d] {
					cur++
				} else if cur > T[d] {
					cur--
				}
			}
			fmt.Fprintln(&sb, cur)
			lastans = cur
		}
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	T := make([]int64, n)
	queries := make([][]int64, n)
	for i := 0; i < n; i++ {
		T[i] = int64(rng.Intn(100))
		k := rng.Intn(5)
		queries[i] = make([]int64, k)
		for j := 0; j < k; j++ {
			queries[i][j] = rng.Int63n(1_000_000_000)
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteString(strconv.FormatInt(T[i], 10))
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(len(queries[i])))
		sb.WriteByte('\n')
		for j, v := range queries[i] {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	inp := sb.String()
	exp := expected(T, queries)
	return inp, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
