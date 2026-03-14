package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(arr []int) string {
	n := len(arr)
	vals := make([]int, n)
	for i, v := range arr {
		vals[i] = (v + 1) / 2
	}
	ans := 1 << 60
	// Two independent sections (distance >= 3): use prefix min of vals
	if n > 0 {
		prefMin := make([]int, n)
		prefMin[0] = vals[0]
		for i := 1; i < n; i++ {
			prefMin[i] = prefMin[i-1]
			if vals[i] < prefMin[i] {
				prefMin[i] = vals[i]
			}
		}
		for j := 3; j < n; j++ {
			cand := prefMin[j-3] + vals[j]
			if cand < ans {
				ans = cand
			}
		}
	}
	// Adjacent pair (distance 1): shoot between them
	for i := 0; i+1 < n; i++ {
		x, y := arr[i], arr[i+1]
		cand := max((x+y+2)/3, max((x+1)/2, (y+1)/2))
		if cand < ans {
			ans = cand
		}
	}
	// Distance-2 pair: shoot at the middle
	for i := 0; i+2 < n; i++ {
		x, z := arr[i], arr[i+2]
		cand := (x + z + 1) / 2
		if cand < ans {
			ans = cand
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(46))
	for t := 0; t < 100; t++ {
		n := rng.Intn(8) + 2
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(30) + 1
			sb.WriteString(fmt.Sprintf("%d ", arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
