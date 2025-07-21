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

type testCaseA struct {
	n   int
	arr []int
}

func generateCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(20) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000) + 1
	}
	// occasionally enforce a valid triple
	if rng.Intn(2) == 0 {
		i := rng.Intn(n)
		j := rng.Intn(n)
		for j == i {
			j = rng.Intn(n)
		}
		k := rng.Intn(n)
		for k == i || k == j {
			k = rng.Intn(n)
		}
		arr[i] = arr[j] + arr[k]
	}
	return testCaseA{n: n, arr: arr}
}

func runCaseA(bin string, tc testCaseA) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 1 && fields[0] == "-1" {
		if hasSolutionA(tc.arr) {
			return fmt.Errorf("expected a solution but got -1")
		}
		return nil
	}
	if len(fields) != 3 {
		return fmt.Errorf("expected three integers or -1, got %q", out.String())
	}
	var i, j, k int
	if _, err := fmt.Sscan(strings.Join(fields, " "), &i, &j, &k); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if i < 1 || i > tc.n || j < 1 || j > tc.n || k < 1 || k > tc.n {
		return fmt.Errorf("indices out of range")
	}
	if i == j || i == k || j == k {
		return fmt.Errorf("indices must be distinct")
	}
	arr := tc.arr
	if arr[i-1] != arr[j-1]+arr[k-1] {
		return fmt.Errorf("invalid triple")
	}
	return nil
}

func hasSolutionA(a []int) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue
				}
				if a[i] == a[j]+a[k] {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := generateCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%v\n", t+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
