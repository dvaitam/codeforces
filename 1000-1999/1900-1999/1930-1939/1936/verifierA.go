package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCase(outR *bufio.Reader, inW *bufio.Writer, perm []int) error {
	n := len(perm)
	queries := 0
	for {
		line, err := outR.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 3*n {
				return fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) != 5 {
				return fmt.Errorf("invalid query: %s", line)
			}
			a, err1 := strconv.Atoi(parts[1])
			b, err2 := strconv.Atoi(parts[2])
			c, err3 := strconv.Atoi(parts[3])
			d, err4 := strconv.Atoi(parts[4])
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return fmt.Errorf("non-integer query")
			}
			if a < 0 || a >= n || b < 0 || b >= n || c < 0 || c >= n || d < 0 || d >= n {
				return fmt.Errorf("index out of range")
			}
			x := perm[a] | perm[b]
			y := perm[c] | perm[d]
			var resp string
			if x < y {
				resp = "<"
			} else if x > y {
				resp = ">"
			} else {
				resp = "="
			}
			fmt.Fprintln(inW, resp)
			inW.Flush()
		} else if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return fmt.Errorf("invalid answer: %s", line)
			}
			i, err1 := strconv.Atoi(parts[1])
			j, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil || i < 0 || i >= n || j < 0 || j >= n {
				return fmt.Errorf("invalid answer indices")
			}
			best := 0
			for a := 0; a < n; a++ {
				for b := 0; b < n; b++ {
					val := perm[a] ^ perm[b]
					if val > best {
						best = val
					}
				}
			}
			if perm[i]^perm[j] != best {
				return fmt.Errorf("wrong answer")
			}
			return nil
		} else {
			return fmt.Errorf("invalid output: %s", line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	t := rng.Intn(3) + 1
	ns := make([]int, t)
	perms := make([][]int, t)
	for i := 0; i < t; i++ {
		// The problem statement guarantees n \ge 2.  Generating n = 1 causes
		// some valid solutions to terminate unexpectedly (e.g. by panicking
		// on empty candidate sets).  Ensure that the verifier only produces
		// test cases within the specified constraints.
		n := rng.Intn(9) + 2
		perm := rng.Perm(n)
		ns[i] = n
		perms[i] = perm
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("stdin pipe:", err)
		os.Exit(1)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("stdout pipe:", err)
		os.Exit(1)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("start:", err)
		os.Exit(1)
	}
	inW := bufio.NewWriter(stdin)
	outR := bufio.NewReader(stdout)

	fmt.Fprintln(inW, t)
	inW.Flush()
	for i := 0; i < t; i++ {
		fmt.Fprintln(inW, ns[i])
		inW.Flush()
		if err := runCase(outR, inW, perms[i]); err != nil {
			cmd.Process.Kill()
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		fmt.Println("program error:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
