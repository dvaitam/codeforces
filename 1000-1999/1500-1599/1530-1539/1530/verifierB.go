package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testB struct{ h, w int }

func genTestsB() []testB {
	rand.Seed(1530002)
	tests := make([]testB, 100)
	for i := range tests {
		tests[i].h = rand.Intn(18) + 3 // 3..20
		tests[i].w = rand.Intn(18) + 3
	}
	return tests
}

func solveB(tc testB) []string {
	h, w := tc.h, tc.w
	grid := make([]string, h)
	var sb strings.Builder
	for i := 0; i < h; i++ {
		sb.Reset()
		for j := 0; j < w; j++ {
			var val byte
			if (i == 0 || i == h-1) && (j == 0 || j == w-1) {
				val = '1'
			} else if i == 0 || i == h-1 {
				if j == 1 || j == w-2 {
					val = '0'
				} else if j%2 == 0 {
					val = '1'
				} else {
					val = '0'
				}
			} else if j == 0 || j == w-1 {
				if i == 1 || i == h-2 {
					val = '0'
				} else if i%2 == 0 {
					val = '1'
				} else {
					val = '0'
				}
			} else {
				val = '0'
			}
			sb.WriteByte(val)
		}
		grid[i] = sb.String()
	}
	return grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.h, tc.w)
	}

	var expected bytes.Buffer
	for _, tc := range tests {
		grid := solveB(tc)
		for _, row := range grid {
			expected.WriteString(row)
			expected.WriteByte('\n')
		}
		expected.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(expected.String())
	if got != want {
		fmt.Fprintln(os.Stderr, "wrong answer")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
