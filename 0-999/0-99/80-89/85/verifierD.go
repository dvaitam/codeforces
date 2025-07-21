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

func buildOracle() (string, error) {
	exe := "oracleD"
	cmd := exec.Command("go", "build", "-o", exe, "./0-999/0-99/80-89/85/85D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) string {
	ops := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", ops)
	present := map[int]bool{}
	keys := []int{}
	for i := 0; i < ops; i++ {
		typ := rng.Intn(3)
		switch typ {
		case 0:
			// add
			var x int
			for {
				x = rng.Intn(1000) + 1
				if !present[x] {
					break
				}
			}
			present[x] = true
			keys = append(keys, x)
			fmt.Fprintf(&sb, "add %d\n", x)
		case 1:
			// del if any else add
			if len(keys) == 0 {
				i--
				continue
			}
			idx := rng.Intn(len(keys))
			x := keys[idx]
			keys = append(keys[:idx], keys[idx+1:]...)
			delete(present, x)
			fmt.Fprintf(&sb, "del %d\n", x)
		default:
			fmt.Fprintf(&sb, "sum\n")
		}
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
