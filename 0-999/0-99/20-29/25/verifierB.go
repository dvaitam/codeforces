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

func isValidGrouping(s, result string) bool {
	cleanResult := strings.ReplaceAll(result, "-", "")
	if cleanResult != s {
		return false
	}
	
	groups := strings.Split(result, "-")
	for _, group := range groups {
		if len(group) != 2 && len(group) != 3 {
			return false
		}
	}
	return true
}

func expected(s string) string {
	n := len(s)
	result := make([]byte, 0, n+n/2)
	i := 0
	
	for i < n {
		rem := n - i
		if rem == 3 {
			result = append(result, s[i], s[i+1], s[i+2])
			break
		}
		if rem == 2 {
			result = append(result, s[i], s[i+1])
			break
		}
		result = append(result, s[i], s[i+1], '-')
		i += 2
	}
	return string(result)
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(99) + 2
		s := make([]byte, n)
		for j := 0; j < n; j++ {
			s[j] = byte('0' + rng.Intn(10))
		}
		
		input := fmt.Sprintf("%d\n%s\n", n, string(s))
		
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		
		if !isValidGrouping(string(s), got) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid grouping %s for input %s\ninput:\n%s", i+1, got, string(s), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}