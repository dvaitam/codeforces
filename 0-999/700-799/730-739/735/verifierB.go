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

func solveB(n, n1, n2 int, a []int) float64 {
	tmp := append([]int(nil), a...)
	sort.Ints(tmp)
	t := n1
	if n2 < t {
		t = n2
	}
	sum1 := 0
	for i := n - t; i < n; i++ {
		sum1 += tmp[i]
	}
	avg1 := float64(sum1) / float64(t)
	t1 := n1
	if n2 > t1 {
		t1 = n2
	}
	sum2 := 0
	for i := n - t - t1; i < n-t; i++ {
		sum2 += tmp[i]
	}
	avg2 := float64(sum2) / float64(t1)
	return avg1 + avg2
}

func genTestB() (int, int, int, []int) {
	n := rand.Intn(40) + 2 // 2..41
	n1 := rand.Intn(n-1) + 1
	n2 := rand.Intn(n-n1) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(100000) + 1
	}
	return n, n1, n2, a
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]
	for i := 0; i < 100; i++ {
		n, n1, n2, a := genTestB()
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d %d\n", n, n1, n2)
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", a[j])
		}
		b.WriteByte('\n')
		input := b.String()
		expected := solveB(n, n1, n2, a)
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		var got float64
		_, err = fmt.Sscanf(gotStr, "%f", &got)
		if err != nil {
			fmt.Printf("test %d: parse output error: %v\ninput:\n%s\noutput:%s\n", i+1, err, input, gotStr)
			os.Exit(1)
		}
		if diff := got - expected; diff < -1e-6 || diff > 1e-6 {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %.6f\ngot: %s\n", i+1, input, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
