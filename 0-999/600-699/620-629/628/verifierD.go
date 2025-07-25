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

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref628D")
	cmd := exec.Command("go", "build", "-o", exe, "628D.go")
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing m for case %d\n", i)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		a := scan.Text()
		scan.Scan()
		b := scan.Text()
		input := fmt.Sprintf("%d %d\n%s\n%s\n", m, d, a, b)
		expect, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", i, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
