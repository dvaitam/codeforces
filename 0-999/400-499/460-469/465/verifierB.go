package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestCase struct {
	n   int
	arr []int
	ans int
}

func compute(arr []int) int {
	total, segments := 0, 0
	for i, v := range arr {
		if v == 1 {
			total++
			if i == 0 || arr[i-1] == 0 {
				segments++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return total + segments - 1
}

func genCases() []TestCase {
	var cases []TestCase
	for n := 1; n <= 7; n++ {
		for mask := 0; mask < (1 << n); mask++ {
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					arr[i] = 1
				}
			}
			cases = append(cases, TestCase{n, arr, compute(arr)})
		}
	}
	return cases
}

func runCase(bin string, tc TestCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

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
	output := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(output)
	if err != nil {
		return fmt.Errorf("invalid output: %s", output)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
