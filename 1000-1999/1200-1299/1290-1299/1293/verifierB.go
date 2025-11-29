package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var testcases = []int{
	138, 583, 868, 822, 783, 65, 262, 121, 508, 780, 461, 484, 668, 389, 808, 215, 97, 500, 30, 915,
	856, 400, 444, 623, 781, 786, 3, 713, 457, 273, 739, 822, 235, 606, 968, 105, 924, 326, 32, 23,
	27, 666, 555, 10, 962, 903, 391, 703, 222, 993, 433, 744, 30, 541, 228, 783, 449, 962, 508, 567,
	239, 354, 237, 694, 225, 780, 471, 976, 297, 949, 23, 427, 858, 939, 570, 945, 658, 103, 191, 645,
	742, 881, 304, 124, 761, 341, 918, 739, 997, 729, 513, 959, 991, 433, 520, 850, 933, 687, 195, 311,
}

const testcasesCount = 100

func solveCase(n int) string {
	var ans float64
	for i := 1; i <= n; i++ {
		ans += 1.0 / float64(i)
	}
	return fmt.Sprintf("%.12f", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	if len(testcases) != testcasesCount {
		fmt.Printf("unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}
	binary := os.Args[1]
	for i, n := range testcases {
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(outBuf.String())
		expected := solveCase(n)
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
