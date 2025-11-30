package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "MTAwCjMKMyAzIDAKMgoxIDIKNiAtMwozCi0zIDMgMwoxCi02IC0xCjQKLTMgMSAwIC0yCjEKLTIgMwo0CjMgLTMgLTMgMwozCi0zIDAKNiAwCi0zIDIKMQotMgoyCjUgMwotNSAxCjYKLTIgMCAwIDIgLTEgLTEKMQo0IDIKMQoyCjMKMiAtMwotMiAwCjQgLTEKMgotMSAtMwozCi0zIC0zCi0zIDIKMCAtMQoxCi0yCjIKMyAtMwowIDMKMgotMyAtMwozCjYgMgoxIDIKMyAtMgozCjMgLTMgLTIKMQoxIC0zCjYKMSAtMSAzIDIgLTMgLTMKMgo0IC0yCjUgMQoxCi0yCjIKLTMgLTEKMiAzCjYKMyAtMSAzIDAgMiAtMgoxCjUgMAozCjIgLTIgLTEKMQoxIDIKNAoxIC0xIDEgMgoyCi0zIDEKLTQgMQozCi0xIDIgMAoyCjMgMwoyIC0yCjEKMQoxCjEgLTEKMwoxIC0xIDAKMQo2IDIKNQowIDIgLTIgLTMgMQoyCjYgMQotMyAtMQo2Ci0zIDIgMSAzIDAgLTIKMwo0IDAKNiAxCi00IDEKMwotMSAwIDIKMgotNSAtMwoyIDIKMQotMwoxCi00IDMKMQotMgoxCjEgLTIKMQotMwozCjQgMwowIDIKMSAxCjUKMiAxIC0xIC0zIC0xCjEKMyAyCjEKLTMKMQo0IDAKMQowCjMKLTYgLTEKMiAyCjAgMwozCjMgMiAtMwozCi0yIDEKLTYgMgotNSAzCjMKLTMgLTIgMQoxCi01IDIKMgoxIDAKMQotMyAxCjYKLTIgMSAyIC0yIDEgMgozCjAgMgotNSAzCjQgMQozCjIgMCAyCjMKMCAtMwotNCAzCi00IDEKNAoxIC0xIDIgLTIKMwozIDAKLTEgLTEKMiAtMQozCjAgMyAyCjIKNSAtMQozIC0yCjUKMSAtMiAwIDMgMQozCi0yIC0xCjEgMAowIC0zCjYKLTMgLTEgMSAtMyAxIC0yCjIKLTYgLTEKLTMgLTIKMwotMyAzIDMKMwotMiAxCjIgMQozIDEKMgotMyAtMwoyCi0xIC0yCi00IC0xCjMKMSAtMiAtMgoyCjYgLTEKLTUgMQo0CjEgLTIgMiAtMwozCi0xIDMKMSAtMgoxIDMKNgoyIDEgMyAtMyAxIDIKMgotNiAyCjYgMQo0Ci0zIDMgLTMgMwoyCi0xIDMKMCAwCjEKMwoxCjUgMgo1Ci0zIDIgMyAyIDAKMQotMiAxCjEKLTEKMwozIC0yCi0xIDEKNSAtMgoyCi0zIDIKMwowIDAKNiAxCjMgMQoxCi0yCjIKLTIgMwowIC0yCjIKLTEgLTMKMgozIDIKLTYgMQo0Ci0xIC0zIDAgMQoyCjYgMgotNCAzCjIKMCAxCjEKMSAwCjYKMSAtMiAxIDEgMiAxCjEKLTMgLTMKMwozIDAgLTEKMwozIC0zCi01IDIKMCAwCjMKLTEgLTIgLTEKMgozIDIKLTYgMwozCi0yIC0xIDEKMgozIC0yCjEgLTMKMgotMiAyCjIKLTQgMwo1IDMKMQotMwozCi02IC0yCjMgLTIKLTIgLTMKMwoyIDMgLTMKMgoxIDAKLTYgLTEKMgoyIDEKMQotMyAxCjUKMyAtMSAxIDAgMQozCi00IDIKNSAyCi02IDAKNQoyIDEgLTMgLTEgLTIKMQotMiAtMgo2CjEgLTEgLTMgMCAyIC0zCjEKMCAyCjMKMiAzIDIKMgotNSAtMwotMiAzCjYKMCAtMiAwIDEgMiAtMwoyCjEgMwozIDEKMQotMgoxCjIgMAo1CjIgLTMgLTMgLTIgMQoyCi0xIDAKMCAyCjIKLTMgLTEKMwotMSAtMQo0IC0xCi0zIDEKMwozIC0zIDIKMgowIDMKLTYgMQozCi0yIC0yIC0xCjIKNCAtMgotNSAyCjQKMSAyIDEgMwozCjQgLTIKMCAtMgozIC0yCjQKLTEgMSAtMSAwCjMKLTQgLTIKNCAtMgozIDEKNgozIC0yIDMgLTIgLTEgMgozCjIgLTIKMCAzCi00IDIKNQotMiAtMSAtMyAxIDEKMwotNSAtMgotMSAzCi00IDIKMQozCjIKMCAwCi0yIDAKNgozIDIgMiAtMSAtMiAzCjEKLTUgLTIKMgotMyAzCjMKMyAtMQozIDEKMyAtMgo1Ci0yIC0yIDIgLTMgLTEKMQo1IDMKNAotMiAxIDMgMQoyCjUgMgotMSAtMgo0CjAgMCAtMiAtMgoxCi0zIC0zCjQKMyAtMyAtMSAwCjEKNSAyCjIKMyAtMgoxCi00IDAKNQowIDIgLTEgMyAwCjEKMCAwCjEKMQozCjAgMgotMSAwCjMgMgo1CjEgMSAxIDIgMgoxCjIgLTMKNAowIDMgLTMgLTIKMwoyIDAKNSAtMwoyIDAKNQozIDIgLTMgLTIgLTEKMgowIC0xCjYgMwozCjAgLTIgMAoyCjMgLTMKNSAxCjIKMCAzCjMKMSAtMwo1IDEKLTUgLTIKNQoyIDEgLTIgMCAwCjIKMSAtMwotNCAwCjUKMiAtMiAwIC0zIC0xCjIKNSAyCjUgMwoyCjIgMgoxCi01IC0zCjYKMCAyIC0xIC0zIDEgMgoyCjUgMwoyIDEKMgowIDEKMgotNiAzCi00IC0xCjUKMyAyIDMgLTIgLTIKMQotNCAzCjYKLTMgMCAtMSAtMSAyIDMKMgotMiAyCi00IC0zCjIKLTIgLTMKMwo2IC0xCi02IDAKMSAzCjQKMSAzIC0zIDEKMgotMyAtMgozIDIKNAowIDAgMCAtMwozCjMgLTMKLTMgMQotMiAxCjYKMSAzIDMgLTIgLTMgLTIKMQotNiAyCjMKLTMgLTIgMwoxCjQgMQo="

func isqrt(x int64) int64 {
	if x < 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func solveQuery(freq map[int64]int64, x, y int64) int64 {
	delta := x*x - 4*y
	if delta < 0 {
		return 0
	}
	s := isqrt(delta)
	if s*s != delta {
		return 0
	}
	if (x+s)%2 != 0 || (x-s)%2 != 0 {
		return 0
	}
	r1 := (x + s) / 2
	r2 := (x - s) / 2
	if r1 == r2 {
		c := freq[r1]
		if c >= 2 {
			return c * (c - 1) / 2
		}
		return 0
	}
	return freq[r1] * freq[r2]
}

func expected(a []int64, queries [][2]int64) string {
	freq := make(map[int64]int64)
	for _, v := range a {
		freq[v]++
	}
	ans := make([]string, len(queries))
	for i, q := range queries {
		ans[i] = fmt.Sprintf("%d", solveQuery(freq, q[0], q[1]))
	}
	return strings.Join(ans, " ")
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
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
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], got)
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
			fmt.Fprintf(os.Stderr, "case %d missing array\n", caseNum)
			os.Exit(1)
		}
		arrVals := tokens[pos : pos+n]
		pos += n
		a := make([]int64, n)
		for i, s := range arrVals {
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid array value on case %d\n", caseNum)
				os.Exit(1)
			}
			a[i] = val
		}
		if pos >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing q\n", caseNum)
			os.Exit(1)
		}
		q, errQ := strconv.Atoi(tokens[pos])
		if errQ != nil {
			fmt.Fprintf(os.Stderr, "invalid q on case %d\n", caseNum)
			os.Exit(1)
		}
		pos++
		if pos+2*q > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing queries\n", caseNum)
			os.Exit(1)
		}
		queries := make([][2]int64, q)
		for i := 0; i < q; i++ {
			x, errX := strconv.ParseInt(tokens[pos+2*i], 10, 64)
			y, errY := strconv.ParseInt(tokens[pos+2*i+1], 10, 64)
			if errX != nil || errY != nil {
				fmt.Fprintf(os.Stderr, "invalid query on case %d\n", caseNum)
				os.Exit(1)
			}
			queries[i] = [2]int64{x, y}
		}
		pos += 2 * q
		want := expected(a, queries)
		lines := []string{"1", fmt.Sprintf("%d", n), strings.Join(arrVals, " "), fmt.Sprintf("%d", q)}
		for _, qv := range queries {
			lines = append(lines, fmt.Sprintf("%d %d", qv[0], qv[1]))
		}
		inputs = append(inputs, strings.Join(lines, "\n")+"\n")
		exps = append(exps, want)
	}
	return inputs, exps
}
