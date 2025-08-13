package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(10) + 1
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", perm[i]))
	}
	sb.WriteString("\n")
	return []byte(sb.String())
}

func verify(in []byte, out string) error {
	inReader := bufio.NewReader(bytes.NewReader(in))
	outReader := bufio.NewReader(strings.NewReader(out))
	var t int
	if _, err := fmt.Fscan(inReader, &t); err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	for tc := 1; tc <= t; tc++ {
		var n int
		if _, err := fmt.Fscan(inReader, &n); err != nil {
			return fmt.Errorf("failed to read n: %v", err)
		}
		a := make([]int, n)
		for i := range a {
			if _, err := fmt.Fscan(inReader, &a[i]); err != nil {
				return fmt.Errorf("failed to read permutation: %v", err)
			}
		}
		x := make([]int, n)
		for i := range x {
			if _, err := fmt.Fscan(outReader, &x[i]); err != nil {
				return fmt.Errorf("not enough output for test case %d", tc)
			}
			if x[i] < 1 || x[i] > n {
				return fmt.Errorf("x[%d]=%d out of range", i+1, x[i])
			}
		}
		for i := 0; i < n; i++ {
			for j := x[i] - 1; j < n; j++ {
				a[j] += i + 1
			}
		}
		for i := 0; i < n-1; i++ {
			if a[i] > a[i+1] {
				return fmt.Errorf("array not sorted after operations")
			}
		}
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	for i := 1; i <= 100; i++ {
		in := genTest()
		got, err := runBinary(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := verify(in, got); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%sreason: %v\noutput:%s\n", i, string(in), err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
