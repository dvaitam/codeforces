package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(a []int) string {
	sum := 0
	zeroReached := false
	ok := true
	for _, x := range a {
		sum += x
		if sum < 0 {
			ok = false
		}
		if zeroReached && sum != 0 {
			ok = false
		}
		if sum == 0 {
			zeroReached = true
		}
	}
	if sum != 0 {
		ok = false
	}
	if ok {
		return "Yes"
	}
	return "No"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	const t = 100
	var input strings.Builder
	var expected strings.Builder
	input.WriteString(fmt.Sprintln(t))
	for i := 0; i < t; i++ {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(7) - 3
		}
		input.WriteString(fmt.Sprintln(n))
		for j, x := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", x))
		}
		input.WriteByte('\n')
		expected.WriteString(solveCase(arr))
		expected.WriteByte('\n')
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\noutput:\n%s", err, out.String())
		os.Exit(1)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(expected.String())
	if normalize(got) != normalize(want) {
		fmt.Fprintf(os.Stderr, "wrong answer\nexpected:\n%s\ngot:\n%s\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}

func normalize(s string) string {
	var res []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		res = append(res, strings.ToLower(strings.TrimSpace(scanner.Text())))
	}
	return strings.Join(res, "\n")
}
