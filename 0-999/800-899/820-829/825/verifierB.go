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

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "825B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genBoard(rng *rand.Rand) string {
	// cross and nought counts equal
	cross := rng.Intn(49) + 1
	nought := cross
	cells := make([]byte, 100)
	for i := 0; i < cross; i++ {
		for {
			pos := rng.Intn(100)
			if cells[pos] == 0 {
				cells[pos] = 'X'
				break
			}
		}
	}
	for i := 0; i < nought; i++ {
		for {
			pos := rng.Intn(100)
			if cells[pos] == 0 {
				cells[pos] = 'O'
				break
			}
		}
	}
	for i := 0; i < 100; i++ {
		if cells[i] == 0 {
			cells[i] = '.'
		}
	}
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.Write(cells[i*10 : i*10+10])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genBoard(rng)
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", in)
			os.Exit(1)
		}
		got, err := run(exe, in)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", in)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", in)
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
