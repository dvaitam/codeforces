package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expectedOutputs(ops []string) []string {
	var mxx, mxy int64
	res := []string{}
	for _, line := range ops {
		parts := strings.Fields(line)
		op := parts[0]
		x, _ := strconv.ParseInt(parts[1], 10, 64)
		y, _ := strconv.ParseInt(parts[2], 10, 64)
		if x > y {
			x, y = y, x
		}
		if op == "+" {
			if x > mxx {
				mxx = x
			}
			if y > mxy {
				mxy = y
			}
		} else {
			if mxx <= x && mxy <= y {
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	t, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	ops := lines[1 : 1+t]
	expected := expectedOutputs(ops)
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outLines := strings.Fields(string(out))
	idx := 0
	for _, line := range ops {
		if line[0] == '?' {
			if idx >= len(outLines) {
				fmt.Printf("missing output for query %d\n", idx+1)
				os.Exit(1)
			}
			if outLines[idx] != expected[idx] {
				fmt.Printf("query %d failed: expected %s got %s\n", idx+1, expected[idx], outLines[idx])
				os.Exit(1)
			}
			idx++
		}
	}
	if idx != len(outLines) {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
