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

type block struct {
	id    int
	start int
	size  int
}

func runOps(m int, ops []string) []string {
	blocks := []block{}
	nextID := 1
	var out []string
	for _, op := range ops {
		if strings.HasPrefix(op, "A:") {
			sz, _ := strconv.Atoi(op[2:])
			pos := -1
			if len(blocks) == 0 {
				if sz <= m {
					pos = 0
				}
			} else {
				if blocks[0].start >= sz {
					pos = 0
				} else {
					for j := 0; j < len(blocks)-1 && pos == -1; j++ {
						freeStart := blocks[j].start + blocks[j].size
						freeEnd := blocks[j+1].start
						if freeEnd-freeStart >= sz {
							pos = j + 1
						}
					}
					if pos == -1 {
						lastEnd := blocks[len(blocks)-1].start + blocks[len(blocks)-1].size
						if m-lastEnd >= sz {
							pos = len(blocks)
						}
					}
				}
			}
			if pos == -1 {
				out = append(out, "NULL")
			} else {
				var start int
				if pos == 0 {
					start = 0
				} else {
					prev := blocks[pos-1]
					start = prev.start + prev.size
				}
				newBlk := block{id: nextID, start: start, size: sz}
				blocks = append(blocks, block{})
				copy(blocks[pos+1:], blocks[pos:])
				blocks[pos] = newBlk
				out = append(out, strconv.Itoa(nextID))
				nextID++
			}
		} else if strings.HasPrefix(op, "E:") {
			x, _ := strconv.Atoi(op[2:])
			found := false
			for j, b := range blocks {
				if b.id == x {
					blocks = append(blocks[:j], blocks[j+1:]...)
					found = true
					break
				}
			}
			if !found {
				out = append(out, "ILLEGAL_ERASE_ARGUMENT")
			}
		} else if op == "D" {
			cur := 0
			for j := range blocks {
				blocks[j].start = cur
				cur += blocks[j].size
			}
		}
	}
	return out
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		t, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		ops := parts[2:]
		if len(ops) != t {
			fmt.Printf("test %d invalid op count\n", idx)
			os.Exit(1)
		}
		exp := runOps(m, ops)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", t, m)
		for _, op := range ops {
			if strings.HasPrefix(op, "A:") {
				fmt.Fprintf(&input, "alloc %s\n", op[2:])
			} else if strings.HasPrefix(op, "E:") {
				fmt.Fprintf(&input, "erase %s\n", op[2:])
			} else {
				input.WriteString("defragment\n")
			}
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(outLines) != len(exp) {
			fmt.Printf("Test %d failed: expected %d lines got %d\n", idx, len(exp), len(outLines))
			os.Exit(1)
		}
		for i := range exp {
			if strings.TrimSpace(outLines[i]) != exp[i] {
				fmt.Printf("Test %d failed line %d: expected %s got %s\n", idx, i+1, exp[i], outLines[i])
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
