package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref628E")
	cmd := exec.Command("go", "build", "-o", exe, "628E.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scan := bufio.NewReader(file)
	var t int
	fmt.Fscan(scan, &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n, m int
		if _, err := fmt.Fscan(scan, &n, &m); err != nil {
			fmt.Printf("bad test file at case %d\n", caseNum)
			os.Exit(1)
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(scan, &rows[i])
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(rows[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
