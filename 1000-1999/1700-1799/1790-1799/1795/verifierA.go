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

func countDup(s string) int {
	c := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			c++
		}
	}
	return c
}

func expected(n, m int, s, t string) string {
	ds := countDup(s)
	dt := countDup(t)
	if ds+dt > 1 || (ds+dt == 1 && s[len(s)-1] == t[len(t)-1]) {
		return "NO"
	}
	return "YES"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		var sBuilder strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sBuilder.WriteByte('B')
			} else {
				sBuilder.WriteByte('R')
			}
		}
		s := sBuilder.String()
		var tBuilder strings.Builder
		for i := 0; i < m; i++ {
			if rng.Intn(2) == 0 {
				tBuilder.WriteByte('B')
			} else {
				tBuilder.WriteByte('R')
			}
		}
		u := tBuilder.String()

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		input.WriteString(s)
		input.WriteByte('\n')
		input.WriteString(u)
		input.WriteByte('\n')

		expectedOut := expected(n, m, s, u)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tcase+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tcase+1, expectedOut, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
