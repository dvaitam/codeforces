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

func expected(n int, x int64, ops []struct {
	op string
	d  int64
}) (int64, int) {
	distressed := 0
	for _, o := range ops {
		if o.op == "+" {
			x += o.d
		} else {
			if x >= o.d {
				x -= o.d
			} else {
				distressed++
			}
		}
	}
	return x, distressed
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		x, _ := strconv.ParseInt(parts[1], 10, 64)
		if len(parts) != 2+n*2 {
			fmt.Printf("test %d: wrong number of operations\n", idx)
			os.Exit(1)
		}
		ops := make([]struct {
			op string
			d  int64
		}, n)
		for i := 0; i < n; i++ {
			ops[i].op = parts[2+i*2]
			d, err := strconv.ParseInt(parts[3+i*2], 10, 64)
			if err != nil {
				fmt.Printf("test %d: invalid number %q\n", idx, parts[3+i*2])
				os.Exit(1)
			}
			ops[i].d = d
		}
		// expected result
		exX, exDist := expected(n, x, ops)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for _, o := range ops {
			input.WriteString(o.op)
			input.WriteByte(' ')
			input.WriteString(strconv.FormatInt(o.d, 10))
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errb.String())
			os.Exit(1)
		}
		r := strings.NewReader(strings.TrimSpace(out.String()))
		var gotX int64
		var gotDist int
		if _, err := fmt.Fscan(r, &gotX, &gotDist); err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if gotX != exX || gotDist != exDist {
			fmt.Printf("test %d failed: expected %d %d got %d %d\n", idx, exX, exDist, gotX, gotDist)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
