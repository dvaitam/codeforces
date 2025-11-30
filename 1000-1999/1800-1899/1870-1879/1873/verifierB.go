package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
3
9 1 4
2
7 7
8
6 3 1 7 0 6 6 9
1
7
5
3 9 1 5 0
1
0
9
0 6 3 6 0 8 3 7 7
9
3 5 3 3 7 4 0 6 8
2
2 4
2
5 8
7
8 3 4 4 9 7 8
7
9 0 7 3 6 6 2
6
8 5 1 7 8 1
3
8 6 5
8
0 7 0 4 9 9 9 6
3
2 8 3
1
3
9
8 3 6 8 5 9 5 7 4
9
9 0 6 8 2 8 8 3 6
1
7
6
9 8 3 8 6 7
6
6 5 0 8 8 9
6
7 9 0 3 2 8
3
1 8 4
1
1
2
0 7
1
4
4
4 1 9 2
6
4 1 2 2 4 8
3
4 4 7
6
7 7 1 0 4 6
6
6 3 4 1 4 8
4
9 6 0 3
1
6
3
0 2 7
9
6 8 3 8 7 3 8 0 6
6
6 0 4 2 3 0
5
1 1 4 4 2
7
9 4 2 0 8 0 9
4
9 7 2 9
9
0 6 3 5 1 3 9 6 9
4
7 1 6 4
9
7 0 5 9 6 4 0 2 3
6
9 2 5 6 3 4
2
6 8
6
8 7 8 3 1 0
2
2 2
3
8 3 4
6
9 8 4 5 5 5
2
4 3
8
2 9 8 1 5 0 6 1
7
2 2 5 1 9 9 6
2
9 8
4
9 1 4 5
5
9 8 1 7 4
2
0 4
1
9
1
1
7
1 0 3 3 9 6 2
2
7 2
4
2 1 6 6
9
4 8 4 7 5 1 3 5 0
1
0
5
9 5 7 6 5
7
1 1 5 9 7 1 4
4
9 8 7 5
5
2 8 3 4 3
4
5 1 4 1
8
1 9 5 3 6 4 0 5
3
5 9 4
4
5 1 8 9
2
3 3
1
3
7
1 4 8 1 1 0 0
5
5 7 7 2 1
9
5 1 8 2 2 2 2 5 4
2
8 9
5
2 3 2 8 0
6
9 8 3 2 4 6
9
2 0 3 4 1 7 6 8 4
9
7 8 7 0 6 5 2 4 7
1
6
1
0
6
9 2 9 2 2 4
5
6 9 6 2 9
2
3 7
1
2
9
5 8 7 3 3 5 7 7 3
7
5 8 9 4 3 0 1
9
5 2 8 3 4 4 4 8 5
3
7 9 1
2
9 8
7
2 2 4 6 3 9 0
8
6 5 6 8 2 8 0 8
2
4 1
5
1 2 9 1 7
4
6 6 6 2
6
7 2 9 7 3 1
7
9 8 6 1 4 4 3
7
8 0 3 8 7 9 0`

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func expected(nums []int) string {
	n := len(nums)
	var maxProd int64
	for i := 0; i < n; i++ {
		prod := int64(1)
		for j := 0; j < n; j++ {
			val := nums[j]
			if i == j {
				val++
			}
			prod *= int64(val)
		}
		if prod > maxProd {
			maxProd = prod
		}
	}
	return fmt.Sprintf("%d", maxProd)
}

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "no embedded testcases")
		os.Exit(1)
	}
	pos := 0
	t, err := strconv.Atoi(tokens[pos])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count header\n")
		os.Exit(1)
	}
	pos++
	var inputs []string
	var exps []string
	for caseNum := 1; caseNum <= t; caseNum++ {
		if pos >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing n\n", caseNum)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on case %d\n", caseNum)
			os.Exit(1)
		}
		pos++
		if pos+n > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing numbers\n", caseNum)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(tokens[pos+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value on case %d\n", caseNum)
				os.Exit(1)
			}
			nums[i] = val
		}
		pos += n
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inputs = append(inputs, sb.String())
		exps = append(exps, expected(nums))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := run(exe, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exps[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
