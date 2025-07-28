package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	h int
	m int
	s string
}

func generateTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 100)
	edge := []Test{{24, 60, "23:59"}, {10, 60, "00:00"}, {12, 12, "11:11"}, {5, 3, "02:02"}}
	tests = append(tests, edge...)
	for len(tests) < 100 {
		h := rand.Intn(99) + 1
		m := rand.Intn(99) + 1
		hh := rand.Intn(h)
		mm := rand.Intn(m)
		tests = append(tests, Test{h, m, fmt.Sprintf("%02d:%02d", hh, mm)})
	}
	return tests
}

var mirror = map[int]int{0: 0, 1: 1, 2: 5, 5: 2, 8: 8}

func solve(h, m int, s string) string {
	hour := int((s[0]-'0')*10 + (s[1] - '0'))
	minute := int((s[3]-'0')*10 + (s[4] - '0'))
	for i := 0; i < h*m; i++ {
		h1 := hour / 10
		h2 := hour % 10
		m1 := minute / 10
		m2 := minute % 10
		d1, ok1 := mirror[m2]
		d2, ok2 := mirror[m1]
		d3, ok3 := mirror[h2]
		d4, ok4 := mirror[h1]
		if ok1 && ok2 && ok3 && ok4 {
			rh := d1*10 + d2
			rm := d3*10 + d4
			if rh < h && rm < m {
				return fmt.Sprintf("%02d:%02d", hour, minute)
			}
		}
		minute++
		if minute == m {
			minute = 0
			hour++
			if hour == h {
				hour = 0
			}
		}
	}
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	var in strings.Builder
	fmt.Fprintln(&in, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&in, "%d %d\n%s\n", t.h, t.m, t.s)
	}
	expectedParts := make([]string, len(tests))
	for i, t := range tests {
		expectedParts[i] = solve(t.h, t.m, t.s) + "\n"
	}
	expect := strings.Join(expectedParts, "")

	got, err := run(binary, in.String())
	if err != nil {
		fmt.Printf("runtime error: %v\noutput:\n%s", err, got)
		os.Exit(1)
	}
	got = strings.ReplaceAll(strings.TrimSpace(got), "\r\n", "\n")
	expect = strings.ReplaceAll(strings.TrimSpace(expect), "\r\n", "\n")
	if got != expect {
		fmt.Println("wrong answer")
		fmt.Println("input:")
		fmt.Print(in.String())
		fmt.Println("expected:")
		fmt.Print(expect)
		fmt.Println("got:")
		fmt.Print(got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0)
}
