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

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1185F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(sc.Text())
	for caseNum := 1; caseNum <= T; caseNum++ {
		if !sc.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		m, _ := strconv.Atoi(sc.Text())
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			sc.Scan()
			k, _ := strconv.Atoi(sc.Text())
			sb.WriteString(strconv.Itoa(k))
			for j := 0; j < k; j++ {
				sc.Scan()
				sb.WriteByte(' ')
				sb.WriteString(sc.Text())
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < m; j++ {
			sc.Scan()
			price := sc.Text()
			sc.Scan()
			r := sc.Text()
			sb.WriteString(price + " " + r)
			rInt, _ := strconv.Atoi(r)
			for x := 0; x < rInt; x++ {
				sc.Scan()
				sb.WriteByte(' ')
				sb.WriteString(sc.Text())
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:%sexpected:%s got:%s\n", caseNum, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
