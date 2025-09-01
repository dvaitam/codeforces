package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Replace oracle-based check with semantic verification that accepts any valid plan.
func verify(vals []int, out string) error {
	// Normalize lines (drop trailing empty lines)
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) != len(vals) {
		return fmt.Errorf("expected %d lines, got %d", len(vals), len(lines))
	}
	// containers
	stack := make([]int, 0)
	queue := make([]int, 0)
	deque := make([]int, 0)
	block := make([]int, 0)

	pushKinds := map[string]bool{
		"pushStack": true,
		"pushQueue": true,
		"pushFront": true,
		"pushBack":  true,
	}
	popKinds := map[string]bool{
		"popStack": true,
		"popQueue": true,
		"popFront": true,
		"popBack":  true,
	}
	containerOf := func(op string) string {
		switch op {
		case "popStack":
			return "stack"
		case "popQueue":
			return "queue"
		case "popFront", "popBack":
			return "deque"
		}
		return ""
	}

	// helpers
	pushFront := func(v int) { deque = append([]int{v}, deque...) }
	pushBack := func(v int) { deque = append(deque, v) }
	popStack := func() (int, bool) {
		if len(stack) == 0 {
			return 0, false
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, true
	}
	popQueue := func() (int, bool) {
		if len(queue) == 0 {
			return 0, false
		}
		v := queue[0]
		queue = queue[1:]
		return v, true
	}
	popFront := func() (int, bool) {
		if len(deque) == 0 {
			return 0, false
		}
		v := deque[0]
		deque = deque[1:]
		return v, true
	}
	popBack := func() (int, bool) {
		if len(deque) == 0 {
			return 0, false
		}
		v := deque[len(deque)-1]
		deque = deque[:len(deque)-1]
		return v, true
	}

	for i, x := range vals {
		line := strings.TrimSpace(lines[i])
		if x != 0 {
			if !pushKinds[line] {
				return fmt.Errorf("line %d: expected push*, got %q", i+1, line)
			}
			// simulate push
			switch line {
			case "pushStack":
				stack = append(stack, x)
			case "pushQueue":
				queue = append(queue, x)
			case "pushFront":
				pushFront(x)
			case "pushBack":
				pushBack(x)
			}
			block = append(block, x)
		} else {
			// zero: must output k ops, where k = min(3, len(block))
			if line == "0" {
				expected := len(block)
				if expected > 3 {
					expected = 3
				}
				if expected != 0 {
					return fmt.Errorf("line %d: reported 0 pops, but need %d", i+1, expected)
				}
			} else {
				parts := strings.Fields(line)
				var k int
				if _, err := fmt.Sscan(parts[0], &k); err != nil {
					return fmt.Errorf("line %d: bad count: %q", i+1, parts[0])
				}
				expected := len(block)
				if expected > 3 {
					expected = 3
				}
				if k != expected {
					return fmt.Errorf("line %d: expected %d pops, got %d", i+1, expected, k)
				}
				if len(parts) != 1+k {
					return fmt.Errorf("line %d: expected %d pop tokens, got %d", i+1, k, len(parts)-1)
				}
				seenContainer := map[string]bool{}
				popped := make([]int, 0, k)
				for j := 0; j < k; j++ {
					op := parts[1+j]
					if !popKinds[op] {
						return fmt.Errorf("line %d: bad op %q", i+1, op)
					}
					c := containerOf(op)
					if seenContainer[c] {
						return fmt.Errorf("line %d: repeated container via %q", i+1, op)
					}
					seenContainer[c] = true
					var v int
					var ok bool
					switch op {
					case "popStack":
						v, ok = popStack()
					case "popQueue":
						v, ok = popQueue()
					case "popFront":
						v, ok = popFront()
					case "popBack":
						v, ok = popBack()
					}
					if !ok {
						return fmt.Errorf("line %d: pop from empty %s", i+1, op)
					}
					popped = append(popped, v)
				}
				// compute expected top k values from current block
				if k > 0 {
					tmp := append([]int(nil), block...)
					// sort desc
					for a := 0; a < len(tmp); a++ {
						for b := a + 1; b < len(tmp); b++ {
							if tmp[b] > tmp[a] {
								tmp[a], tmp[b] = tmp[b], tmp[a]
							}
						}
					}
					expVals := append([]int(nil), tmp[:k]...)
					// sort both popped and expected asc for multiset compare
					asc := func(s []int) {
						for a := 0; a < len(s); a++ {
							for b := a + 1; b < len(s); b++ {
								if s[b] < s[a] {
									s[a], s[b] = s[b], s[a]
								}
							}
						}
					}
					asc(popped)
					asc(expVals)
					if len(popped) != len(expVals) {
						return fmt.Errorf("line %d: internal size mismatch", i+1)
					}
					for j := 0; j < k; j++ {
						if popped[j] != expVals[j] {
							return fmt.Errorf("line %d failed: expected top %v, got %v", i+1, expVals, popped)
						}
					}
				}
			}
			// clear containers and block for next segment
			stack = stack[:0]
			queue = queue[:0]
			deque = deque[:0]
			block = block[:0]
		}
	}
	return nil
}

func run358C(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if len(fields) != 1+n {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[1+i], &vals[i])
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", vals[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run358C(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := verify(vals, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nexpected %s\ngot %s\n", idx, fmt.Errorf("%w", err), "(any valid sequence)", strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
