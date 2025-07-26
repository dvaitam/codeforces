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

func expected(k int, primes []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", k))
	for _, p := range primes {
		sb.WriteByte(' ')
		sb.WriteString(p)
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randPrimeString() string {
	// generate a random 3-digit prime string for simplicity
	primes := []string{"2", "3", "5", "7", "11", "13", "17", "19", "23", "29"}
	return primes[rand.Intn(len(primes))]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		k := rand.Intn(5) + 2
		primes := make([]string, k)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", k))
		for i := 0; i < k; i++ {
			primes[i] = randPrimeString()
			input.WriteString(primes[i])
			if i+1 < k {
				input.WriteByte(' ')
			}
		}
		input.WriteByte('\n')
		expect := expected(k, primes)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Println("input:\n", input.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("wrong answer on test %d\n", t+1)
			fmt.Println("input:\n", input.String())
			fmt.Printf("expected: %s\n got: %s\n", expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
