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

const testcasesB64 = "MTAwCjUKLTQgMyA0IC00IC0xCi00IC00IDIgLTIgNAozCi0xIDMgLTIKMCAxIC0xCjMKLTUgMyAyCi0yIDIgMAozCi01IC01IDUKLTQgLTQgNAo1Ci0zIC00IDMgLTQgLTQKLTEgNSAtMiAyIC0xCjQKMiAtNSAtNCAtMwotNSAtMSAwIDUKNAoxIC00IC00IC01Ci01IC0zIDUgNQozCjAgMCAyCjQgLTEgLTQKNAowIC0zIC00IDEKMSAyIC0xIDEKNQoxIDUgLTMgLTQgLTMKLTUgLTMgLTQgMSA0CjUKNCA1IDIgLTMgNAoxIDAgNCAtNSA1CjMKMiAyIC00CjEgMiAtNQoyCi0xIDUKMCAtNQo2CjQgNCAtMiA1IDAgMwozIDQgLTQgMSA1IC0yCjUKMCA1IDUgNSAxCi0zIDMgNCAtNSAtNQozCjMgMiAyCi0zIC00IDIKNAotMiAwIC0xIC01CjMgLTIgMyA1CjUKMSAtMiAtNCAyIDIKNCAyIC00IDIgMgo0Ci00IDIgNSAtNQotNCAxIDEgLTUKNgozIC01IC00IDQgNCA1Ci0xIDMgMyAtMiA1IDUKNQowIDAgLTUgLTIgMgowIDMgNCA0IDIKNAoyIC00IC00IDUKLTIgLTUgMCAwCjQKMyAzIDAgLTQKLTUgLTMgMiAtMwo1Ci00IC0xIDQgLTIgLTIKLTQgMSAtMiAyIC0yCjQKLTQgMSAtNSA1CjUgNCAtNCAyCjUKNSA0IDMgLTMgMgotNSAzIC01IDMgMQo2CjIgMyAtMyA0IC0zIC0zCi00IDEgNSA0IDQgMAo2CjEgMSA0IC0yIC0xIDEKMCAtMSAyIC0zIC0zIDEKNgozIC0xIDMgMCAzIDUKLTIgLTIgLTIgNCAtMSA1CjQKNSAtMyA1IC0zCjMgMyA1IDAKMgowIDMKMiA0CjYKNSAxIDUgMyAtMSAxCi01IDMgLTQgLTUgMSAtMwoyCjUgLTUKLTIgLTEKMwotMSA1IC0xCi0zIC01IC0xCjMKMyAtNSAwCjIgLTQgNAo2Ci0yIDUgMyAxIC0yIDMKLTEgNSAtNSAxIDEgMQo0CjMgLTUgLTUgLTEKLTIgNCAyIC0yCjQKNCAzIDQgLTIKLTIgLTEgMiAtMwoyCi0zIC0zCjMgLTQKNQotNSAxIC0xIDMgLTEKLTMgLTMgNCAtMSAtNQo0CjIgLTMgLTUgLTMKMCA0IC01IDUKNgo1IC0xIDIgNCAzIDAKLTQgMCAzIC0yIC0zIDMKMgozIC0zCjEgMwo2CjEgLTQgMCAtMiAtMiA1CjAgMCAwIC0xIC0yIDQKNgoyIDMgNSAtNSAtNCA1CjAgMiAtMiA1IDQgNAozCi01IDUgMAoxIC00IDEKNAozIC0xIC01IDMKMCAzIDMgMgo1Ci01IC0xIDUgLTIgMAozIC00IC00IC0zIDMKNgotNSA1IC00IC0yIDUgNQotMiA1IDEgLTQgLTIgMwo1CjUgLTQgNSAtMyAtNQoyIDAgLTUgLTQgLTQKMgotMyAtMQozIC00CjYKLTIgLTEgMSAxIDIgNQoxIDEgMCAtNSAxIDUKMgozIC01CjUgNAoyCjQgLTUKLTQgMAo2Ci00IC0xIDUgLTEgNCAtNAotMSAyIDEgMSA1IC01CjUKLTMgMyAtMiAtMyAxCjMgLTEgNSAtMyAtMQo0Ci01IDMgMyAtMwotNCAtNSAtNSA0CjYKMSAzIDMgLTQgLTEgLTIKMCAtMiAxIDUgMyAtMwozCi0zIC0yIDQKLTEgLTIgLTQKMwo1IDQgLTUKMyAyIDUKNgotNCAtMSA1IC00IC0xIC00Ci0yIDQgMiAwIDAgLTMKMwotNCAtMSA0Ci00IC0yIDAKNQoyIDAgNCA1IC0xCjQgLTMgNCA0IDAKNAoxIC0zIC01IDAKLTIgLTIgNSAxCjQKMCAtMyA1IDAKMiAyIDIgMAo0CjIgMyAtNCAtMwo1IDQgLTQgLTEKMwozIC0yIC0xCjUgLTQgLTQKMgozIDUKLTUgNAo2CjEgNCAtMiAtNSAzIC0xCjMgMyAyIDUgLTMgMAo1Ci0yIDQgLTEgLTMgMgozIDMgLTQgMSAwCjYKLTUgNSAtMiA0IDUgLTMKMyAtMyAyIC0zIC0zIDQKNgo1IC0zIDIgMCAtNSAtMgoyIC0yIC00IC0xIDAgLTIKMgo1IC0yCjMgMAo1CjIgMiAtNSAtNSA1CjAgLTQgMyA1IC0xCjQKNCAzIDQgLTMKMSA1IDEgMAo2CjUgLTQgMyAtMSAxIC0yCjMgMSAwIDAgMiAyCjYKLTUgMyAtMyAyIC0zIC0yCi00IDMgNSAtMSAtMiAtMwozCi0zIDAgNAotNCA1IC0yCjUKMCAtNCAyIDUgMgoyIC0yIDQgLTMgNQo1Ci01IC0yIC01IC0zIDQKNCAtMyAtMyA1IDMKMgo0IC0zCi01IC0zCjQKLTMgMyAxIDQKNCA1IC01IDUKNAotNCAtMiAyIDQKMiAyIC0zIC0yCjUKNCAtMyA0IC0zIDQKNSAzIC0xIC0zIDUKNAoyIDQgNCAtNAotMiAxIDEgLTMKMgotNSAtMgozIDMKNQotMyA0IC00IC0zIDIKNCA0IC01IDMgMwoyCi00IDIKNSAtMwo2CjMgLTUgNSA1IDEgMwoxIDAgLTUgMyAxIC0yCjUKMSAwIDQgNCAtNAoxIC0zIC0xIDIgLTIKMgotMSA0CjMgMQo0Ci0yIC01IC01IDQKMyAtNCAzIC01CjUKNSAtMiAwIDEgMAotNCAtMyAxIDUgNAo0Ci0xIDMgMSA0CjAgNCAtMSA1CjUKLTIgMSAzIDMgMgotMyAwIDMgMyAyCjQKLTEgNCAwIDEKLTMgNSAtNCAtMgo1Ci0zIC00IC0yIDMgMQowIDMgNCA1IC01CjMKMSA1IC01Ci00IC0zIDIKNgoxIC00IC00IC0xIC0yIDAKLTQgNCAtNSAxIDEgNAo2Ci0xIDQgLTUgLTMgNSAyCjIgLTUgLTIgNSAzIC0xCg=="

func expected(a, b []int64) string {
	n := len(a)
	diff := make([]int64, n)
	maxd := int64(-1 << 63)
	for i := 0; i < n; i++ {
		diff[i] = a[i] - b[i]
		if diff[i] > maxd {
			maxd = diff[i]
		}
	}
	var idx []string
	for i := 0; i < n; i++ {
		if diff[i] == maxd {
			idx = append(idx, fmt.Sprintf("%d", i+1))
		}
	}
	return fmt.Sprintf("%d\n%s", len(idx), strings.Join(idx, " "))
}

func runCase(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exps[idx] {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

func loadCases() ([]string, []string) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(string(data))
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "no testcases found")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count header\n")
		os.Exit(1)
	}
	pos := 1
	var inputs []string
	var exps []string
	for caseNum := 1; caseNum <= t; caseNum++ {
		if pos >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing n\n", caseNum)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on case %d\n", caseNum)
			os.Exit(1)
		}
		pos++
		if pos+n > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing array a\n", caseNum)
			os.Exit(1)
		}
		aVals := tokens[pos : pos+n]
		pos += n
		if pos+n > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing array b\n", caseNum)
			os.Exit(1)
		}
		bVals := tokens[pos : pos+n]
		pos += n
		arrA := make([]int64, n)
		arrB := make([]int64, n)
		for i := 0; i < n; i++ {
			av, err := strconv.ParseInt(aVals[i], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid a value on case %d\n", caseNum)
				os.Exit(1)
			}
			bv, err := strconv.ParseInt(bVals[i], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid b value on case %d\n", caseNum)
				os.Exit(1)
			}
			arrA[i] = av
			arrB[i] = bv
		}
		want := expected(arrA, arrB)
		inputs = append(inputs, fmt.Sprintf("1\n%d\n%s\n%s\n", n, strings.Join(aVals, " "), strings.Join(bVals, " ")))
		exps = append(exps, want)
	}
	return inputs, exps
}
