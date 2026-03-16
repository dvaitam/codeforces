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

const testcasesRaw = `541 8
636 10
399 8
201 0
467 3
298 0
882 3
994 3
768 0
640 7
785 3
885 1
603 0
780 9
85 4
573 7
796 1
855 6
479 0
96 2
248 7
684 3
2 5
832 8
349 2
670 9
495 2
262 4
105 5
653 4
205 7
837 4
372 4
374 5
756 10
611 0
563 9
11 8
811 1
91 10
860 7
93 2
834 2
191 1
692 8
398 7
569 8
330 3
367 4
311 2
243 5
692 5
31 0
156 1
985 3
142 4
107 5
447 5
715 2
785 3
530 5
943 7
656 7
132 9
194 4
434 9
509 6
275 8
926 0
448 0
306 10
330 2
337 0
881 9
479 9
204 6
169 10
240 1
160 9
3 5
724 1
544 7
365 9
746 1
85 9
619 5
881 3
392 9
284 9
27 7
489 10
943 2
309 4
365 8
395 0
334 2
874 0
37 7
94 2
249 8`

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
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

const mod = 1000000007

func modPow(a, b int) int {
	res := 1
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			panic("invalid case")
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		expected := modPow(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
