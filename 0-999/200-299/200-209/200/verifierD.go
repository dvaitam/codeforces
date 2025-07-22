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

type Proc struct {
	name  string
	types []string
}

type Call struct {
	name string
	args []string
}

func parseProc(line string) Proc {
	line = strings.TrimSpace(line)
	l := strings.Index(line, "(")
	r := strings.LastIndex(line, ")")
	header := strings.Fields(line[:l])
	pname := header[len(header)-1]
	args := strings.Split(line[l+1:r], ",")
	types := make([]string, 0, len(args))
	for _, a := range args {
		types = append(types, strings.TrimSpace(a))
	}
	return Proc{name: pname, types: types}
}

func parseVar(line string) (string, string) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 2 {
		return "", ""
	}
	return fields[1], fields[0]
}

func parseCall(line string) Call {
	line = strings.TrimSpace(line)
	l := strings.Index(line, "(")
	r := strings.LastIndex(line, ")")
	name := strings.TrimSpace(line[:l])
	argsStr := strings.TrimSpace(line[l+1 : r])
	args := []string{}
	if len(argsStr) > 0 {
		parts := strings.Split(argsStr, ",")
		for _, p := range parts {
			args = append(args, strings.TrimSpace(p))
		}
	}
	return Call{name, args}
}

func solveCase(procs []Proc, vars map[string]string, calls []Call) []int {
	res := make([]int, len(calls))
	for i, c := range calls {
		count := 0
		for _, p := range procs {
			if p.name != c.name || len(p.types) != len(c.args) {
				continue
			}
			ok := true
			for j, typ := range p.types {
				if typ != "T" && typ != vars[c.args[j]] {
					ok = false
					break
				}
			}
			if ok {
				count++
			}
		}
		res[i] = count
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	lineScan := bufio.NewScanner(bytes.NewReader(data))
	if !lineScan.Scan() {
		fmt.Println("bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(lineScan.Text()))
	var expected []int
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		lineScan.Scan()
		n, _ := strconv.Atoi(strings.TrimSpace(lineScan.Text()))
		procs := make([]Proc, n)
		for i := 0; i < n; i++ {
			lineScan.Scan()
			procs[i] = parseProc(lineScan.Text())
		}
		lineScan.Scan()
		m, _ := strconv.Atoi(strings.TrimSpace(lineScan.Text()))
		vars := make(map[string]string, m)
		varNames := make([]string, m)
		for i := 0; i < m; i++ {
			lineScan.Scan()
			name, typ := parseVar(lineScan.Text())
			vars[name] = typ
			varNames[i] = name
		}
		lineScan.Scan()
		k, _ := strconv.Atoi(strings.TrimSpace(lineScan.Text()))
		calls := make([]Call, k)
		for i := 0; i < k; i++ {
			lineScan.Scan()
			calls[i] = parseCall(lineScan.Text())
		}
		res := solveCase(procs, vars, calls)
		expected = append(expected, res...)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < len(expected); i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for line %d\n", i+1)
			os.Exit(1)
		}
		v, err := strconv.Atoi(outScan.Text())
		if err != nil || v != expected[i] {
			fmt.Printf("line %d failed: expected %d got %s\n", i+1, expected[i], outScan.Text())
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
