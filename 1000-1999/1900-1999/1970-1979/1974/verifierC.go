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

type key struct{ x, y int }

func expected(arr []int) int64 {
	n := len(arr)
	if n < 3 {
		return 0
	}
	triples := make([][3]int, n-2)
	for i := 0; i < n-2; i++ {
		triples[i] = [3]int{arr[i], arr[i+1], arr[i+2]}
	}
	var ans int64
	ab := make(map[key]int)
	bc := make(map[key]int)
	ac := make(map[key]int)
	abc := make(map[[3]int]int)
	for _, t := range triples {
		k1 := key{t[0], t[1]}
		k2 := key{t[1], t[2]}
		k3 := key{t[0], t[2]}
		ans += int64(ab[k1] - abc[t])
		ans += int64(bc[k2] - abc[t])
		ans += int64(ac[k3] - abc[t])
		ab[k1]++
		bc[k2]++
		ac[k3]++
		abc[t]++
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
	var out bytes.Buffer
	var stderr bytes.Buffer
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
		n := rng.Intn(8) + 3
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(6)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedVal := fmt.Sprintf("%d", expected(arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expectedVal, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
