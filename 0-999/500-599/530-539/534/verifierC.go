package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, A int64, d []int64) []int64 {
	sumd := int64(0)
	for _, v := range d {
		sumd += v
	}
	sumMinExcl := int64(n - 1)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		low := A - (sumd - d[i])
		if low < 1 {
			low = 1
		}
		high := A - sumMinExcl
		if high > d[i] {
			high = d[i]
		}
		var possible int64
		if low > high {
			possible = 0
		} else {
			possible = high - low + 1
		}
		res[i] = d[i] - possible
	}
	return res
}

func check(n int, out string, exp []int64) error {
	fields := strings.Fields(out)
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	for i := 0; i < n; i++ {
		val, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return fmt.Errorf("bad number")
		}
		if val != exp[i] {
			return fmt.Errorf("expected %d got %d at pos %d", exp[i], val, i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		d := make([]int64, n)
		for j := 0; j < n; j++ {
			d[j] = int64(rng.Intn(6) + 1) // dice faces 1..6 to keep numbers small
		}
		sumd := int64(0)
		for _, v := range d {
			sumd += v
		}
		A := int64(rng.Intn(int(sumd-int64(n)+1)) + n) // ensure n <= A <= sumd
		exp := expected(n, A, d)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, A))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(d[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(n, out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
