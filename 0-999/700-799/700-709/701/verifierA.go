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

func generateCaseA(rng *rand.Rand) (int, []int) {
	n := rng.Intn(50)*2 + 2 // even between 2 and 100
	d := rng.Intn(100) + 2
	arr := make([]int, n)
	for i := 0; i < n/2; i++ {
		a := rng.Intn(d-1) + 1
		b := d - a
		arr[2*i] = a
		arr[2*i+1] = b
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return n, arr
}

func runCaseA(bin string, n int, arr []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(arr[i]))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	used := make([]bool, n+1)
	idx := func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			return -1
		}
		return v
	}
	d := arr[idx(fields[0])-1] + arr[idx(fields[1])-1]
	for i := 0; i < n; i += 2 {
		a := idx(fields[i])
		b := idx(fields[i+1])
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("index out of range")
		}
		if used[a] || used[b] {
			return fmt.Errorf("index reused")
		}
		used[a], used[b] = true, true
		if arr[a-1]+arr[b-1] != d {
			return fmt.Errorf("pair sum mismatch")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr := generateCaseA(rng)
		if err := runCaseA(bin, n, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%v\n", i+1, err, n, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
