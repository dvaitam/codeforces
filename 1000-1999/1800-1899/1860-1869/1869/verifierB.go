package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "OSA0IDEgOSA1IDQgNCAtMiAtMiAtNCA1IDMgMyAxIDMgLTEgLTQgLTMgMSA0IDEgLTQKMyAzIDEgMiAxIC0zIC01IDIgMSA1CjggMCA4IDYgLTEgLTQgMCAtNCAtNCAwIC01IDAgMCAtMyAtNSAtMiA0IC0zIC0yIC01CjUgNSAxIDIgLTEgMCAtNSA0IC0yIC0zIC0zIDIgLTQgMgo3IDQgMiAxIC0yIDAgMCAyIC0xIC0xIDMgNSAwIC0zIDQgLTQgLTQgMwo2IDEgNCAyIC0zIC0yIDAgMyAtMiAtMiAtMyAtMSAwIDEgNSAtNQo0IDQgMSA0IC00IC00IC0zIDEgLTEgMyAxIC0zCjggNCA2IDIgLTIgMiA1IDAgNSAzIC01IDEgMSAtNSAxIDAgMiAtMiAwIC0xCjkgMSAzIDIgLTEgLTQgMyA0IC0zIDIgMSAtMyAxIDEgLTMgLTIgMiAwIDMgLTMgMCAyCjMgMyAxIDIgLTUgMiA0IDIgLTUgLTIKNiAwIDYgMyAzIDQgLTMgMSAyIC00IDUgMiAtMiAzIDEgLTEKMiAwIDIgMSAtNSAtMSAxIDMKOCA3IDIgNSAwIC0xIDUgLTIgNCAtNCAtNSAtNCAtMSAtMSAzIDAgLTQgMyAtMiAtMwozIDMgMiAxIDMgLTMgNCAzIDUgLTIKMTAgMSA3IDkgMSAtMSAtMSAyIDAgNCA1IC0zIC0zIC00IC00IDEgMSA0IDIgLTMgMyA1IC0xIDAKOSA2IDQgOCAyIDMgMCAyIDUgLTUgMiAtMSAtMyAyIC01IDQgLTIgLTUgMSAtNSAzIC00CjMgMyAxIDIgLTUgLTQgNCAtNSAtMSA1CjYgNSAyIDEgNCAtMSAtMiAtNCAxIDIgMCAxIC0zIDAgMSA1CjggMiA4IDMgMyAwIC0zIC0yIC0zIDIgMCAxIDEgMiAxIC0yIC0yIDIgLTIgNAoyIDEgMSAyIDUgLTQgLTMgMAoyIDIgMSAyIDQgLTEgNCAtNAoxMCA0IDYgNyAyIC01IDUgMyA1IDUgMyAxIDQgMiAyIC0xIDIgLTIgMCAtMSAtNSAtNSAtNSAtMwo3IDAgMyA2IC01IC0zIC00IDEgNSAtMiA0IDEgMyAtMiAyIC0yIDAgNAozIDAgMiAxIDMgMiAwIC0xIC01IDMKMiAwIDIgMSAtMiAzIDAgLTIKNSAyIDMgMSAzIDEgLTEgMiAwIC0yIC01IC0xIDMgLTQKMiAxIDIgMSAtNSAxIDIgMgo5IDEgMiAxIC0yIC00IC0zIDEgLTIgMiA0IC00IDEgMyAxIC01IC0zIC0yIDIgLTIgLTMgLTEKNyA1IDQgMSAzIC0xIDQgMyAtMiAtMSAyIDMgNCAyIDMgNSAtMSAtMQo1IDAgMSA1IC00IC0zIDEgLTIgLTIgLTEgNSAtNSAzIDMKOCAwIDIgNyA1IC0xIC00IDQgMCAtMiA1IDMgLTIgLTIgLTQgMyAtMSA1IDAgNQo5IDQgMyAxIC01IDMgMyAwIDAgNCA1IC01IC0zIDEgLTMgLTMgMyAtNCAtMyAtMiAyIDQKNSAxIDIgMSAxIDAgNCA0IC0zIDUgMiAtNCA0IC01CjEwIDkgNiA4IDIgLTEgLTUgLTIgMyA1IC0zIDUgMiAyIDMgMCAtNCAtMSAtMyA0IDEgLTIgMCAtMQo4IDAgNCAxIDAgLTIgMCAyIDUgNSA1IC0yIC0xIDAgNSAtMyAtMSAtNSAwIDQKMTAgMCAzIDYgLTUgMiA1IC01IC01IC0yIC01IC01IC0yIDUgMCAtNCAtNSAwIDUgMSAtMyAtMiAyIDEKNCAyIDMgMiA1IDAgMSAxIC01IDEgLTEgMwoxMCAxMCA4IDEgNCAtNCAxIDEgLTMgLTUgMyAtMyA0IDUgLTQgMCAtMiAtMyAtMiAtNSAtMyA1IC00IDEKMyAzIDMgMSA0IDQgLTUgLTEgMCAxCjIgMiAxIDIgLTQgMCAtMSA1CjQgMyAyIDMgLTMgMSAwIC0xIDIgMSAtNSAtMQoxMCA0IDkgOCAtNSAzIDQgMyAtMSA1IC01IDIgMSAtNCAxIDAgMiAtNSAtNSAtMSA1IDUgNCAtMQo1IDQgNSAzIDEgLTEgLTIgLTQgNCAwIC0yIDQgNSAzCjcgMiA3IDIgMCAtNSA0IC01IDQgLTMgMCAwIC0xIDUgLTEgMCAyIDEKOCAyIDEgMyA0IC01IDIgLTMgMCAtNSAyIDUgNSAtMSA0IC0yIC00IDMgMSAtMQo0IDQgMiAxIDUgNSAtMyA0IC00IDMgNSAzCjggNiA1IDEgLTEgLTUgMSAtMSAtMSAzIDMgMyAwIDAgLTIgMSAtMyAtNSAzIC0zCjggNSA4IDEgMyAxIDUgNCAtMiAtNSAwIDMgLTMgNSAtMiA1IDAgNSAyIC01CjUgNCAyIDMgLTMgMSAtNCA0IDIgLTIgMiAzIC00IC0yCjQgMyAxIDQgNSAxIC0xIC0xIDEgMCA0IDAKMyAyIDEgMiAtNSAtMSAtMiAxIDEgMQo4IDAgOCA2IDQgLTMgNCAtMSAwIC01IDEgMiAzIC0zIC01IC00IDQgMCAtNCAtMgozIDMgMSAyIC01IDAgMSAtMyA1IC0xCjggMiAzIDcgLTEgMyAtNSAtMyAtMyAtMyAyIDUgLTUgMyA1IDEgLTMgMCA0IC00CjMgMSAyIDEgLTEgMCAtMSAtMSAzIDIKNCAzIDIgMSA1IDQgLTMgNSAzIC01IDAgLTQKNSA1IDQgNSAtMiAyIDMgLTMgMCA1IC0zIDIgMyAtNQoxMCAxIDkgNiAtNSAtNCAtNCAxIDQgMCA0IDIgMCAxIDMgMCA1IC00IC0zIDAgLTUgLTMgLTMgLTUKNyAzIDEgNCA1IC01IC0xIDEgLTUgNCAtMyAwIC00IDEgLTUgMiAwIDQKNiA1IDMgNSAyIDEgLTMgLTUgMiAtMSAtMiAxIC00IDAgLTQgLTQKMiAxIDEgMiAxIDQgNSAtNQo3IDcgNyA1IDIgMiAtNCAtNSAzIDEgLTEgLTUgNSAzIC00IC00IDAgMAozIDMgMSAyIDMgNSAtMSAtNSAtNSAxCjcgMiA1IDYgLTMgLTMgNSA1IC0yIDQgMCAtNSAyIDUgNSAxIC01IC0yCjUgNSAzIDEgLTMgLTIgMCAtMiAtMyAxIDIgMCA0IC0zCjggMCAzIDEgNSAxIC0zIC0zIC01IC01IDAgMyAtNSAtNCA0IDQgLTUgMSAxIDQKNyAzIDIgMyAzIC0yIDMgMSAtNCAtMyAxIDQgNSAwIC00IDEgMSAtMgo5IDYgNCA3IC0yIDUgMiAxIDQgLTQgLTEgLTEgMyAwIDMgLTUgNCA0IDIgLTIgLTEgLTUKNyA2IDYgMSAzIC01IC0zIDEgLTUgMSAxIDEgLTQgMiA0IDIgLTMgLTMKNyA3IDQgMiA0IC0xIDMgLTQgMCAwIC0zIDUgMCAyIDUgMyAtNSAtMgo2IDEgNiA1IDAgLTEgMSA1IC01IC0xIDMgMSAtNSA1IDEgLTEKOCAzIDYgMyAtMyAtNCA0IDAgLTMgLTUgMSA0IDEgMiAtNCA1IDUgLTQgMSAzCjEwIDIgMyAxIC0yIC0zIC0yIC01IDMgLTMgMiAwIDQgLTEgMCA1IC00IDEgLTEgLTMgMyAwIDEgMwo2IDEgNCAzIDEgMiAzIC0xIDEgMSAtNCA1IC0zIC0zIC01IDQKMTAgMSA0IDEwIDUgNCAzIC00IC0xIDQgLTMgMSAtNCAtNSAtNSAtNCAwIDIgMCAtNCA1IDIgMCA0CjYgNSA0IDIgLTMgMyAzIC00IDMgLTMgLTUgNSAxIDQgNSAtMgo5IDYgNSAxIDQgLTMgMSAtMyAyIDQgLTUgMSAtNCA1IDQgMSAwIC0yIDMgMiAtNSAyCjMgMiAzIDIgMyA1IDEgMiAtMSAtMwo1IDQgMyAyIC0xIDQgLTMgMiAtNCAtNCAyIDEgNCAxCjEwIDEgNSA4IC0yIC00IC0xIC0zIDAgLTQgLTMgNSAtNSA1IC0zIDQgMyAtMiAtNSAtNSAxIDMgNCAyCjMgMyAzIDIgMCAtNCA1IC01IC0yIC0yCjkgNCA1IDQgLTUgMiAwIDMgMCAtNCAtNCAtMSA0IDEgLTIgMCAxIC0zIC0yIC0xIC0yIDIKNyA0IDQgNSAtMyAtNCAxIDAgMyAyIC0yIDUgMCA1IDAgMSAtMSAwCjggNCAyIDggLTEgLTQgMiAtMyAwIC0yIC0zIDAgMiAtMiAtNCA1IDEgMSAyIDMKOSA5IDQgNyAzIC0xIDIgLTIgMCAzIDUgLTUgLTQgMiAwIDEgLTIgMSAtNSA0IC01IDEKMyAyIDEgMyAwIC0zIC00IC0zIDAgLTUKNSAwIDEgNCAzIC01IC0zIC00IDQgNCAtMiAtNCAyIC0yCjIgMiAyIDEgMyAtMyAtMiAtMgo4IDYgOCAxIDEgLTIgMSAtNSA0IC0xIC01IDQgMCAwIDAgNSAyIDUgLTMgNAoxMCAxIDUgMiAtNCAtMSAtNSAtMyA0IDUgLTMgMSAtMiA0IDUgMCAtMiAxIDMgMyAtNCAzIC00IDIKMyAzIDIgMSAyIDMgMCAtMyAxIC0xCjggMSA2IDQgMiAtMiAwIDIgMSAtNSAyIC01IDMgMSAxIC0yIC0xIDIgMiAtMwo1IDMgMyAxIDUgMiA0IC0zIDMgNSAtNCAtMiAtMSAzCjMgMCAxIDIgLTUgMCA0IC0zIDIgMQo1IDMgNCAyIC0zIDEgLTUgMCA1IDMgLTIgMCAtMyAzCjEwIDkgMTAgNyAtMyA1IDQgMiA1IDQgLTIgNCAyIC01IC0yIDIgNCAwIDMgLTIgLTUgLTUgLTUgLTMKNSAzIDUgMiA1IDUgLTQgMyAxIDUgNSAtMSA1IDAKOSAzIDMgNiAzIDAgMiAxIDIgNSAwIDIgLTQgLTMgLTIgLTQgLTMgMSAyIDMgLTIgMAoyIDAgMiAxIC00IDEgNCA0CjggOCA0IDEgNSAtMyAzIDUgMiAtMSAxIDMgMyAtNCAtNCA0IDIgLTUgLTIgMQo="

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n, k, a, b int, xs, ys []int64) int64 {
	direct := abs(xs[a]-xs[b]) + abs(ys[a]-ys[b])
	const inf int64 = 1 << 62
	dA, dB := inf, inf
	for i := 1; i <= k; i++ {
		da := abs(xs[a]-xs[i]) + abs(ys[a]-ys[i])
		if da < dA {
			dA = da
		}
		db := abs(xs[b]-xs[i]) + abs(ys[b]-ys[i])
		if db < dB {
			dB = db
		}
	}
	if dA+dB < direct {
		return dA + dB
	}
	return direct
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", idx+1, out)
			os.Exit(1)
		}
		if val != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, exps[idx], val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

func loadCases() ([]string, []int64) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var inputs []string
	var exps []int64
	for caseIdx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", caseIdx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		a, _ := strconv.Atoi(fields[2])
		b, _ := strconv.Atoi(fields[3])
		if len(fields) != 4+2*n {
			fmt.Fprintf(os.Stderr, "case %d: wrong number of coordinates\n", caseIdx+1)
			os.Exit(1)
		}
		xs := make([]int64, n+1)
		ys := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			x64, _ := strconv.ParseInt(fields[4+2*(i-1)], 10, 64)
			y64, _ := strconv.ParseInt(fields[4+2*(i-1)+1], 10, 64)
			xs[i] = x64
			ys[i] = y64
		}
		input := buildInput(n, k, a, b, xs, ys)
		inputs = append(inputs, input)
		exps = append(exps, expected(n, k, a, b, xs, ys))
	}
	return inputs, exps
}

func buildInput(n, k, a, b int, xs, ys []int64) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, a, b))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	return sb.String()
}
