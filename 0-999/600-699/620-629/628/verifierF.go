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
	exe := filepath.Join(os.TempDir(), "ref628F")
	cmd := exec.Command("go", "build", "-o", exe, "628F.go")
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scan := bufio.NewReader(file)
	var t int
	fmt.Fscan(scan, &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n, b, q int
		if _, err := fmt.Fscan(scan, &n, &b, &q); err != nil {
			fmt.Printf("bad test file at case %d\n", caseNum)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, b, q))
		for i := 0; i < q; i++ {
			var up, cnt int
			fmt.Fscan(scan, &up, &cnt)
			sb.WriteString(fmt.Sprintf("%d %d\n", up, cnt))
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
