package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "NyA1IDIgNiAzIDEgNCA3CjcgMiAxIDUgNiA0IDMgNwo5IDcgNiA0IDkgOCAxIDIgNSAzCjUgNCAyIDMgNSAxCjggMSA2IDUgOCA3IDQgMyAyCjggNCAyIDMgNiA3IDggMSA1CjEgMQoxMCAxMCA1IDcgMiAxIDMgOSA0IDYgOAo0IDEgNCAzIDIKMiAyIDEKNiA2IDMgMiAxIDQgNQoyIDEgMgo5IDMgNiAyIDcgOCAxIDkgNSA0CjUgNCAzIDEgNSAyCjEwIDQgMyA3IDEwIDkgNiAxIDIgOCA1CjkgMSA5IDMgNCA2IDIgOCA1IDcKOCA1IDQgMiAxIDMgNyA2IDgKOCA0IDUgMyAxIDggMiA3IDYKMiAyIDEKNiA2IDUgMSA0IDMgMgozIDIgMSAzCjEgMQoxMCAzIDQgNSA3IDEwIDYgOCAxIDIgOQo2IDIgMyA0IDUgNiAxCjIgMSAyCjQgMiA0IDMgMQo5IDYgNCA1IDggMSAzIDkgMiA3CjYgNiA1IDMgMSAyIDQKMSAxCjEwIDUgMSA2IDEwIDggOSAzIDQgNyAyCjEgMQozIDMgMiAxCjkgNiAzIDcgMSA4IDQgOSAyIDUKMTAgNyAxIDEwIDIgMyA4IDQgNiA1IDkKOCA4IDQgNSA3IDEgNiAzIDIKOCAxIDQgNyAyIDggMyA1IDYKNyAyIDcgMyA1IDEgNCA2CjYgMSA0IDMgNSA2IDIKNyAzIDIgMSA0IDUgNyA2CjEwIDYgNSAxMCAyIDggOSA0IDMgMSA3CjkgNSAyIDggNyA2IDMgNCA5IDEKNCAyIDQgMyAxCjEgMQo3IDIgNSA3IDEgMyA0IDYKMSAxCjkgOCA3IDYgMyA5IDUgMSA0IDIKMiAxIDIKNyA0IDUgMiAzIDcgMSA2CjUgNSA0IDEgMyAyCjUgMiAzIDQgNSAxCjUgNCAxIDIgMyA1CjEwIDEgOSA0IDYgMiA1IDMgNyAxMCA4CjEgMQozIDMgMiAxCjYgNSA0IDIgMSA2IDMKMSAxCjUgMSA0IDMgNSAyCjYgMyAxIDYgMiA1IDQKOCA4IDYgNSA3IDIgMyAxIDQKNiA2IDIgNSAxIDQgMwo4IDggMiA1IDQgMSAzIDcgNgo3IDUgNCA3IDIgNiAzIDEKMiAyIDEKMiAyIDEKNCAzIDQgMiAxCjIgMSAyCjkgMyAxIDYgMiA4IDcgNCA1IDkKNiA1IDEgNCA2IDMgMgo2IDQgMiAzIDUgNiAxCjQgNCAzIDIgMQo1IDUgMyA0IDEgMgoyIDIgMQo4IDQgNiA4IDEgMiA1IDMgNwozIDEgMiAzCjcgNyAxIDIgNiAzIDQgNQo5IDIgOSA4IDUgNyA2IDMgMSA0CjUgMSAzIDQgMiA1CjggOCAzIDQgNiA1IDEgNyAyCjMgMyAyIDEKMSAxCjggMSAzIDUgNCA4IDIgNyA2CjkgNCAyIDMgNiA3IDggNSAxIDkKNCAxIDIgNCAzCjggNiA0IDMgMSAyIDggNSA3CjMgMSAzIDIKNiA2IDUgMiAzIDQgMQo3IDYgNSAzIDIgMSA0IDcKNCAxIDIgNCAzCjEgMQoxIDEKOCA4IDYgMyAxIDUgNCAyIDcKNiA1IDMgMiA2IDEgNAoxIDEKMiAxIDIKMyAyIDMgMQoxMCA0IDIgNyA5IDUgOCAxIDYgMyAxMAo4IDggNyA2IDMgNCA1IDEgMgoxMCA2IDQgMTAgMiAxIDkgOCA3IDMgNQo4IDggMSA1IDMgNCAyIDcgNgoxIDEK"

func expectedMoves(p []int) int {
	fixed := 0
	for i, v := range p {
		if v == i+1 {
			fixed++
		}
	}
	moves := fixed / 2
	if fixed%2 == 1 {
		moves++
	}
	return moves
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func loadCases() ([]string, []int) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var inputs []string
	var exps []int
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		n, errN := strconv.Atoi(fields[0])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "line %d expected %d numbers, got %d\n", idx+1, n+1, len(fields))
			os.Exit(1)
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value on line %d\n", idx+1)
				os.Exit(1)
			}
			p[i] = val
		}
		expected := expectedMoves(p)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", n)
		for i, v := range p {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')
		inputs = append(inputs, input.String())
		exps = append(exps, expected)
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", exps[idx]) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
