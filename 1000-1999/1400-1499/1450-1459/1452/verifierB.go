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

type testB struct {
	arr []int64
}

func solveB(arr []int64) int64 {
	n := int64(len(arr))
	var sum, mx int64
	for _, v := range arr {
		sum += v
		if v > mx {
			mx = v
		}
	}
	need1 := mx*(n-1) - sum
	if need1 < 0 {
		need1 = 0
	}
	sum += need1
	rem := sum % (n - 1)
	var need2 int64
	if rem != 0 {
		need2 = (n - 1) - rem
	}
	return need1 + need2
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	tests := make([]testB, t)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Int63n(20)
		}
		tests[i].arr = arr
		fmt.Fprintf(&b, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", arr[j])
		}
		b.WriteByte('\n')
	}
	out, err := runBinary(binary, b.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != t {
		fmt.Printf("expected %d lines, got %d\noutput:\n%s\n", t, len(fields), out)
		os.Exit(1)
	}
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			fmt.Printf("invalid output on test %d: %s\n", i+1, f)
			os.Exit(1)
		}
		exp := solveB(tests[i].arr)
		if v != exp {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, exp, v)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
