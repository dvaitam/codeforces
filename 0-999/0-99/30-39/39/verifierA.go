package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// solveA implements the logic of 39A using the code from 39A.go.
func solveA(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var a0 int
	var expr string
	if _, err := fmt.Fscan(reader, &a0, &expr); err != nil {
		return ""
	}
	type Task struct {
		w    int
		c    int
		pre  bool
		sign int
	}
	tasks := make([]Task, 0, len(expr))
	i := 0
	n := len(expr)
	sign := 1
	for i < n {
		coef := 1
		start := i
		tmp := 0
		for i < n && expr[i] >= '0' && expr[i] <= '9' {
			tmp = tmp*10 + int(expr[i]-'0')
			i++
		}
		if i < n && expr[i] == '*' && i > start {
			coef = tmp
			i++
		} else {
			i = start
			coef = 1
		}
		pre := false
		if i+1 < n && expr[i] == '+' && expr[i+1] == '+' {
			pre = true
			i += 2
			i++
		} else if i < n && expr[i] == 'a' {
			i++
			if i+1 < n && expr[i] == '+' && expr[i+1] == '+' {
				pre = false
				i += 2
			}
		}
		tasks = append(tasks, Task{w: sign * coef, c: coef, pre: pre, sign: sign})
		if i < n {
			if expr[i] == '+' {
				sign = 1
			} else if expr[i] == '-' {
				sign = -1
			}
			i++
		}
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].w < tasks[j].w })
	a := a0
	var total int64
	for _, t := range tasks {
		if t.pre {
			a++
			total += int64(t.sign) * int64(t.c) * int64(a)
		} else {
			total += int64(t.sign) * int64(t.c) * int64(a)
			a++
		}
	}
	return fmt.Sprintf("%d\n", total)
}

func generateCaseA(rng *rand.Rand) string {
	a := rng.Intn(2001) - 1000
	terms := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", a))
	for i := 0; i < terms; i++ {
		if i > 0 {
			if rng.Intn(2) == 0 {
				sb.WriteByte('+')
			} else {
				sb.WriteByte('-')
			}
		}
		coef := 1
		if rng.Intn(3) == 0 {
			coef = rng.Intn(5) + 1
		}
		if coef != 1 {
			sb.WriteString(fmt.Sprintf("%d*", coef))
		}
		if rng.Intn(2) == 0 {
			sb.WriteString("a++")
		} else {
			sb.WriteString("++a")
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseA(rng)
	}
	for i, tc := range cases {
		expect := solveA(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
