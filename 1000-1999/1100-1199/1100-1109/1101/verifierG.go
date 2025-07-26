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

type caseG struct {
	n   int
	arr []int
}

func readCasesG(path string) ([]caseG, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]caseG, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		cases[i] = caseG{n, arr}
	}
	return cases, nil
}

func solveCaseG(cs caseG) int {
	a := append([]int{}, cs.arr...)
	x := 0
	for _, v := range a {
		x ^= v
	}
	if x == 0 {
		return -1
	}
	ans := 0
	for i := 29; i >= 0; i-- {
		id := -1
		mask := 1 << uint(i)
		for j, v := range a {
			if v&mask != 0 {
				id = j
				break
			}
		}
		if id == -1 {
			continue
		}
		ans++
		last := len(a) - 1
		a[id], a[last] = a[last], a[id]
		pivot := a[last]
		for j := 0; j < last; j++ {
			if a[j]&mask != 0 {
				a[j] ^= pivot
			}
		}
		a = a[:last]
	}
	return ans
}

func run(bin, stringInput string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(stringInput)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	cases, err := readCasesG("testcasesG.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, cs := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCaseG(cs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
