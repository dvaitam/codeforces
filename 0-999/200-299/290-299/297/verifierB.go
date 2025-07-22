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

func aliceGreater(a, b []int) bool {
	diff := make(map[int]int)
	for _, x := range a {
		diff[x]++
	}
	for _, x := range b {
		diff[x]--
	}
	keys := make([]int, 0, len(diff))
	for k := range diff {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	suffix := 0
	for i := len(keys) - 1; i >= 0; i-- {
		suffix += diff[keys[i]]
		if suffix > 0 {
			return true
		}
	}
	return false
}

func randArray(rng *rand.Rand, n int, k int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(k) + 1
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		k := rng.Intn(10) + 1
		a := randArray(rng, n, k)
		b := randArray(rng, m, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := "NO"
		if aliceGreater(a, b) {
			expect = "YES"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		res := strings.ToUpper(strings.TrimSpace(out))
		if res != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
