package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n, m   int
	x1, y1 int
	x2, y2 int
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		n := rng.Intn(8) + 2
		m := rng.Intn(8) + 2
		x1 := rng.Intn(n) + 1
		y1 := rng.Intn(m) + 1
		x2 := rng.Intn(n) + 1
		y2 := rng.Intn(m) + 1
		for x1 == x2 && y1 == y2 {
			x2 = rng.Intn(n) + 1
			y2 = rng.Intn(m) + 1
		}
		cases[i] = testCase{n: n, m: m, x1: x1, y1: y1, x2: x2, y2: y2}
	}
	return cases
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func interact(bin string, cases []testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Start(); err != nil {
		return err
	}
	w := bufio.NewWriter(stdin)
	r := bufio.NewReader(stdout)

	fmt.Fprintf(w, "%d\n", len(cases))
	w.Flush()
	for idx, tc := range cases {
		fmt.Fprintf(w, "%d %d\n", tc.n, tc.m)
		w.Flush()
		queries := 0
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return fmt.Errorf("case %d: failed to read: %v", idx+1, err)
			}
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "?") {
				queries++
				if queries > 4 {
					return fmt.Errorf("case %d: too many queries", idx+1)
				}
				var x, y int
				if _, err := fmt.Sscanf(line, "? %d %d", &x, &y); err != nil {
					return fmt.Errorf("case %d: invalid query %q", idx+1, line)
				}
				d1 := abs(x-tc.x1) + abs(y-tc.y1)
				d2 := abs(x-tc.x2) + abs(y-tc.y2)
				d := d1
				if d2 < d {
					d = d2
				}
				fmt.Fprintf(w, "%d\n", d)
				w.Flush()
			} else if strings.HasPrefix(line, "!") {
				var x, y int
				if _, err := fmt.Sscanf(line, "! %d %d", &x, &y); err != nil {
					return fmt.Errorf("case %d: invalid answer %q", idx+1, line)
				}
				if (x != tc.x1 || y != tc.y1) && (x != tc.x2 || y != tc.y2) {
					return fmt.Errorf("case %d: wrong cell", idx+1)
				}
				fmt.Fprintln(w, "Ok")
				w.Flush()
				break
			} else if line != "" {
				return fmt.Errorf("case %d: unexpected output %q", idx+1, line)
			}
		}
	}
	stdin.Close()
	outLeft, _ := io.ReadAll(stdout)
	errLeft, _ := io.ReadAll(&errBuf)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, string(errLeft)+string(errBuf.Bytes()))
	}
	if strings.TrimSpace(string(outLeft)) != "" {
		return fmt.Errorf("extra output: %s", string(outLeft))
	}
	if errBuf.Len() > 0 {
		return fmt.Errorf(errBuf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	if err := interact(bin, cases); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
