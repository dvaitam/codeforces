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

func parseTestcases(path string) ([]int, [][][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ns := []int{}
	mats := [][][]byte{}
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		mat := make([][]byte, n)
		for i := 0; i < n; i++ {
			scanner.Scan()
			row := strings.TrimSpace(scanner.Text())
			mat[i] = []byte(row)
		}
		ns = append(ns, n)
		mats = append(mats, mat)
		scanner.Scan() // blank line
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return ns, mats, nil
}

func reachability(mat [][]byte) [][]bool {
	n := len(mat)
	reach := make([][]bool, n)
	for i := range reach {
		reach[i] = make([]bool, n)
	}
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if mat[i][j] == '1' {
				reach[j][i] = true
			}
		}
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if reach[i][k] {
				for j := 0; j < n; j++ {
					if reach[k][j] {
						reach[i][j] = true
					}
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		reach[i][i] = true
	}
	return reach
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ns, mats, err := parseTestcases("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, n := range ns {
		mat := mats[idx]
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(string(mat[i]))
			sb.WriteByte('\n')
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(strings.NewReader(out))
		scanner.Split(bufio.ScanWords)
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: missing colors\n", idx+1)
				os.Exit(1)
			}
			v := scanner.Text()
			if v != "0" && v != "1" {
				fmt.Fprintf(os.Stderr, "case %d: invalid color %s\n", idx+1, v)
				os.Exit(1)
			}
			if v == "1" {
				colors[i] = 1
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		reach := reachability(mat)
		for j := 0; j < n; j++ {
			for i := 0; i < j; i++ {
				if colors[i] == colors[j] {
					if !reach[j][i] {
						fmt.Fprintf(os.Stderr, "case %d: coloring invalid\n", idx+1)
						os.Exit(1)
					}
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(ns))
}
