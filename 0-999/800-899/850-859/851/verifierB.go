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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(ax, ay, bx, by, cx, cy int64) string {
	abx := bx - ax
	aby := by - ay
	bcx := cx - bx
	bcy := cy - by
	d1 := abx*abx + aby*aby
	d2 := bcx*bcx + bcy*bcy
	cross := abx*bcy - aby*bcx
	if d1 == d2 && cross != 0 {
		return "Yes"
	}
	return "No"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	tCases, _ := strconv.Atoi(scan.Text())
	for i := 0; i < tCases; i++ {
		var vals [6]int64
		for j := 0; j < 6; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			vals[j] = v
		}
		ax, ay, bx, by, cx, cy := vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]
		input := fmt.Sprintf("%d %d %d %d %d %d\n", ax, ay, bx, by, cx, cy)
		expected := solveCase(ax, ay, bx, by, cx, cy)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
