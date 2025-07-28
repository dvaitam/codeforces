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

func expected(n int) (aW, aB, bW, bB int) {
	pos := 1
	step := 1
	remaining := n
	for remaining > 0 {
		take := step
		if take > remaining {
			take = remaining
		}
		l := pos
		r := pos + take - 1
		white := (r+1)/2 - (l)/2
		black := take - white
		var alice bool
		if step == 1 {
			alice = true
		} else {
			pair := (step - 2) / 2
			if pair%2 == 1 {
				alice = true
			} else {
				alice = false
			}
		}
		if alice {
			aW += white
			aB += black
		} else {
			bW += white
			bB += black
		}
		pos += take
		remaining -= take
		step++
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA2.txt")
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
		n, _ := strconv.Atoi(line)
		aW, aB, bW, bB := expected(n)
		input := fmt.Sprintf("1\n%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		expect := fmt.Sprintf("%d %d %d %d", aW, aB, bW, bB)
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
