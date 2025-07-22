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

type testCaseD struct {
	n   int
	arr []int
}

func isAP(seq []int) bool {
	if len(seq) <= 1 {
		return len(seq) > 0
	}
	d := seq[1] - seq[0]
	for i := 2; i < len(seq); i++ {
		if seq[i]-seq[i-1] != d {
			return false
		}
	}
	return true
}

func existsSolution(a []int) bool {
	n := len(a)
	if n < 2 {
		return false
	}
	m := 1 << n
	for mask := 1; mask < m-1; mask++ {
		var s1, s2 []int
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				s1 = append(s1, a[i])
			} else {
				s2 = append(s2, a[i])
			}
		}
		if isAP(s1) && isAP(s2) {
			return true
		}
	}
	return false
}

func checkOutput(a []int, out string, hasSol bool) error {
	out = strings.TrimSpace(out)
	if out == "No solution" {
		if hasSol {
			return fmt.Errorf("solution exists but output says none")
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		return fmt.Errorf("expected two lines of sequences")
	}
	parse := func(s string) ([]int, error) {
		if strings.TrimSpace(s) == "" {
			return nil, fmt.Errorf("empty sequence")
		}
		f := strings.Fields(s)
		res := make([]int, len(f))
		for i, x := range f {
			v, err := strconv.Atoi(x)
			if err != nil {
				return nil, err
			}
			res[i] = v
		}
		return res, nil
	}
	seq1, err := parse(lines[0])
	if err != nil {
		return fmt.Errorf("bad numbers in first line")
	}
	seq2, err := parse(lines[1])
	if err != nil {
		return fmt.Errorf("bad numbers in second line")
	}
	if !isAP(seq1) || !isAP(seq2) {
		return fmt.Errorf("sequences are not arithmetic")
	}
	// check order and usage
	i1, i2 := 0, 0
	for _, v := range a {
		if i1 < len(seq1) && seq1[i1] == v {
			i1++
			continue
		}
		if i2 < len(seq2) && seq2[i2] == v {
			i2++
			continue
		}
		return fmt.Errorf("elements do not match input order")
	}
	if i1 != len(seq1) || i2 != len(seq2) {
		return fmt.Errorf("unused elements in sequences")
	}
	return nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = v
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		expectedExist := existsSolution(arr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if err := checkOutput(arr, out, expectedExist); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
