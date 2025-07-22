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

func run(bin, input string) (string, error) {
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

type query struct{ l, r int }

func expected(arr []int, qs []query) string {
	cnt1 := 0
	for _, v := range arr {
		if v == 1 {
			cnt1++
		}
	}
	cntNeg := len(arr) - cnt1
	var sb strings.Builder
	for i, q := range qs {
		length := q.r - q.l + 1
		res := 0
		if length%2 == 0 && cnt1 >= length/2 && cntNeg >= length/2 {
			res = 1
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(res))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := range arr {
			if rng.Intn(2) == 0 {
				arr[j] = -1
			} else {
				arr[j] = 1
			}
		}
		m := rng.Intn(20) + 1
		qs := make([]query, m)
		for j := 0; j < m; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			qs[j] = query{l, r}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, q := range qs {
			fmt.Fprintf(&sb, "%d %d\n", q.l, q.r)
		}
		input := sb.String()
		expect := expected(arr, qs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n got:\n%s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
