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

func lis(a []int) int {
	d := make([]int, 0, len(a))
	for _, v := range a {
		i := sort.SearchInts(d, v)
		if i == len(d) {
			d = append(d, v)
		} else {
			d[i] = v
		}
	}
	return len(d)
}

func expectedD(n int, arr []int, t int) int {
	a := make([]int, 0, n*t)
	for i := 0; i < t; i++ {
		a = append(a, arr...)
	}
	return lis(a)
}

func genCaseD(rng *rand.Rand) (string, string) {
	k := rng.Intn(3) + 1
	n := rng.Intn(5) + 1
	maxb := 20
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", k, n, maxb, t)
	var out strings.Builder
	for i := 0; i < k; i++ {
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(maxb) + 1
		}
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&out, "%d", expectedD(n, arr, t))
		if i+1 < k {
			out.WriteByte('\n')
		}
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
