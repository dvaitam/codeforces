package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expected(n, a, b int, groups []int) int {
	denied := 0
	half := 0
	for _, t := range groups {
		if t == 1 {
			if a > 0 {
				a--
			} else if b > 0 {
				b--
				half++
			} else if half > 0 {
				half--
			} else {
				denied++
			}
		} else { // t==2
			if b > 0 {
				b--
			} else {
				denied += 2
			}
		}
	}
	return denied
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	if exe == "--" && len(os.Args) == 3 {
		exe = os.Args[2]
	}
	const testcasesRaw = `100
7 3 0
2 2 2 2 2 2 1
9 1 2
1 1 2 1 2 1 1 2 2
9 0 2
2 2 1 2 2 2 1 1 1
7 5 5
1 2 2 1 2 1 1
10 1 1
1 2 1 1 2 2 1 2 2 1
9 2 4
1 2 2 1 2 2 1 2 1
4 1 0
2 2 1 1
3 1 0
1 2 2
9 1 1
2 2 2 2 2 1 2 1 2
10 5 2
1 1 1 2 1 1 2 1 2 2
1 0 1
1
1 4 5
1
1 0 5
1
10 4 0
2 1 2 1 1 1 1 1 1 2
4 5 0
1 2 1 2
2 1 0
2 2
7 1 0
2 1 1 2 1 2 2
8 4 1
1 1 1 1 2 2 1 2
3 0 3
2 2 2
7 5 2
1 1 2 1 2 1 2
3 1 3
2 2 2
10 5 4
1 2 2 2 1 1 1 2 1 1
4 5 3
2 2 1 2
10 3 5
1 1 2 1 2 1 2 2 1 1
8 2 2
2 1 2 1 1 1 1 2
7 2 0
1 1 1 1 1 1 1
5 2 5
1 1 2 2 1
1 2 3
1
5 1 5
2 1 1 2 1
1 0 1
2
9 2 2
1 2 2 2 2 1 1 2 2
1 1 1
2
6 2 2
1 2 1 1 2 1
3 4 2
2 2 1
5 0 3
1 1 2 1 1
5 3 2
2 2 1 1 2
8 2 2
1 2 1 2 2 1 2 2
3 1 5
2 1 1
2 1 5
1 1
7 0 0
2 2 2 2 1 2 1
6 1 2
1 2 1 2 1 1
1 4 3
1
2 3 3
2 1
1 1 4
1
2 1 3
2 2
9 1 0
2 1 2 2 2 2 2 2 1
9 4 1
1 2 2 2 1 1 2 1 2
10 2 5
2 1 1 1 1 1 1 1 2 1
8 2 1
1 2 2 2 1 1 1 2
4 2 4
2 2 2 1
1 5 0
1
5 4 4
2 2 1 2 2
5 3 3
2 1 1 1 1
5 5 2
1 1 2 2 1
8 4 5
1 1 2 2 1 2 2 2
1 0 3
1
1 0 4
1
6 0 5
2 1 2 2 1 1
10 5 3
2 1 2 1 2 2 1 1 1 2
8 2 1
2 2 1 1 1 2 2 1
9 5 4
1 1 1 2 2 1 1 2 2
2 0 3
1 1
7 1 5
1 2 2 2 1 2 2
5 0 2
1 1 2 1 2
6 1 0
1 2 1 1 1 1
4 5 5
1 1 2 2
1 4 1
1
3 3 0
2 2 2
3 0 1
2 2 2
5 2 4
2 1 1 1 2
3 3 1
1 1 2
9 1 1
1 2 2 2 1 1 1 2 1
2 1 3
2 2
3 4 3
2 2 1
4 3 5
2 1 2 2
1 3 5
2
8 1 2
2 2 1 1 1 2 1 1
8 3 1
2 2 1 1 2 2 1 2
8 5 5
1 2 1 2 1 2 2 1
4 2 0
2 1 2 2
2 5 3
1 2
5 5 0
1 2 1 1 2
7 4 4
2 2 1 2 2 2 1
10 0 0
1 2 2 2 1 1 1 1 2 2
5 4 1
1 1 2 2 2
5 3 2
1 1 1 1 2
7 4 3
1 2 2 2 2 1 2
8 5 4
2 2 1 2 2 1 2 1
10 1 0
2 2 2 1 1 1 2 1 2 1
2 4 0
2 2
4 1 4
2 1 1 2
8 5 2
2 1 2 2 2 1 2 1
9 2 1
1 1 2 2 2 2 2 2 1
4 3 1
1 2 1 1
2 1 2
1 1
4 4 2
1 2 2 2
8 0 5
2 2 2 2 2 1 2 2
1 0 0
1
6 0 2
1 1 1 2 1 2
9 1 3
1 2 1 1 2 2 2 2 2`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		a, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		b, _ := strconv.Atoi(scan.Text())
		groups := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			groups[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(groups[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := fmt.Sprintf("%d\n", expected(n, a, b, groups))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
