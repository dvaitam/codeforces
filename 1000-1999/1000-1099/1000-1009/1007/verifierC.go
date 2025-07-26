package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCase(bin string, n, a, b int) error {
	cmd := exec.Command(bin)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	fmt.Fprintln(writer, n)
	writer.Flush()

	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 600 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid query: %s", line)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			var res int
			if x < a {
				res = 1
			} else if y < b {
				res = 2
			} else if x > a || y > b {
				res = 3
			} else {
				res = 3
			}
			fmt.Fprintln(writer, res)
			writer.Flush()
		} else if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid answer: %s", line)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			if x == a && y == b {
				stdin.Close()
				cmd.Wait()
				return nil
			}
			cmd.Process.Kill()
			return fmt.Errorf("wrong answer %d %d expected %d %d", x, y, a, b)
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("invalid output: %s", line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
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
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("bad testcase line: %s\n", line)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		a, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		idx++
		if err := runCase(bin, n, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
