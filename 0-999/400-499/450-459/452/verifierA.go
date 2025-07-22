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

var pokemons = []string{"vaporeon", "jolteon", "flareon", "espeon", "umbreon", "leafeon", "glaceon", "sylveon"}

func expected(n int, pattern string) string {
	for _, p := range pokemons {
		if len(p) != n {
			continue
		}
		ok := true
		for i := 0; i < n; i++ {
			if pattern[i] != '.' && pattern[i] != p[i] {
				ok = false
				break
			}
		}
		if ok {
			return p
		}
	}
	return ""
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		name := pokemons[rng.Intn(len(pokemons))]
		n := len(name)
		b := []byte(name)
		for i := range b {
			if rng.Intn(2) == 0 {
				b[i] = '.'
			}
		}
		pattern := string(b)
		// ensure unique solution
		cnt := 0
		var ans string
		for _, p := range pokemons {
			if len(p) != n {
				continue
			}
			match := true
			for i := 0; i < n; i++ {
				if pattern[i] != '.' && pattern[i] != p[i] {
					match = false
					break
				}
			}
			if match {
				cnt++
				ans = p
			}
		}
		if cnt == 1 {
			input := fmt.Sprintf("%d\n%s\n", n, pattern)
			return input, ans
		}
	}
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
