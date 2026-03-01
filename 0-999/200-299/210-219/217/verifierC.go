package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesC.txt.
const embeddedTestcasesC = `2 1
3 (((?^1)&((1^?)^(1^0)))^(((?&0)&(0|0))&0))
3 ((((?^1)&(1&1))^1)&(((?&?)|(?&1))&((?^?)&(0^1))))
3 (1^(((?^?)^(0^1))^((1|?)^0)))
2 ?
2 0
2 ((((0^0)|(0|?))|0)|(0|((?|0)|(1|?))))
2 (0|?)
2 1
2 ((1^((0&0)^(0&0)))^(((0&1)&?)|(?^(0^1))))
2 0
2 (((1|0)^((?^1)^1))|0)
3 1
3 ?
2 0
2 0
2 ?
5 ((((1&?)|?)|1)^(0^((0^?)&(1^1))))
2 (1&?)
2 ((((1|0)^(0^0))&(1|(0&1)))|(0|1))
5 (?^((1&1)|?))
5 ((((1|?)&1)&((?^?)^1))^?)
3 (((1^(?^?))^((0|?)|(0&1)))|(((?&0)|(?&?))^1))
2 ((((1&?)|(?&?))|((?^0)^(0^?)))&(((0^1)&(1^?))&0))
2 ((((1|?)|(1|?))^(?^?))|(?&?))
3 ((((1^?)^(1&?))&((0&1)|?))|(0^((0|1)&(1|0))))
5 ((((?^1)|1)&((?&?)^(?|?)))|?)
3 (0^1)
3 ((((1|1)&1)|((1&0)&(?&0)))|(0&((1|0)^(1|?))))
5 ((((1^?)^(?^?))^(?^(1|1)))^(1|((1|1)|0)))
2 (?&(1&?))
5 ((((0|0)|(1|1))|1)&?)
3 (((0^0)&(?&(1|?)))|(((1&1)|(?&?))^((1^1)^(0&1))))
5 (((?|(0|1))&(0^(0|0)))|?)
2 ((1&0)^?)
2 ?
2 ((((0|?)|0)^((1^1)|(?&?)))^(0^((0|1)|1)))
5 0
5 ((1^1)&(((0^0)|(0|0))|((0^?)|?)))
5 ((((0^0)^?)|((?|1)|(0|1)))^(((1&?)^(1&?))&((0&0)^(?&?))))
5 ((((1&1)|(1^?))|1)^(((1|?)&(0&?))&((1^0)^(0&0))))
2 1
2 (((1&(1|?))^(?|(0^?)))&(((0|?)&0)&?))
4 ((0&((?|?)&(0&1)))^((0&?)&((?^0)|?)))
5 (1^(((1^?)&(0&0))^((1|0)|(0^?))))
4 ((?^?)|1)
4 ((((?|1)|(1^1))&((0|0)^(0^1)))^((?^(1&1))^0))
3 (((0^(0|?))^((0&?)|(0|?)))|0)
5 1
3 ?
5 (0|(?&1))
5 ((((0|1)|(?|1))&((1&0)&(0&1)))^(((1&0)^(1^?))|1))
3 (((?^1)^((?^0)&(?^1)))&?)
5 (?&(((?&1)^0)&(1^(0^1))))
4 ((((0^1)&(1|0))&1)^(((0&?)|(?&0))&((0&1)|(?^0))))
2 ((((0^0)&(0&1))^((1^1)&(?|0)))&(((0|?)&(?^0))|1))
3 ((((?&1)|(0^?))^((0^?)|0))&(0|(0^?)))
2 0
4 1
5 ((1|(0^(0^?)))|(0&(1^(?|1))))
5 ((((1^?)|(?^0))&((0|1)|0))|(((0&1)^(1|?))|((?|0)|(0|?))))
5 ((0^((0|0)|(?&0)))^((0^(0^1))^((1|?)|(1&1))))
4 ((((0&1)|(0^?))^?)^(?|(?|(?|1))))
3 ((((1|0)|(0^?))^((1&0)&0))|(((?|1)|(0&0))^(0^0)))
3 ((1|(?^(1&?)))&(((0&?)&(?&1))&((0^?)^1)))
5 ((((0&?)|0)^((1^0)^(1&0)))^((?^?)|1))
4 0
5 ((((1^0)&?)|((0|1)^1))|((1^0)^1))
5 ?
4 1
2 (?|((0&(1|0))|(0^(?^?))))
5 (((0&0)^(0^0))^?)
2 (?&((0&(?^1))^((0|?)|0)))
3 ((0|(1^(0^?)))|0)
2 (0^(?|(1&(1|0))))
5 ((?&((0^?)^?))|?)
3 0
3 (((0^1)&(1|(0&?)))|(0|1))
3 0
2 1
2 (((0^(0|0))^((1^0)&(?|1)))^?)
2 0
2 ((((0|?)^(0|?))&0)^0)
2 ((?|1)&(((0&1)|(?&0))^((?^?)&(1^?))))
2 ((1&((0^?)&0))^(((?&?)|?)&?))
2 0
2 (0^1)
2 ((1|?)^1)
2 ((((0|1)&(?^?))^((1|1)|(0^?)))|?)
3 ((((1^?)|(?&?))|(0|1))|(((0|1)|?)|(0^(0&1))))
3 (1|(((1|?)&0)^((1|1)^(?&1))))
5 (0|((0|(1^1))^(1&(0&?))))
3 1
4 1
5 (0&(1^((1^1)^(1&0))))
5 ((0^?)&0)
4 ((((?&1)|(1&0))&1)|0)
3 1
3 (((0|(1|1))&0)^(((0&0)&(?|?))|((1^?)&(1&0))))
5 ((?&1)|(((0|0)^(1|1))^((0^0)|(1^1))))`

func expected(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ""
	}
	formula := parts[1]
	// Track which unary boolean functions can be represented by each sub-expression.
	// Function ids use two-bit truth tables:
	// 0 -> 00 (const 0), 1 -> 01 (x), 2 -> 10 (!x), 3 -> 11 (const 1).
	// A bitmask keeps a set of possible ids.
	op := make([]byte, 0, len(formula))
	val := make([]int, 0, len(formula))
	apply := func(x, y int, oper byte) int {
		res := 0
		for i := 0; i < 4; i++ {
			if x>>i&1 == 0 {
				continue
			}
			for j := 0; j < 4; j++ {
				if y>>j&1 == 0 {
					continue
				}
				k := 0
				switch oper {
				case '^':
					k = i ^ j
				case '&':
					k = i & j
				case '|':
					k = i | j
				}
				res |= 1 << k
			}
		}
		return res
	}

	for i := 0; i < len(formula); i++ {
		c := formula[i]
		switch c {
		case '0':
			val = append(val, 1<<0)
		case '1':
			val = append(val, 1<<3)
		case '?':
			// A leaf can be any colony, so relative to two distinct colonies it may
			// behave as x or !x.
			val = append(val, (1<<1)|(1<<2))
		case '(', '|', '&', '^':
			op = append(op, c)
		case ')':
			l := len(val)
			v2 := val[l-1]
			val = val[:l-1]
			l = len(val)
			v1 := val[l-1]
			val = val[:l-1]
			oper := op[len(op)-1]
			op = op[:len(op)-1]
			if len(op) > 0 && op[len(op)-1] == '(' {
				op = op[:len(op)-1]
			}
			val = append(val, apply(v1, v2, oper))
		}
	}
	if len(val) > 0 && val[len(val)-1]&((1<<1)|(1<<2)) != 0 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesC))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expect := expected(line)
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
