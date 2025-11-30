package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "MiA1MzAgOSAzCjMgODY4IDggMTAgMgo1IDM0MiA4IDUgOSA0IDQKNCAxMTQzIDkgOCA3IDMKMiAyOTAgOSA3CjEgMjUgMwo1IDM1MSA1IDEgNSA4IDEwCjQgOTI2IDcgMTAgOCAzCjMgMTM0IDEgMyA4CjIgMjkwIDcgNQo0IDEzMDYgNyAxMCA2IDkKNSA5MTQgMTAgNCA2IDEgNQo1IDY5NyA2IDkgMTAgMTAgMgoyIDQ1MCA1IDUKMSAxMDAgOAo0IDE4NiA2IDIgNyAzCjEgMTY5IDcKNCAzMDYgMSAxMCAxMCAxCjQgMTIwMyA2IDkgNSA5CjIgNTggNSAxCjEgMTQ0IDEwCjUgMzU5IDQgNyA1IDEwIDUKMiAxMjggNiA2CjMgMzg2IDcgNyA4CjUgMTMyNiAxMCA5IDIgMTAgOQozIDUzOCA0IDUgNwozIDg0MiA1IDkgNgoxIDMyNCAxMAozIDM2OSA3IDEwIDEwCjIgMTY0IDYgOAozIDU3MyAxMCA1IDgKMSAxMjEgMQoxIDEyMSA1CjQgNzM3IDEwIDEwIDYgMwozIDM5NiA2IDYgMTAKMyAyODIgNyAyIDEKNSA0NTkgNSA5IDQgNSA0CjMgMTkzIDcgMiAyCjUgNTg1IDYgNCA4IDMgMgozIDQyMSAxMCA4IDUKMiAxMzAgMSA5CjIgMzM3IDEwIDMKMyA0NjQgMiAxMCA2CjUgNTk2IDcgNSA5IDUgOAozIDcxOCA1IDcgMTAKNCAxNTEgNyAzIDQgMQo0IDEyMDcgOSA3IDkgNAoxIDI4OSA5CjMgNTk2IDYgNCAyCjUgNDg3IDIgNCAxIDEgOQoyIDQwNSAxMCAxCjEgMTAwIDIKMiA0MjEgNSA0CjEgMzYxIDkKNCAyNDkgMTAgMiA2IDMKMyA3MDEgOCAxIDYKMiAyMDUgMiA5CjEgNjQgNAozIDM2NSAxIDggMTAKNCAyNzggNSA0IDUgMTAKNSAxMTExIDcgMSA4IDYgMQoxIDI1IDEKMSAxNiAyCjQgMzU4IDIgOSA5IDgKMyAyMzYgNiAyIDYKNCA4NTggMTAgNSA2IDUKMiAyMzMgNyAyCjIgNDEwIDEgNwoxIDE2OSAzCjEgMTk2IDgKNSAxMjIwIDcgMSAxMCA3IDEKMyA2NDYgNiA3IDcKNCAyNDIgNCA0IDkgNQo1IDIzMiA3IDQgNyAzIDEKMyA0MTAgOSA1IDIKNCAyODIgOSA3IDIgNgo1IDExNTggMiAxMCAxIDggMwoyIDM3MCAxIDkKMSAxNDQgMgo0IDE4NiAxIDYgMiAxCjEgMTY5IDUKNSA4MTkgMiAxIDEwIDkgOQoyIDEzNyA5IDIKNSAzNzAgOSA2IDEwIDMgMgoyIDIwOCA0IDgKNSAxMTEwIDUgNiAxMCA3IDYKNSA5ODMgMiA3IDkgNCA3CjIgNjQ4IDEwIDEwCjUgNjg4IDMgNyAzIDMgMgo0IDk5MCA5IDggMTAgMwoyIDE4MSA0IDMKNSAxMzI3IDYgNCA5IDUgNwo1IDExNjcgMTAgNSA0IDUgMQozIDQ5MCA3IDQgMwo1IDY5MCA0IDYgOCAzIDcKNCAxMjgxIDQgOCAxMCA5CjEgMTAwIDIKNCAxODggOCA0IDQgMgoyIDIwMCA0IDQKMyAyNzAgMyAxMCAxCjMgMTc0IDEgNiAzCjQgOTcgMiAyIDIgNQozIDMwOCA2IDggMTAKMyAxMzcgMSA2IDYKNCA4MjQgOCAyIDQgMTAKNCA3MDYgMyA5IDYgMgozIDE5NyA3IDIgOAo1IDc3MyAyIDkgNiA2IDgKMyA0NjQgMiA2IDEwCjUgMTIwNiAyIDggOSA2IDEKMyA1MDcgMyAzIDMKMyA0ODkgMiAyIDkKMiA0MjUgMTAgNwo1IDY3NSAzIDggOCA1IDMKMSAyNSAzCjUgMTMxNCAxMCA3IDYgMiA1CjMgMjgzIDEgMyAxCjQgMTAzOCA1IDQgOSA2CjMgNjQ1IDggOSAyCjMgMzkwIDIgMyA1Cg=="

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solveCase(arr []int64, c int64) int64 {
	calc := func(w int64) int64 {
		var sum int64
		for _, v := range arr {
			d := v + 2*w
			sum += d * d
			if sum > c {
				return sum
			}
		}
		return sum
	}

	low, high := int64(0), int64(1)
	for calc(high) < c {
		high <<= 1
	}
	for low+1 < high {
		mid := (low + high) / 2
		if calc(mid) >= c {
			high = mid
		} else {
			low = mid
		}
	}
	return high
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
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "invalid line %d\n", idx+1)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(fields[0])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx+1)
			os.Exit(1)
		}
		c, errC := strconv.ParseInt(fields[1], 10, 64)
		if errC != nil {
			fmt.Fprintf(os.Stderr, "invalid c on line %d\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != n+2 {
			fmt.Fprintf(os.Stderr, "line %d expected %d numbers, got %d\n", idx+1, n+2, len(fields))
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[i+2], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value on line %d\n", idx+1)
				os.Exit(1)
			}
			arr[i] = val
		}
		input := fmt.Sprintf("1\n%d %d\n", n, c)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		inputs = append(inputs, input)
		exps = append(exps, solveCase(arr, c))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		var got int64
		fmt.Fscan(strings.NewReader(out), &got)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %d\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
