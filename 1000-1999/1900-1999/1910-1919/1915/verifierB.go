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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(rows [3]string) string {
	letters := []rune{'A', 'B', 'C'}
	for i := 0; i < 3; i++ {
		if strings.ContainsRune(rows[i], '?') {
			used := map[rune]bool{}
			for _, ch := range rows[i] {
				if ch != '?' {
					used[ch] = true
				}
			}
			for _, l := range letters {
				if !used[l] {
					return string(l)
				}
			}
		}
	}
	return ""
}

func generateCase(rng *rand.Rand) (string, string) {
	letters := []byte{'A', 'B', 'C'}
	rng.Shuffle(3, func(i, j int) { letters[i], letters[j] = letters[j], letters[i] })
	rows := [3]string{}
	for i := 0; i < 3; i++ {
		row := []byte{letters[(i)%3], letters[(i+1)%3], letters[(i+2)%3]}
		rows[i] = string(row)
	}
	r := rng.Intn(3)
	c := rng.Intn(3)
	rowBytes := []byte(rows[r])
	rowBytes[c] = '?'
	rows[r] = string(rowBytes)

	var sb strings.Builder
	sb.WriteString("1\n")
	for i := 0; i < 3; i++ {
		sb.WriteString(rows[i])
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solveCase(rows)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
