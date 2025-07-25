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

type Passenger struct {
	f int
	t int
}

func expected(n, s int, passengers []Passenger) int {
	arrival := make([]int, s+1)
	for _, p := range passengers {
		if p.t > arrival[p.f] {
			arrival[p.f] = p.t
		}
	}
	time := 0
	for floor := s; floor > 0; floor-- {
		if arrival[floor] > time {
			time = arrival[floor]
		}
		time++
	}
	return time
}

func runCandidate(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func parseCase(line string) (int, int, []Passenger, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0, 0, nil, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("bad n")
	}
	s, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("bad s")
	}
	if len(parts) != 2+n*2 {
		return 0, 0, nil, fmt.Errorf("expected %d pairs got %d", n, (len(parts)-2)/2)
	}
	passengers := make([]Passenger, n)
	for i := 0; i < n; i++ {
		f, err := strconv.Atoi(parts[2+2*i])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("bad floor on pair %d", i+1)
		}
		t, err := strconv.Atoi(parts[2+2*i+1])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("bad time on pair %d", i+1)
		}
		passengers[i] = Passenger{f: f, t: t}
	}
	return n, s, passengers, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
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
		n, s, passengers, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect := expected(n, s, passengers)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, s))
		for _, p := range passengers {
			input.WriteString(fmt.Sprintf("%d %d\n", p.f, p.t))
		}
		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
