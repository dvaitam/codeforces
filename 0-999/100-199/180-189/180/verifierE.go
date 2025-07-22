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

func expected(n, m, k int, colors []int) int {
	pos := make(map[int][]int)
	for i, c := range colors {
		pos[c] = append(pos[c], i+1)
	}
	best := 1
	for _, p := range pos {
		if len(p) <= best {
			continue
		}
		left := 0
		for right := 0; right < len(p); right++ {
			for left <= right && (p[right]-p[left]+1-(right-left+1) > k) {
				left++
			}
			if cur := right - left + 1; cur > best {
				best = cur
			}
		}
	}
	return best
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

func runCase(bin string, n, m, k int, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	exp := expected(n, m, k, arr)
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	val, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct {
		n, m, k int
		arr     []int
	}{{1, 1, 0, []int{1}}, {5, 2, 1, []int{1, 2, 1, 2, 1}}}

	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(5) + 1
		k := rng.Intn(n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(m) + 1
		}
		cases = append(cases, struct {
			n, m, k int
			arr     []int
		}{n, m, k, arr})
	}

	for idx, c := range cases {
		if err := runCase(bin, c.n, c.m, c.k, c.arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
