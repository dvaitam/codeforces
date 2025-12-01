package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./207D2.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			reportFailure("reference failed on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			reportFailure("candidate failed on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
		}

		if normalize(refOut) != normalize(candOut) {
			reportFailure("wrong answer on test %d (%s)\nInput:\n%sExpected: %s\nGot: %s",
				idx+1, tc.name, tc.input, normalize(refOut), normalize(candOut))
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func reportFailure(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "207D2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []testCase {
	return []testCase{
		{
			name: "culture opera",
			input: doc(11, "Opera Review",
				"The opera theater reopened with a stunning ballet sequence and new choreography.",
				"Artists and musicians described the orchestra rehearsals and costume design.",
			),
		},
		{
			name: "culture museum",
			input: doc(22, "Museum Guide",
				"A national museum gallery features sculptures, literature readings, and film retrospectives.",
				"Curators invite visitors to explore design studios and art workshops all weekend.",
			),
		},
		{
			name: "culture festival",
			input: doc(33, "Festival Chronicle",
				"A citywide art festival showcases poetry recitals, theater rehearsals, and nightly concerts.",
				"Young writers debut novels while choirs perform folk songs in public squares.",
			),
		},
		{
			name: "politics parliament",
			input: doc(44, "Parliament Briefing",
				"The parliament gathered to vote on election reforms and minister appointments after debate.",
				"A coalition of senators discussed federal policy and national security council reports.",
			),
		},
		{
			name: "politics diplomacy",
			input: doc(55, "Diplomatic Note",
				"Foreign ministers met at the embassy to negotiate border security and defense treaties.",
				"The president instructed the cabinet to brief the security council and governors.",
			),
		},
		{
			name: "politics campaign",
			input: doc(66, "Campaign Trail",
				"Municipal governors launched campaigns with referendum pledges and policy manifestos.",
				"Opposition parties debated constitutional amendments and military spending targets.",
			),
		},
		{
			name: "economy bulletin",
			input: doc(77, "Market Bulletin",
				"Global markets monitored trade balances, corporate profits, and stock indexes all week.",
				"Logistics firms reported revenue growth from shipping contracts and supply deals.",
			),
		},
		{
			name: "economy finance",
			input: doc(88, "Finance Report",
				"The central bank adjusted credit rates to stabilize currency exchange and inflation.",
				"Investors discussed capital budgets, venture funding, and long term loans.",
			),
		},
		{
			name: "economy agriculture",
			input: doc(99, "Harvest Ledger",
				"Agricultural companies estimated harvest output, export demand, and commodity pricing.",
				"Suppliers negotiated energy contracts for oil and gas deliveries to rural factories.",
			),
		},
	}
}

func doc(id int, title string, lines ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n%s\n", id, title)
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}
