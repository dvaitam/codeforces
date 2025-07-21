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

var months = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

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

func generateCase(rng *rand.Rand) (string, string) {
	month := months[rng.Intn(len(months))]
	k := rng.Intn(300)
	input := fmt.Sprintf("%s\n%d\n", month, k)
	idx := 0
	for i, m := range months {
		if m == month {
			idx = i
			break
		}
	}
	expected := months[(idx+k)%12]
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	ref := "45A.go"
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		outRef, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if outRef != exp {
			fmt.Fprintf(os.Stderr, "reference output mismatch on case %d\n", i+1)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong answer: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
