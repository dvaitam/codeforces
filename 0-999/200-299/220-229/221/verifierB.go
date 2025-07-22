package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func shares(d int, has []bool) bool {
	if d == 0 {
		return has[0]
	}
	for d > 0 {
		if has[d%10] {
			return true
		}
		d /= 10
	}
	return false
}

func solve(x int) int {
	s := fmt.Sprint(x)
	has := make([]bool, 10)
	for _, ch := range s {
		has[ch-'0'] = true
	}
	ans := 0
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			d1 := i
			d2 := x / i
			if shares(d1, has) {
				ans++
			}
			if d2 != d1 && shares(d2, has) {
				ans++
			}
		}
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
	var bin string
	switch {
	case len(os.Args) == 2:
		bin = os.Args[1]
	case len(os.Args) == 3 && os.Args[1] == "--":
		bin = os.Args[2]
	default:
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var x int
		fmt.Sscan(line, &x)
		input := fmt.Sprintf("%d\n", x)
		want := fmt.Sprintf("%d", solve(x))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
