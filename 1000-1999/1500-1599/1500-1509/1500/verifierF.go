package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1500F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, bytes.Buffer, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errBuf, err
}

func checkHeights(h []int64, w []int64) bool {
	n := len(h)
	for i := 0; i < n-2; i++ {
		mx := h[i]
		mn := h[i]
		if h[i+1] > mx {
			mx = h[i+1]
		}
		if h[i+2] > mx {
			mx = h[i+2]
		}
		if h[i+1] < mn {
			mn = h[i+1]
		}
		if h[i+2] < mn {
			mn = h[i+2]
		}
		if mx-mn != w[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	idx := 0
	for scan.Scan() {
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		C, _ := strconv.Atoi(scan.Text())
		w := make([]int64, n-2)
		for i := 0; i < n-2; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			w[i] = v
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, C)
		if n > 2 {
			for i := 0; i < n-2; i++ {
				if i > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprintf(&input, "%d", w[i])
			}
			input.WriteByte('\n')
		} else {
			input.WriteByte('\n')
		}
		idx++
		expect, _, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
			os.Exit(1)
		}
		got, stderr, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, stderr.String())
			os.Exit(1)
		}
		exp := strings.ToLower(strings.Fields(expect)[0])
		fields := strings.Fields(got)
		if len(fields) == 0 {
			fmt.Printf("case %d empty output\n", idx)
			os.Exit(1)
		}
		ans := strings.ToLower(fields[0])
		if ans != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, exp, ans)
			os.Exit(1)
		}
		if ans == "yes" {
			if len(fields) != n+1 {
				fmt.Printf("case %d: expected %d numbers\n", idx, n)
				os.Exit(1)
			}
			h := make([]int64, n)
			for i := 0; i < n; i++ {
				v, err := strconv.ParseInt(fields[i+1], 10, 64)
				if err != nil {
					fmt.Printf("case %d invalid number\n", idx)
					os.Exit(1)
				}
				if v < 0 {
					fmt.Printf("case %d negative height\n", idx)
					os.Exit(1)
				}
				h[i] = v
			}
			if !checkHeights(h, w) {
				fmt.Printf("case %d heights do not match constraints\n", idx)
				os.Exit(1)
			}
		}
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
