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

func expected(n, k int, a []int) int {
	maxA := 0
	minA := int(1e9)
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
		if v < minA {
			minA = v
		}
	}
	cnt := make([]int, maxA+1)
	for _, v := range a {
		cnt[v]++
	}
	ps := make([]int, maxA+1)
	for i := 1; i <= maxA; i++ {
		ps[i] = ps[i-1] + cnt[i]
	}
	ans := 1
	for d := 2; d <= minA; d++ {
		ok := true
		for m := 0; m <= maxA; m += d {
			low := m + k + 1
			if low > maxA {
				break
			}
			high := m + d - 1
			if high > maxA {
				high = maxA
			}
			if low <= high && ps[high]-ps[low-1] > 0 {
				ok = false
				break
			}
		}
		if ok {
			ans = d
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(20) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(50) + 1
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j, v := range a {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expect := strconv.Itoa(expected(n, k, a))
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
