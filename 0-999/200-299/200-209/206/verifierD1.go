package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseD1 struct {
	input  string
	expect string
}

func solveD1(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var id int
	if _, err := fmt.Fscan(reader, &id); err != nil {
		return ""
	}
	var name string
	if _, err := fmt.Fscan(reader, &name); err != nil {
		return ""
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
	}
	return "1"
}

func makeDoc(id int, name string, lines []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", id))
	sb.WriteString(fmt.Sprintf("%s\n", name))
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genTests() []testCaseD1 {
	rand.Seed(42)
	base := []testCaseD1{
		{
			input: makeDoc(1, "Trade",
				[]string{"Wheat prices rising worldwide.", "Investors stay cautious."}),
			expect: "1",
		},
		{
			input: makeDoc(1024, "Research",
				[]string{"New materials exhibit superconductivity.", "More studies required."}),
			expect: "1",
		},
		{
			input: makeDoc(999999, "LogisticsReport",
				[]string{"Shipment delayed due to weather"}),
			expect: "1",
		},
	}
	tests := make([]testCaseD1, 0, 100)
	for _, t := range base {
		tests = append(tests, testCaseD1{
			input:  t.input,
			expect: solveD1(t.input),
		})
	}
	for i := 0; i < 50; i++ {
		id := rand.Intn(1_000_000)
		name := fmt.Sprintf("Doc%d", i)
		lines := rand.Intn(5) + 1
		var body []string
		for j := 0; j < lines; j++ {
			body = append(body, fmt.Sprintf("Line %d content %d", j, rand.Int()))
		}
		input := makeDoc(id, name, body)
		tests = append(tests, testCaseD1{
			input:  input,
			expect: solveD1(input),
		})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
