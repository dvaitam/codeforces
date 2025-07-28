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

func rotateLeft(arr []int, l, r, d int) {
	d = d % (r - l)
	if d == 0 {
		return
	}
	tmp := append([]int(nil), arr[l:l+d]...)
	copy(arr[l:r-d], arr[l+d:r])
	copy(arr[r-d:r], tmp)
}

func runCaseB(bin string, n int, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k > n {
		return fmt.Errorf("k too large")
	}
	if len(fields) != 1+3*k {
		return fmt.Errorf("expected %d numbers got %d", 1+3*k, len(fields))
	}
	a := append([]int(nil), arr...)
	idx := 1
	for op := 0; op < k; op++ {
		l, _ := strconv.Atoi(fields[idx])
		r, _ := strconv.Atoi(fields[idx+1])
		d, _ := strconv.Atoi(fields[idx+2])
		idx += 3
		if l < 1 || l >= r || r > n {
			return fmt.Errorf("invalid operation")
		}
		rotateLeft(a, l-1, r, d)
	}
	for i := 0; i+1 < n; i++ {
		if a[i] > a[i+1] {
			return fmt.Errorf("array not sorted")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(21) - 10
		}
		if err := runCaseB(bin, n, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
