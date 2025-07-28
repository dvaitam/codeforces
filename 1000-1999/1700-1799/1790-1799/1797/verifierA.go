package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func degree(n, m, x, y int) int {
	d := 4
	if x == 1 {
		d--
	}
	if x == n {
		d--
	}
	if y == 1 {
		d--
	}
	if y == m {
		d--
	}
	return d
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierA.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(1)
	const T = 100
	type test struct {
		n, m   int
		x1, y1 int
		x2, y2 int
	}
	tests := make([]test, T)
	var input strings.Builder
	input.WriteString(strconv.Itoa(T) + "\n")
	for i := 0; i < T; i++ {
		n := rand.Intn(7) + 4 // 4..10
		m := rand.Intn(7) + 4
		x1 := rand.Intn(n) + 1
		y1 := rand.Intn(m) + 1
		var x2, y2 int
		for {
			x2 = rand.Intn(n) + 1
			y2 = rand.Intn(m) + 1
			if abs(x1-x2)+abs(y1-y2) >= 2 {
				break
			}
		}
		tests[i] = test{n, m, x1, y1, x2, y2}
		input.WriteString(fmt.Sprintf("%d %d\n%d %d %d %d\n", n, m, x1, y1, x2, y2))
	}

	out, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Println("binary error:", err)
		os.Exit(1)
	}
	fields := strings.Fields(out)
	if len(fields) != T {
		fmt.Printf("wrong number of outputs: got %d want %d\n", len(fields), T)
		os.Exit(1)
	}

	for i := 0; i < T; i++ {
		expect := degree(tests[i].n, tests[i].m, tests[i].x1, tests[i].y1)
		d2 := degree(tests[i].n, tests[i].m, tests[i].x2, tests[i].y2)
		if d2 < expect {
			expect = d2
		}
		got, err := strconv.Atoi(fields[i])
		if err != nil {
			fmt.Println("invalid output")
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
