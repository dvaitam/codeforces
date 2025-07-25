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

type eventA struct {
	t    int
	v, u int64
	w    int64
}

func solveA(events []eventA) string {
	weights := make(map[int64]int64)
	var sb strings.Builder
	for _, e := range events {
		if e.t == 1 {
			v, u := e.v, e.u
			w := e.w
			for v != u {
				if v > u {
					weights[v] += w
					v /= 2
				} else {
					weights[u] += w
					u /= 2
				}
			}
		} else {
			v, u := e.v, e.u
			var res int64
			for v != u {
				if v > u {
					res += weights[v]
					v /= 2
				} else {
					res += weights[u]
					u /= 2
				}
			}
			fmt.Fprintln(&sb, res)
		}
	}
	return sb.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		events := make([]eventA, q)
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			typ, _ := strconv.Atoi(scan.Text())
			if typ == 1 {
				scan.Scan()
				v, _ := strconv.ParseInt(scan.Text(), 10, 64)
				scan.Scan()
				u, _ := strconv.ParseInt(scan.Text(), 10, 64)
				scan.Scan()
				w, _ := strconv.ParseInt(scan.Text(), 10, 64)
				events[i] = eventA{t: 1, v: v, u: u, w: w}
				inSB.WriteString(fmt.Sprintf("1 %d %d %d\n", v, u, w))
			} else {
				scan.Scan()
				v, _ := strconv.ParseInt(scan.Text(), 10, 64)
				scan.Scan()
				u, _ := strconv.ParseInt(scan.Text(), 10, 64)
				events[i] = eventA{t: 2, v: v, u: u}
				inSB.WriteString(fmt.Sprintf("2 %d %d\n", v, u))
			}
		}
		input := inSB.String()
		expected := solveA(events)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
