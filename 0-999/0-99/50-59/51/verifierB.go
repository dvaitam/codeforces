package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func genTable(rng *rand.Rand, depth int) (string, []int) {
	var sb strings.Builder
	counts := []int{}
	cell := 0
	sb.WriteString("<table>")
	rows := rng.Intn(2) + 1
	for i := 0; i < rows; i++ {
		sb.WriteString("<tr>")
		cols := rng.Intn(2) + 1
		for j := 0; j < cols; j++ {
			sb.WriteString("<td>")
			cell++
			if depth < 2 && rng.Float32() < 0.3 {
				sub, subc := genTable(rng, depth+1)
				sb.WriteString(sub)
				counts = append(counts, subc...)
			}
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")
	counts = append(counts, cell)
	return sb.String(), counts
}

func generateCase(rng *rand.Rand) (string, string) {
	markup, counts := genTable(rng, 0)
	sort.Ints(counts)
	var exp strings.Builder
	for i, v := range counts {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	expected := exp.String()
	input := markup + "\n"
	return input, expected
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
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
