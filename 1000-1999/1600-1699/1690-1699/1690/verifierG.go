package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func computeTrains(a []int) int {
	n := len(a)
	speeds := make([]int, n)
	speeds[0] = a[0]
	trains := 1
	for i := 1; i < n; i++ {
		if a[i] < speeds[i-1] {
			speeds[i] = a[i]
			trains++
		} else {
			speeds[i] = speeds[i-1]
		}
	}
	return trains
}

func processCase(n, m int, a []int, ops [][2]int) []int {
	res := make([]int, m)
	cur := append([]int(nil), a...)
	for i, op := range ops {
		k := op[0] - 1
		d := op[1]
		cur[k] -= d
		res[i] = computeTrains(cur)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(48))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(10) + 1
		}
		ops := make([][2]int, m)
		for i := 0; i < m; i++ {
			k := rng.Intn(n) + 1
			d := rng.Intn(a[k-1] + 1)
			a[k-1] -= d
			ops[i] = [2]int{k, d}
		}
		// rebuild original a for expected computation
		orig := make([]int, n)
		for i := range orig {
			orig[i] = a[i]
		}
		for i := 0; i < m; i++ {
			orig[ops[i][0]-1] += ops[i][1]
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", orig[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d %d\n", ops[i][0], ops[i][1])
		}
		input := sb.String()
		expected := processCase(n, m, orig, ops)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < m {
			fmt.Fprintf(os.Stderr, "case %d wrong output size: got %q\n", tc+1, out)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			var v int
			fmt.Sscanf(fields[i], "%d", &v)
			if v != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at op %d: expected %d got %d\ninput:\n%s", tc+1, i+1, expected[i], v, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
