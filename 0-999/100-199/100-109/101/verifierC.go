package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	ax, ay, bx, by, cx, cy int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCaseC{
		{1, 0, 2, 0, 1, 0},
		{1, 1, 0, 0, 1, 0},
	}
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d %d %d %d\n", t.ax, t.ay, t.bx, t.by, t.cx, t.cy)
		expect := solveC(t.ax, t.ay, t.bx, t.by, t.cx, t.cy)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveC(ax, ay, bx, by, cx, cy int64) string {
	t := cx*cx + cy*cy
	check := func(x, y int64) bool {
		if t == 0 {
			return x == bx && y == by
		}
		expr1 := cx*bx + cy*by - x*cx - y*cy
		expr2 := bx*cy - by*cx - x*cy + y*cx
		return expr1%t == 0 && expr2%t == 0
	}
	if check(ax, ay) || check(ay, -ax) || check(-ax, -ay) || check(-ay, ax) {
		return "YES\n"
	}
	return "NO\n"
}
