package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
}

const solution1054ESource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ys, xs int
	fmt.Fscan(in, &ys, &xs)
	jed := make([]int, ys)
	zer := make([]int, ys)
	chce1 := make([]int, ys)
	chce0 := make([]int, ys)

	type op struct{ y1, x1, y2, x2 int }
	R := make([]op, 0)
	K := make([]op, 0)

	for y := 0; y < ys; y++ {
		opp := ys - 1 - y
		for x := 0; x < xs; x++ {
			var buf string
			fmt.Fscan(in, &buf)
			for z := len(buf) - 1; z >= 0; z-- {
				if buf[z] == '1' {
					if x == 1 {
						jed[opp]++
						R = append(R, op{y, x, opp, x})
					} else {
						jed[y]++
						R = append(R, op{y, x, y, 1})
					}
				} else {
					if x == 0 {
						zer[opp]++
						R = append(R, op{y, x, opp, x})
					} else {
						zer[y]++
						R = append(R, op{y, x, y, 0})
					}
				}
			}
		}
	}
	for y := 0; y < ys; y++ {
		opp := ys - 1 - y
		for x := 0; x < xs; x++ {
			var buf string
			fmt.Fscan(in, &buf)
			for z := len(buf) - 1; z >= 0; z-- {
				if buf[z] == '1' {
					if x == 1 {
						chce1[opp]++
						K = append(K, op{opp, x, y, x})
					} else {
						chce1[y]++
						K = append(K, op{y, 1, y, x})
					}
				} else {
					if x == 0 {
						chce0[opp]++
						K = append(K, op{opp, x, y, x})
					} else {
						chce0[y]++
						K = append(K, op{y, 0, y, x})
					}
				}
			}
		}
	}
	// balance ones
	for ym, yd := 0, 0; ; {
		for ym < ys && jed[ym] >= chce1[ym] {
			ym++
		}
		if ym == ys {
			break
		}
		for jed[yd] <= chce1[yd] {
			yd++
		}
		jed[yd]--
		jed[ym]++
		R = append(R, op{yd, 1, ym, 1})
	}
	// balance zeros
	for ym, yd := 0, 0; ; {
		for ym < ys && zer[ym] >= chce0[ym] {
			ym++
		}
		if ym == ys {
			break
		}
		for zer[yd] <= chce0[yd] {
			yd++
		}
		zer[yd]--
		zer[ym]++
		R = append(R, op{yd, 0, ym, 0})
	}
	total := len(R) + len(K)
	fmt.Fprintln(out, total)
	for _, v := range R {
		// convert to 1-based
		fmt.Fprintf(out, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
	}
	for _, v := range K {
		fmt.Fprintf(out, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1054ESource

var testcases = []testCase{
	{input: "3 2 1 1 0 1 0 1 1 1 0 1 0 1"},
	{input: "3 3 0 01 11 10 00 11 01 0 0 0 01 11 10 00 11 01 0 0"},
	{input: "3 3 1 10 00 0 0 11 1 0 00 1 10 00 0 0 11 1 0 00"},
	{input: "2 3 0 01 10 00 0 1 0 01 10 00 0 1"},
	{input: "3 3 01 0 0 00 0 01 1 0 01 01 0 0 00 0 01 1 0 01"},
	{input: "3 3 10 1 11 10 0 10 0 0 00 10 1 11 10 0 10 0 0 00"},
	{input: "3 3 11 11 0 11 0 1 1 00 0 11 11 0 11 0 1 1 00 0"},
	{input: "3 2 01 11 00 10 1 1 01 11 00 10 1 1"},
	{input: "3 3 0 01 11 1 01 10 1 1 01 0 01 11 1 01 10 1 1 01"},
	{input: "2 3 0 00 10 1 0 01 0 00 10 1 0 01"},
	{input: "3 3 1 10 0 0 01 11 10 0 01 1 10 0 0 01 11 10 0 01"},
	{input: "3 3 11 0 0 1 0 0 1 1 0 11 0 0 1 0 0 1 1 0"},
	{input: "3 3 0 0 0 11 1 0 1 0 1 0 0 0 11 1 0 1 0 1"},
	{input: "3 2 01 0 0 1 0 01 01 0 0 1 0 01"},
	{input: "2 2 11 1 0 11 11 1 0 11"},
	{input: "3 3 1 00 00 1 1 00 1 11 1 1 00 00 1 1 00 1 11 1"},
	{input: "3 3 0 1 10 01 0 01 0 0 0 0 1 10 01 0 01 0 0 0"},
	{input: "3 3 0 0 00 0 1 11 1 00 1 0 0 00 0 1 11 1 00 1"},
	{input: "2 2 0 1 01 11 0 1 01 11"},
	{input: "3 2 1 00 1 0 10 00 1 00 1 0 10 00"},
	{input: "3 3 1 0 11 00 01 0 10 1 00 1 0 11 00 01 0 10 1 00"},
	{input: "2 2 10 1 01 0 10 1 01 0"},
	{input: "3 2 0 1 11 0 0 00 0 1 11 0 0 00"},
	{input: "2 2 0 0 1 1 0 0 1 1"},
	{input: "3 3 00 1 1 11 01 0 0 0 10 00 1 1 11 01 0 0 0 10"},
	{input: "3 3 1 01 1 01 01 1 1 01 10 1 01 1 01 01 1 1 01 10"},
	{input: "2 2 01 01 11 1 01 01 11 1"},
	{input: "3 2 0 10 1 10 10 0 0 10 1 10 10 0"},
	{input: "3 3 01 11 00 00 1 01 1 1 1 01 11 00 00 1 01 1 1 1"},
	{input: "3 2 00 0 11 0 11 1 00 0 11 0 11 1"},
	{input: "3 2 0 1 1 1 0 00 0 1 1 1 0 00"},
	{input: "2 2 00 10 01 0 00 10 01 0"},
	{input: "2 3 1 1 01 01 1 1 1 1 01 01 1 1"},
	{input: "3 3 11 10 0 10 11 0 1 0 10 11 10 0 10 11 0 1 0 10"},
	{input: "2 3 11 00 0 00 10 1 11 00 0 00 10 1"},
	{input: "2 3 00 11 00 1 11 1 00 11 00 1 11 1"},
	{input: "3 3 1 10 11 11 10 01 11 0 01 1 10 11 11 10 01 11 0 01"},
	{input: "2 3 0 00 0 0 1 01 0 00 0 0 1 01"},
	{input: "3 3 00 00 1 11 01 10 1 1 01 00 00 1 11 01 10 1 1 01"},
	{input: "2 3 0 1 1 00 0 0 0 1 1 00 0 0"},
	{input: "3 2 0 0 0 0 1 00 0 0 0 0 1 00"},
	{input: "3 3 0 1 00 0 1 11 00 1 11 0 1 00 0 1 11 00 1 11"},
	{input: "3 2 0 1 0 0 10 10 0 1 0 0 10 10"},
	{input: "2 2 10 1 1 00 10 1 1 00"},
	{input: "2 3 01 10 01 1 01 1 01 10 01 1 01 1"},
	{input: "2 3 01 0 10 1 0 10 01 0 10 1 0 10"},
	{input: "3 2 0 01 0 0 11 10 0 01 0 0 11 10"},
	{input: "2 3 01 1 0 11 0 11 01 1 0 11 0 11"},
	{input: "2 3 0 1 0 1 1 11 0 1 0 1 1 11"},
	{input: "3 3 01 11 10 01 10 1 1 0 11 01 11 10 01 10 1 1 0 11"},
	{input: "3 3 00 01 00 0 01 10 11 00 11 00 01 00 0 01 10 11 00 11"},
	{input: "2 3 0 0 00 1 00 1 0 0 00 1 00 1"},
	{input: "2 3 1 01 01 0 1 10 1 01 01 0 1 10"},
	{input: "2 3 10 1 11 10 01 1 10 1 11 10 01 1"},
	{input: "2 2 11 11 0 01 11 11 0 01"},
	{input: "2 3 0 10 1 1 00 11 0 10 1 1 00 11"},
	{input: "3 3 0 00 1 0 00 01 10 0 0 0 00 1 0 00 01 10 0 0"},
	{input: "2 3 1 0 1 01 10 10 1 0 1 01 10 10"},
	{input: "2 2 1 01 11 0 1 01 11 0"},
	{input: "3 2 00 0 10 00 00 1 00 0 10 00 00 1"},
	{input: "3 3 11 0 00 10 0 11 10 10 1 11 0 00 10 0 11 10 10 1"},
	{input: "2 2 0 0 1 00 0 0 1 00"},
	{input: "2 2 1 11 1 0 1 11 1 0"},
	{input: "2 3 1 0 0 11 11 11 1 0 0 11 11 11"},
	{input: "3 3 11 10 1 11 0 00 0 00 01 11 10 1 11 0 00 0 00 01"},
	{input: "2 3 0 0 0 0 1 00 0 0 0 0 1 00"},
	{input: "3 3 1 0 10 1 0 0 11 1 00 1 0 10 1 0 0 11 1 00"},
	{input: "3 2 0 10 10 01 10 01 0 10 10 01 10 01"},
	{input: "3 2 1 00 1 01 11 1 1 00 1 01 11 1"},
	{input: "2 3 0 11 1 01 1 1 0 11 1 01 1 1"},
	{input: "2 2 1 01 01 0 1 01 01 0"},
	{input: "2 2 0 00 11 11 0 00 11 11"},
	{input: "2 3 10 01 01 00 00 10 10 01 01 00 00 10"},
	{input: "3 3 1 0 1 00 0 1 00 0 10 1 0 1 00 0 1 00 0 10"},
	{input: "2 3 0 0 1 1 0 01 0 0 1 1 0 01"},
	{input: "3 2 01 01 10 01 0 01 01 01 10 01 0 01"},
	{input: "2 3 10 1 0 0 1 00 10 1 0 0 1 00"},
	{input: "2 2 1 10 01 01 1 10 01 01"},
	{input: "2 2 10 1 00 0 10 1 00 0"},
	{input: "2 2 11 00 10 01 11 00 10 01"},
	{input: "2 3 1 01 00 1 00 0 1 01 00 1 00 0"},
	{input: "2 3 1 0 1 01 0 1 1 0 1 01 0 1"},
	{input: "3 3 0 01 1 11 1 1 01 0 0 0 01 1 11 1 1 01 0 0"},
	{input: "3 2 01 00 01 1 0 10 01 00 01 1 0 10"},
	{input: "2 2 10 10 0 01 10 10 0 01"},
	{input: "3 3 00 10 1 11 10 01 1 0 1 00 10 1 11 10 01 1 0 1"},
	{input: "2 3 1 00 11 0 0 0 1 00 11 0 0 0"},
	{input: "3 3 01 0 00 00 00 1 11 0 1 01 0 00 00 00 1 11 0 1"},
	{input: "3 3 1 01 0 1 11 01 01 11 0 1 01 0 1 11 01 01 11 0"},
	{input: "3 3 00 01 0 0 01 0 1 10 1 00 01 0 0 01 0 1 10 1"},
	{input: "3 3 0 00 0 0 00 1 11 10 0 0 00 0 0 00 1 11 10 0"},
	{input: "2 2 0 0 1 11 0 0 1 11"},
	{input: "2 3 0 11 1 1 10 11 0 11 1 1 10 11"},
	{input: "3 3 1 01 00 00 0 1 00 01 10 1 01 00 00 0 1 00 01 10"},
	{input: "3 3 11 1 1 0 11 10 1 11 0 11 1 1 0 11 10 1 11 0"},
	{input: "2 2 11 01 01 01 11 01 01 01"},
	{input: "3 3 01 1 1 11 1 0 00 1 10 01 1 1 11 1 0 00 1 10"},
	{input: "3 2 1 00 01 1 0 0 1 00 01 1 0 0"},
	{input: "3 3 01 11 00 1 1 00 0 1 0 01 11 00 1 1 00 0 1 0"},
	{input: "3 3 01 00 11 01 0 0 00 0 0 01 00 11 01 0 0 00 0 0"},
}

type operation struct {
	y1 int
	x1 int
	y2 int
	x2 int
}

func solveCase(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var ys, xs int
	if _, err := fmt.Fscan(reader, &ys, &xs); err != nil {
		return "", err
	}
	jed := make([]int, ys)
	zer := make([]int, ys)
	chce1 := make([]int, ys)
	chce0 := make([]int, ys)

	R := make([]operation, 0)
	K := make([]operation, 0)

	for y := 0; y < ys; y++ {
		opp := ys - 1 - y
		for x := 0; x < xs; x++ {
			var buf string
			if _, err := fmt.Fscan(reader, &buf); err != nil {
				return "", err
			}
			for z := len(buf) - 1; z >= 0; z-- {
				if buf[z] == '1' {
					if x == 1 {
						jed[opp]++
						R = append(R, operation{y, x, opp, x})
					} else {
						jed[y]++
						R = append(R, operation{y, x, y, 1})
					}
				} else {
					if x == 0 {
						zer[opp]++
						R = append(R, operation{y, x, opp, x})
					} else {
						zer[y]++
						R = append(R, operation{y, x, y, 0})
					}
				}
			}
		}
	}
	for y := 0; y < ys; y++ {
		opp := ys - 1 - y
		for x := 0; x < xs; x++ {
			var buf string
			if _, err := fmt.Fscan(reader, &buf); err != nil {
				return "", err
			}
			for z := len(buf) - 1; z >= 0; z-- {
				if buf[z] == '1' {
					if x == 1 {
						chce1[opp]++
						K = append(K, operation{opp, x, y, x})
					} else {
						chce1[y]++
						K = append(K, operation{y, 1, y, x})
					}
				} else {
					if x == 0 {
						chce0[opp]++
						K = append(K, operation{opp, x, y, x})
					} else {
						chce0[y]++
						K = append(K, operation{y, 0, y, x})
					}
				}
			}
		}
	}

	for ym, yd := 0, 0; ; {
		for ym < ys && jed[ym] >= chce1[ym] {
			ym++
		}
		if ym == ys {
			break
		}
		for jed[yd] <= chce1[yd] {
			yd++
		}
		jed[yd]--
		jed[ym]++
		R = append(R, operation{yd, 1, ym, 1})
	}

	for ym, yd := 0, 0; ; {
		for ym < ys && zer[ym] >= chce0[ym] {
			ym++
		}
		if ym == ys {
			break
		}
		for zer[yd] <= chce0[yd] {
			yd++
		}
		zer[yd]--
		zer[ym]++
		R = append(R, operation{yd, 0, ym, 0})
	}

	var sb strings.Builder
	total := len(R) + len(K)
	fmt.Fprintf(&sb, "%d\n", total)
	for _, v := range R {
		fmt.Fprintf(&sb, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
	}
	for _, v := range K {
		fmt.Fprintf(&sb, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
	}
	return sb.String(), nil
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := tc.input + "\n"
		expected, err := solveCase(input)
		if err != nil {
			fmt.Printf("test %d: failed to compute expected output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := runCase(bin, input, expected); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
