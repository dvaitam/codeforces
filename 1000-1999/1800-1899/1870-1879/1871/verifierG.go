package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "ODgwMgotNzM2MAo1ODkzCi0xNDI4Ci04Nzk0Ci05OTg2Ci01MjMwCjkyMTcKNTQwOQoyMjI3CjQ2NQotOTI4MwotMTA2NQo2MDE4Ci0zNTA5CjM1NjYKNzYzOQo3NjcwCi02OTEzCi0zNjc3Cjg0NjIKODE0MAotMTMwMAo5OTkwCi03MTEzCjM5MTYKOTk2Ci02OTQ5CjE4NjYKMzQzMwotMTc5Mwo0NTg4Ci02OTI4Ci0zNTQwCi00NTkKLTY4MDQKLTg1MDQKOTI4NgotMzQxOAoxODA4CjU5NDgKLTM2NTMKNjg2NAo4ODcwCjY1MDQKLTkwNjkKMTg0MQotMTk3Ngo5NzE4CjQwODAKLTI4CjE3MTAKOTMwNQotNjA4MgotNzA2Mgo2Mzg5CjcyMTEKLTM0NzYKLTYxNzIKOTkyOQotMTIyMgoyMzIKLTM1NzEKMjQ0NQo1ODUxCi0yNzAyCi01NDc1Cjk1NzMKLTMxMjMKNzEwOQotOTU5MAotMzgxNgotNDQ1OAotOTQ0Mgo5MDgKODI3NgozMQoyMjc1CjIzNjYKNzI4MAoyNzUyCi01MjYKLTU4NDQKNjAxNQotODIzNwotMzkyNwozODU3Cjk0NjEKMzAzMwotNjg1NQo0NDU4Ci0xOTYyCi03MDgyCjk3NjUKNDYxOAo0NjUzCjI0NTIKLTc0OTUKNjkzMAozOTcyCg=="

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func reverse(x int) int {
	neg := x < 0
	if neg {
		x = -x
	}
	res := 0
	for x > 0 {
		res = res*10 + x%10
		x /= 10
	}
	if neg {
		res = -res
	}
	return res
}

func loadCases() ([]string, []int) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "no embedded testcases")
		os.Exit(1)
	}
	var inputs []string
	var exps []int
	for _, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid number in embedded tests")
			os.Exit(1)
		}
		inputs = append(inputs, fmt.Sprintf("%d\n", val))
		exps = append(exps, reverse(val))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", exps[idx]) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
