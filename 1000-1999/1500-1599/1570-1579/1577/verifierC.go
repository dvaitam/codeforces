package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isSorted(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		if rng.Intn(2) == 0 {
			base := rng.Intn(100)
			for j := 0; j < n; j++ {
				arr[j] = base + j*rng.Intn(3)
			}
		} else {
			for j := 0; j < n; j++ {
				arr[j] = rng.Intn(100)
			}
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expect := "NO"
		if isSorted(arr) {
			expect = "YES"
		}
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(out)) != expect {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", i, expect, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
