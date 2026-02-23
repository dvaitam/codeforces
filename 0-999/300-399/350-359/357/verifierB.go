package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// rawTestcases embeds the contents of testcasesB.txt.
var rawTestcases = []string{

	`3 1
1 2 3`,

	`3 1
2 3 1`,

	`3 1
2 3 1`,
}

func verifyOutput(input string, output string) error {
	inReader := bufio.NewReader(strings.NewReader(input))
	outReader := bufio.NewReader(strings.NewReader(output))

	var n, m int
	if _, err := fmt.Fscan(inReader, &n, &m); err != nil {
		return fmt.Errorf("bad input: %v", err)
	}

	dances := make([][3]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(inReader, &dances[i][0], &dances[i][1], &dances[i][2]); err != nil {
			return fmt.Errorf("bad input dances: %v", err)
		}
	}

	colors := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var c int
		if _, err := fmt.Fscan(outReader, &c); err != nil {
			return fmt.Errorf("failed to read color for dancer %d: %v", i, err)
		}
		if c < 1 || c > 3 {
			return fmt.Errorf("invalid color %d for dancer %d", c, i)
		}
		colors[i] = c
	}

	var extra string
	if _, err := fmt.Fscan(outReader, &extra); err == nil {
		return fmt.Errorf("unexpected extra output: %q", extra)
	}

	for i := 0; i < m; i++ {
		c1, c2, c3 := colors[dances[i][0]], colors[dances[i][1]], colors[dances[i][2]]
		if c1 == c2 || c2 == c3 || c1 == c3 {
			return fmt.Errorf("dance %d (%d, %d, %d) has colors (%d, %d, %d) which are not distinct", i+1, dances[i][0], dances[i][1], dances[i][2], c1, c2, c3)
		}
	}

	return nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range rawTestcases {
		input := tc + "\n"
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed\n got: %s\n error: %v\n", idx+1, got, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
