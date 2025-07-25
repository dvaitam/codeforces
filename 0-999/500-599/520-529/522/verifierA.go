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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomCase(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		if rand.Intn(2) == 0 {
			sb.WriteByte(byte(strings.ToLower(string(ch))[0]))
		} else {
			sb.WriteByte(byte(strings.ToUpper(string(ch))[0]))
		}
	}
	return sb.String()
}

type testCaseA struct {
	input    string
	expected int
}

func computeA(events [][2]string) int {
	depths := map[string]int{"polycarp": 1}
	maxDepth := 1
	for _, ev := range events {
		l1 := strings.ToLower(ev[0])
		l2 := strings.ToLower(ev[1])
		d := depths[l2] + 1
		depths[l1] = d
		if d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}

func generateCaseA() testCaseA {
	n := rand.Intn(20) + 1 // 1..20
	names := []string{"polycarp"}
	events := make([][2]string, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("user%d", i+1)
		parent := names[rand.Intn(len(names))]
		events[i] = [2]string{randomCase(name), randomCase(parent)}
		names = append(names, name)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("%s reposted %s\n", ev[0], ev[1]))
	}
	return testCaseA{input: sb.String(), expected: computeA(events)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCaseA()
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
