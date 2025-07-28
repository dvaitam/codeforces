package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func simulate(n int) int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	sort.Ints(arr)
	for len(arr) > 1 {
		mn := arr[0]
		mx := arr[len(arr)-1]
		arr = arr[1 : len(arr)-1]
		d := mx - mn
		i := sort.SearchInts(arr, d)
		arr = append(arr[:i], append([]int{d}, arr[i:]...)...)
	}
	return arr[0]
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprint(simulate(n))
	if got != exp {
		return fmt.Errorf("n=%d expected %s got %s", n, exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(15) + 2
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
