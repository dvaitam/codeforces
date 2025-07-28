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

func generateCase(rng *rand.Rand) (string, []int, bool) {
	n := rng.Intn(100) + 1
	arr := make([]int, n)
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
		counts[arr[i]]++
	}
	dup := 0
	for _, c := range counts {
		if c > 1 {
			dup++
		}
	}
	expect := dup >= 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), arr, expect
}

func checkOutput(arr []int, out string, expect bool) error {
	out = strings.TrimSpace(out)
	if !expect {
		if out != "-1" {
			return fmt.Errorf("expected -1 got %q", out)
		}
		return nil
	}
	fields := strings.Fields(out)
	if len(fields) != len(arr) {
		return fmt.Errorf("expected %d numbers got %d", len(arr), len(fields))
	}
	b := make([]int, len(arr))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		if v < 1 || v > 3 {
			return fmt.Errorf("number out of range: %d", v)
		}
		b[i] = v
	}
	cond1, cond2, cond3 := false, false, false
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] != arr[j] {
				continue
			}
			if b[i] == 1 && b[j] == 2 || b[i] == 2 && b[j] == 1 {
				cond1 = true
			}
			if b[i] == 1 && b[j] == 3 || b[i] == 3 && b[j] == 1 {
				cond2 = true
			}
			if b[i] == 2 && b[j] == 3 || b[i] == 3 && b[j] == 2 {
				cond3 = true
			}
		}
	}
	cnt := 0
	if cond1 {
		cnt++
	}
	if cond2 {
		cnt++
	}
	if cond3 {
		cnt++
	}
	if cnt != 2 {
		return fmt.Errorf("exactly two conditions should hold, got %d", cnt)
	}
	return nil
}

func runCase(bin, input string, arr []int, expect bool) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutput(arr, out.String(), expect)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, arr, expect := generateCase(rng)
		if err := runCase(bin, in, arr, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
