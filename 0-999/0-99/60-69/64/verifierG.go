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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func canonical(path string) (string, bool) {
	parts := strings.Split(path, "/")
	var stack []string
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		if p == ".." {
			if len(stack) == 0 {
				return "-1", false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, p)
		}
	}
	if len(stack) == 0 {
		return "/", true
	}
	return "/" + strings.Join(stack, "/"), true
}

func randToken(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = byte('a' + rng.Intn(26))
		} else {
			b[i] = byte('0' + rng.Intn(10))
		}
	}
	return string(b)
}

func generatePath(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var parts []string
	parts = append(parts, "") // start with root
	for i := 0; i < n; i++ {
		choice := rng.Intn(5)
		switch choice {
		case 0:
			parts = append(parts, ".")
		case 1:
			parts = append(parts, "..")
		default:
			parts = append(parts, randToken(rng))
		}
	}
	return strings.Join(parts, "/")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		path := generatePath(rng)
		input := path + "\n"
		exp, ok := canonical(path)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if !ok {
			if out != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:%s", i+1, out, input)
				os.Exit(1)
			}
			continue
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
