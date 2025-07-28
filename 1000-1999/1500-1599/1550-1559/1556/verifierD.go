package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func runCase(bin string, input string) error {
	parts := strings.Split(strings.TrimSpace(input), "\n")
	header := strings.Fields(parts[0])
	if len(header) != 2 {
		return fmt.Errorf("invalid input")
	}
	n, _ := strconv.Atoi(header[0])
	k, _ := strconv.Atoi(header[1])
	arrStr := strings.Fields(parts[1])
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i], _ = strconv.Atoi(arrStr[i])
	}
	sorted := append([]int(nil), arr...)
	sort.Ints(sorted)
	target := sorted[k-1]

	cmd := exec.Command(bin)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, fmt.Sprintf("%d %d\n", n, k))
	rdr := bufio.NewReader(stdout)
	queries := 0
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			return fmt.Errorf("runtime error: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "and", "or":
			if len(parts) != 3 {
				return fmt.Errorf("invalid query")
			}
			i, _ := strconv.Atoi(parts[1])
			j, _ := strconv.Atoi(parts[2])
			if i < 1 || i > n || j < 1 || j > n || i == j {
				return fmt.Errorf("invalid indices")
			}
			var val int
			if parts[0] == "and" {
				val = arr[i-1] & arr[j-1]
			} else {
				val = arr[i-1] | arr[j-1]
			}
			fmt.Fprintln(stdin, val)
			queries++
			if queries > 2*n {
				return fmt.Errorf("too many queries")
			}
		case "finish":
			if len(parts) != 2 {
				return fmt.Errorf("invalid finish")
			}
			guess, _ := strconv.Atoi(parts[1])
			if guess != target {
				return fmt.Errorf("wrong answer")
			}
			stdin.Close()
			return cmd.Wait()
		default:
			return fmt.Errorf("invalid command %s", parts[0])
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, _ := filepath.Glob("tests/D/*.in")
	total := 0
	for _, f := range cases {
		data, _ := os.ReadFile(f)
		if err := runCase(bin, string(data)); err != nil {
			fmt.Printf("%s: %v\n", f, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("OK %d cases\n", total)
}
