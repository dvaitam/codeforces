package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesB64 = "YQpiYQpiY2FjYQpiY2IKY2JjYmFhCmJiYmJjYQphYWEKYmFhCmNjY2FiYgpjYmJiYWIKY2NhYmJiY2MKY2JiYmNjCmJjYWJjYWNiCmJiY2NjY2NjCmJjYWJjYmMKYmMKYQphYwpiCmNhY2EKYWFhYmMKYQpiYWFjYWEKYWEKYwpiCmFhY2FjCmIKYQphYWIKYmIKYWJiY2NjYWIKY2NhYmFhYwphYWJhY2MKYmNiYWJiYgpjYWNjYWNhCmFhYWFhCmNhY2NhYWFjCmFiYWNhY2NjCmJjYmJjYQphYmIKYWNjCmFhCmFhCmFhYWMKYmJjY2JhY2EKYmNhY2NhYgphY2IKYWNhY2JiCmJhY2JhYQphY2FiYQpjYmJhY2NhCmIKYgpjYWFiYQpjYgpjYmFiYmFiCmFhCmNiYWNhYQphY2NiYmJiYQpiYmNiY2MKYmNiYmFhYgpjYmNjYwpjYwphYgpjYWIKYWMKYQpiYWNjYwpiYWNiYWEKYWJjYWNjYgphYmMKYmJjYwpiYmMKYwpjYWJjYWJiCmJjY2JjY2IKY2JjY2NhCmNjYWIKY2FiYmNjYgphYWIKY2NjYgphYmNhCmNjY2FiCmFjYwpiYmNjYWJjCmFiYWIKY2IKYmJjYmFhCmFiYmNjYwphYWIKYmNiYWFiYwphYmFjYmNiYQpiY2FiYwphYmJiYgpjY2JjYQpiY2NhCmJhCmJjYWMKYWFhY2IK"

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

func expected(s string) string {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return "NO"
		}
	}
	return "YES"
}

func loadCases() ([]string, []string) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var inputs []string
	var exps []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		inputs = append(inputs, line+"\n")
		exps = append(exps, expected(line))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		if strings.TrimSpace(out) != exps[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
