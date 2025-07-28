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

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1845D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	n   int
	arr []int64
}

func readCases() ([]Case, error) {
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	var cases []Case
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+n {
			return nil, fmt.Errorf("bad line %q", line)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+i], 10, 64)
			arr[i] = v
		}
		cases = append(cases, Case{n, arr})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func finalRating(arr []int64, k int64) int64 {
	rating := int64(0)
	for _, v := range arr {
		rating += v
		if rating < k {
			rating = k
		}
	}
	return rating
}

func runCase(bin, ref string, c Case) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for i, v := range c.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expStr, err := run(ref, input)
	if err != nil {
		return fmt.Errorf("reference: %v", err)
	}
	expK, err := strconv.ParseInt(strings.TrimSpace(expStr), 10, 64)
	if err != nil {
		return fmt.Errorf("ref output invalid: %v", err)
	}
	refRating := finalRating(c.arr, expK)

	gotStr, err := run(bin, input)
	if err != nil {
		return err
	}
	gotK, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	candRating := finalRating(c.arr, gotK)
	if candRating != refRating {
		return fmt.Errorf("final rating %d expected %d", candRating, refRating)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	cases, err := readCases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
