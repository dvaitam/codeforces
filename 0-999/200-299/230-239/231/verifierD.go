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

type TestCaseD struct {
	x, y, z    int
	x1, y1, z1 int
	a          [6]int
	ans        int
}

func solveCaseD(tc TestCaseD) int {
	sum := 0
	if tc.y < 0 {
		sum += tc.a[0]
	}
	if tc.y > tc.y1 {
		sum += tc.a[1]
	}
	if tc.z < 0 {
		sum += tc.a[2]
	}
	if tc.z > tc.z1 {
		sum += tc.a[3]
	}
	if tc.x < 0 {
		sum += tc.a[4]
	}
	if tc.x > tc.x1 {
		sum += tc.a[5]
	}
	return sum
}

func readCasesD(path string) ([]TestCaseD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]TestCaseD, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		x, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		y, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		z, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x1, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		y1, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		z1, _ := strconv.Atoi(scan.Text())
		var a [6]int
		for j := 0; j < 6; j++ {
			scan.Scan()
			a[j], _ = strconv.Atoi(scan.Text())
		}
		tc := TestCaseD{x, y, z, x1, y1, z1, a, 0}
		tc.ans = solveCaseD(tc)
		cases[i] = tc
	}
	return cases, nil
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := readCasesD("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", tc.x, tc.y, tc.z)
		fmt.Fprintf(&sb, "%d %d %d\n", tc.x1, tc.y1, tc.z1)
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		expected := fmt.Sprintf("%d", tc.ans)
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
