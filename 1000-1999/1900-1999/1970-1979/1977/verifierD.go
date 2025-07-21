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

func best(matrix [][]byte) (int, string) {
	n := len(matrix)
	m := len(matrix[0])
	best := -1
	bestMask := 0
	for mask := 0; mask < 1<<n; mask++ {
		count := 0
		for j := 0; j < m; j++ {
			ones := 0
			for i := 0; i < n; i++ {
				bit := matrix[i][j]
				if mask>>i&1 == 1 {
					if bit == '1' {
						bit = '0'
					} else {
						bit = '1'
					}
				}
				if bit == '1' {
					ones++
				}
			}
			if ones == 1 {
				count++
			}
		}
		if count > best {
			best = count
			bestMask = mask
		}
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		if bestMask>>i&1 == 1 {
			res[i] = '1'
		} else {
			res[i] = '0'
		}
	}
	return best, string(res)
}

func apply(matrix [][]byte, maskStr string) int {
	n := len(matrix)
	m := len(matrix[0])
	count := 0
	for j := 0; j < m; j++ {
		ones := 0
		for i := 0; i < n; i++ {
			bit := matrix[i][j]
			if maskStr[i] == '1' {
				if bit == '1' {
					bit = '0'
				} else {
					bit = '1'
				}
			}
			if bit == '1' {
				ones++
			}
		}
		if ones == 1 {
			count++
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m int
		fmt.Sscan(line, &n, &m)
		matrix := make([][]byte, n)
		for i := 0; i < n; i++ {
			scanner.Scan()
			row := scanner.Text()
			matrix[i] = []byte(strings.TrimSpace(row))
		}
		scanner.Scan() // consume blank line

		bestCount, _ := best(matrix)
		// build input for candidate
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(string(matrix[i]))
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		scanOut := bufio.NewScanner(strings.NewReader(got))
		scanOut.Split(bufio.ScanWords)
		if !scanOut.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing answer\n", idx)
			os.Exit(1)
		}
		ansReported, err := strconv.Atoi(scanOut.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid answer\n", idx)
			os.Exit(1)
		}
		if !scanOut.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing mask\n", idx)
			os.Exit(1)
		}
		mask := scanOut.Text()
		if len(mask) != n {
			fmt.Fprintf(os.Stderr, "case %d: mask length mismatch\n", idx)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if mask[i] != '0' && mask[i] != '1' {
				fmt.Fprintf(os.Stderr, "case %d: invalid mask character\n", idx)
				os.Exit(1)
			}
		}
		if scanOut.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx)
			os.Exit(1)
		}
		gotCount := apply(matrix, mask)
		if gotCount != ansReported {
			fmt.Fprintf(os.Stderr, "case %d: reported %d but got %d\n", idx, ansReported, gotCount)
			os.Exit(1)
		}
		if gotCount != bestCount {
			fmt.Fprintf(os.Stderr, "case %d: answer not optimal\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
