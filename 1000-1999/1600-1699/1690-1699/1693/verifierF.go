package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func minCost(s string) int {
	n := len(s)
	first := -1
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			first = i
			break
		}
	}
	if first == -1 {
		return 0
	}
	last := -1
	for i := n - 1; i >= 0; i-- {
		if s[i] == '0' {
			last = i
			break
		}
	}
	if last == -1 || first > last {
		return 0
	}
	zerosBefore := first
	onesAfter := n - 1 - last
	zerosInside := 0
	onesInside := 0
	for i := first; i <= last; i++ {
		if s[i] == '0' {
			zerosInside++
		} else {
			onesInside++
		}
	}
	delta := onesInside - zerosInside
	diff := 0
	if delta > zerosBefore {
		diff = delta - zerosBefore
	} else if -delta > onesAfter {
		diff = -delta - onesAfter
	}
	return diff + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(47)
	const t = 100
	var sb strings.Builder
	var exp strings.Builder
	sb.WriteString(fmt.Sprintln(t))
	for i := 0; i < t; i++ {
		n := rand.Intn(10) + 1
		var b strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b.WriteByte('0')
			} else {
				b.WriteByte('1')
			}
		}
		s := b.String()
		sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
		exp.WriteString(fmt.Sprintf("%d\n", minCost(s)))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\noutput:\n%s", err, out.String())
		os.Exit(1)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(exp.String())
	if got != want {
		fmt.Fprintf(os.Stderr, "wrong answer\nexpected:\n%s\ngot:\n%s\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
