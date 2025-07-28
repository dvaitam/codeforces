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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

func solveCase(a []int) int {
	n := len(a)
	used := make([]bool, n+1)
	extras := make([]int, 0)
	for _, v := range a {
		if v >= 1 && v <= n && !used[v] {
			used[v] = true
		} else {
			extras = append(extras, v)
		}
	}
	sort.Ints(extras)
	missing := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !used[i] {
			missing = append(missing, i)
		}
	}
	if len(extras) != len(missing) {
		return -1
	}
	for i, m := range missing {
		if extras[i] <= 2*m {
			return -1
		}
	}
	return len(extras)
}

func genTest(rng *rand.Rand) []int {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(30) + 1
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	var input bytes.Buffer
	input.WriteString(fmt.Sprintf("%d\n", tests))
	expected := make([]int, tests)
	for i := 0; i < tests; i++ {
		arr := genTest(rng)
		expected[i] = solveCase(append([]int(nil), arr...))
		input.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')
	}
	out, err := run(bin, input.String())
	if err != nil {
		fmt.Printf("runtime error: %v\n%s", err, out)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != tests {
		fmt.Printf("expected %d lines of output got %d\n", tests, len(lines))
		os.Exit(1)
	}
	for i := 0; i < tests; i++ {
		got := strings.TrimSpace(lines[i])
		if got != fmt.Sprintf("%d", expected[i]) {
			fmt.Printf("test %d failed expected:%d got:%s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
