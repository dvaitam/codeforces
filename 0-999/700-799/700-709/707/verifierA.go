package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateTest() string {
	n := rand.Intn(100) + 1
	m := rand.Intn(100) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	colors := []string{"C", "M", "Y", "W", "G", "B"}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(colors[rand.Intn(len(colors))])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func expectedOutput(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	in.Scan()
	n := atoi(in.Text())
	in.Scan()
	m := atoi(in.Text())
	isColor := false
	for i := 0; i < n*m; i++ {
		if !in.Scan() {
			break
		}
		c := in.Text()[0]
		if c == 'C' || c == 'M' || c == 'Y' {
			isColor = true
		}
	}
	if isColor {
		return "#Color"
	}
	return "#Black&White"
}

func atoi(s string) int {
	v := 0
	for i := 0; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for t := 0; t < 100; t++ {
		input := generateTest()
		want := expectedOutput(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "test", t+1, "error running binary:", err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", t+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
