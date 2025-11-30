package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "NTMyNzQ3NjYwNDAwCjk4ODkzOTM1NDAwMwo3MzIzNzgwOTQ2NzEKNjI5NTEwMTY0MzUzCjIzOTk4NTA2NTU3Ngo5OTc3MDYwOTUzOQo5MDA5ODI1MzExMjMKOTc4NzU3MjczODM0Cjg1NzM5MzI3MzU4Mgo3NjY2NjIwMjU0ODIKNDYxODEwODI0NTgwCjMzNzE4NTQ1NTg3OQoxMjg2MTQ4ODExMTgKNDY4Nzc3MTcwODYzCjYyMjMwNTcyMzMwNgo0NTU3MTY0OTM5OTUKMTA3NjQ0NDMzOTE1Cjg1MjE4NjYyOTQ0Mwo4MDM4Mjk1NzM0MTcKMzQxNjQ1NDg2MTUKNDkzMDI1MDYzMzk3Cjc1MzQ3MDc3ODQwMQozMTg1NTg1NjI0Ngo3OTE2Njc2OTU2MTQKODY5ODQxOTQyNDAKNzg4MjM0MDAyOTYKMzkxMzYzMzM4NDc0CjMzMDM1OTIyNzI0CjM4MzczNTc2NDc3MAo5MzU0MDQyMTc4Cjg5ODYzODEyNzA1Mgo3ODg4MDUyMTAxNwo5ODE4MTQ4NzM4NzEKMjI4MjQ4Njk3NDY2CjIyMzM1MjE4MTE2OAo3NDE1NjQ4Mzc3NjYKMTM3MDEzMzM5MTQ4Cjc1MDc1MDMyNzkKNDA0OTg2NDgwNjMxCjI4NzMxMDM5NzQ4CjY2NTQzMTI1NzY5Mgo5NDE1OTc5NzcyMzgKMjAyNDcyODY0ODI2CjEyMjIwOTIwNTkxOAozODAwMDQ1MDY3OTMKMTQyODQyOTYyMjcxCjM5NjAzMDg2MDM4Mwo1MjExMjk5MTM0MDYKMzIxOTU0NjQ0NTAwCjY5ODE2MDY5NzI3NAoyMDMyNjc4NDkxMDcKMTEyMDE2MjIxNzY3CjYzNzk0NTM1NjA5OQoxNzMxMjA1MjE3MjgKOTg1MTY0OTQ3OTI3Cjg4MTAwNjE4MDE4MAozNDQ1NTQzNzcyMzkKMjY4NDcxMjQzNTkwCjgyOTk0NTEyMzAyOAozMTg2MTc2MTUxODUKNDYxMTYxMjEyNDYxCjUwMDkzMDQzMzI0CjE0NTQzMTc3ODU1OAoyNDA1NzAwNjEyMQo4MzI5NTc0MTc4Ngo4MDMyNTIwNjQ2OAo0NjAxMjg1MjM0NDQKNjAyNTgxMjc3MzYwCjgxMzUzODc0MzY1MgoxNTg1OTc2NjQ5MTAKNDY2Mzk1MTIyMzAzCjcwMTM1OTY5MDIzOQo5MTcxNzE1NDUwNAo0ODYzOTY3MDc2MTYKNDA2NDQ0NjkwOTMwCjYyNDAyNDcwMTk2CjQ0ODI5MzM2MDg1NQo0NTUzMDI3Mjk2OTAKODAyOTk1NDg5MDEyCjM1NjA1ODIxMDUyMwoyMjUyMzM2NTcyNDIKMzIzNzE4NTQyMzQzCjUxOTQ4NjM0ODY4Mwo4NzI2NzQ4OTgwMDAKMzAxMTE0MjA5NzUzCjYxMDM2NzA3NDE1Mgo3NTg1MTUwNDM3NjMKODcyNTM5MjY1NjM1CjQ5MjY0MzkwNzU1NAo0NDIwNzM1NDU4MDUKODQyNjA5Nzk0MDE1CjQ3NDI1NzYwODc3NgoyNzEzMzM0MzAyMTEKNTAyMzM3MzY2MTY0CjE1Njg2NjU4Nzk5Nwo1MDgzMzI0MzI1MDUKNzAyNzkyMDEzNzYxCjUyODY1MjQ4OTQ5NwoyMjY1ODExNTAzNTEKMTI2NTg0ODUzMAo="

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func buildLCMs(limit uint64) []uint64 {
	lcms := []uint64{1}
	for i := 1; ; i++ {
		prev := lcms[len(lcms)-1]
		g := gcd(prev, uint64(i))
		l := prev / g * uint64(i)
		if l > limit {
			break
		}
		lcms = append(lcms, l)
	}
	return lcms
}

func expectedAnswer(n uint64, lcms []uint64) int {
	ans := 0
	for i := 1; i < len(lcms); i++ {
		if n%lcms[i] == 0 {
			ans = i
		} else {
			break
		}
	}
	return ans
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
	lcms := buildLCMs(1e18)
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n64, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid number on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expectedAnswer(n64, lcms)
		inputs = append(inputs, fmt.Sprintf("1\n%d\n", n64))
		exps = append(exps, exp)
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
