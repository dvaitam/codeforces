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

type op struct {
	cmd string
	id  int
}

func solve(ops []op) []int {
	pos := make(map[int]int)
	l, r := 0, 0
	size := 0
	res := []int{}
	for _, o := range ops {
		switch o.cmd {
		case "L":
			if size == 0 {
				pos[o.id] = 0
				l = 0
				r = 0
				size = 1
			} else {
				l--
				pos[o.id] = l
				size++
			}
		case "R":
			if size == 0 {
				pos[o.id] = 0
				l = 0
				r = 0
				size = 1
			} else {
				r++
				pos[o.id] = r
				size++
			}
		case "?":
			p := pos[o.id]
			left := p - l
			right := r - p
			if left < right {
				res = append(res, left)
			} else {
				res = append(res, right)
			}
		}
	}
	return res
}

func runCase(bin string, ops []op) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, o := range ops {
		input.WriteString(fmt.Sprintf("%s %d\n", o.cmd, o.id))
	}
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	expectVals := solve(ops)
	if len(gotLines) != len(expectVals) {
		return fmt.Errorf("expected %d lines got %d", len(expectVals), len(gotLines))
	}
	for i, v := range expectVals {
		if gotLines[i] != strconv.Itoa(v) {
			return fmt.Errorf("mismatch at %d: expected %d got %s", i+1, v, gotLines[i])
		}
	}
	return nil
}

func parseOps(parts []string) ([]op, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	q, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	if len(parts) != 1+2*q {
		return nil, fmt.Errorf("need %d tokens", 1+2*q)
	}
	ops := make([]op, q)
	idx := 1
	for i := 0; i < q; i++ {
		cmd := parts[idx]
		id, _ := strconv.Atoi(parts[idx+1])
		idx += 2
		ops[i] = op{cmd, id}
	}
	return ops, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		ops, err := parseOps(parts)
		if err != nil {
			fmt.Printf("invalid test line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := runCase(bin, ops); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
