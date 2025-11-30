package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "NiAwIDUgMyAtNSAyIC0yIDUgLTUgLTMgLTQgMCAyCjUgMSAzIC00IDQgLTIgLTUgLTIgMSAtMSAtMwo4IC0zIC00IC0zIDQgNCAyIC0zIC0zIC01IC01IC0yIC0yIC0xIDAgLTIgMwo1IC0zIC0yIDEgLTEgLTUgMCAxIC0zIC0zIC0xCjMgMCAtMSA0IDQgLTUgNAo3IC00IC0xIDAgLTEgMiAwIC0zIDIgMiAtMyAtNSAtMSAtNSAwCjggLTUgMyAxIDAgMSA0IC01IDIgLTUgLTMgNCAtMiAtNCAtMiAyIDAKMTAgMCAzIC0xIDIgLTQgNCAwIC0xIC01IDEgLTQgLTIgNCAwIC0zIDAgLTEgMyAtNCAtMQo3IC0xIC0zIC00IDUgLTMgLTEgMiAtMyAtNSAtNCA0IDMgMSAtNQo1IDQgMCAtMSAyIDUgMSAtMyAtNSA1IC01CjkgMCAtMiAtMyA0IC0zIDUgMSAtNCAtMyAxIDAgLTMgLTUgMSAtMSAtMyAyIDQKNCAzIDIgMiAwIDIgLTEgLTEgMgo4IC0zIC00IDEgMyAtMyA1IDIgMCAyIC0xIDMgMyAzIDAgLTQgMAoyIC0xIDAgMyA1CjYgMiAtMSAtMSAwIDUgLTMgNCAtNSAyIDMgNSAtMQo5IC0xIDMgNSA1IDAgMCAtMSA1IDAgMSAwIC0zIDIgMCAwIDMgLTMgMwo0IC0yIDAgMiAtMSAtNCA1IDEgLTMKMTAgNSAxIC0xIDQgMyA1IC0xIC01IC0yIC0zIDQgMiA0IDUgLTMgLTIgNSAtMyA1IC01CjkgLTIgLTMgLTUgLTMgLTQgMCAtMyAyIC0yIDMgLTUgMSAyIDAgMSA1IDQgLTQKNSAtMiAwIC01IDAgMSAtMSAxIC00IDMgMAoyIDMgNCAtMSAtNAo2IDMgMyAwIDQgLTEgMCAtMyAxIDEgNCA1IDMKNyAyIC0zIC0zIDQgMSA0IDIgLTIgLTQgMCA1IC01IDEgLTQKNyA0IDQgMyAtMyAwIDUgNCAxIDEgMSAtMiAyIC0xIDIKOCAxIC0zIDQgNCAtMSAtMSAyIC0xIDEgLTUgMCAtMSAtMyAyIC01IC00CjkgLTIgLTEgLTUgLTMgMSAtNSAyIDMgMyAtMSAtMiAyIC01IC0yIDIgLTEgLTMgLTEKNiAyIDQgMiAzIDUgNCAtNCAtNSAtMyAtMSAtMSAzCjcgNCAtMSAzIC01IDIgMCAwIDUgNCAtMyAtNSAtNSAtMSAzCjkgNSAtNCA1IDMgLTIgLTUgMSAxIDQgNCA1IDUgMiAxIDUgLTIgLTEgMgozIC0xIC01IDEgNCAtMSA1CjkgLTEgLTMgLTMgMiAzIDIgMCAzIC0zIDEgNCAzIC01IC00IC0yIC0xIC00IC00CjIgMCAxIC00IDEKOSAtNSAtNCAtNCAtMiA0IDUgLTQgLTMgLTEgMiAtMyAtMyA0IC0zIDEgLTMgLTQgNAo1IC01IDMgLTQgNSAxIC00IC0xIC01IDQgNAozIDEgNCAtMyAtNSAxIC00CjcgNSA0IDIgMiAwIDUgMCAtNSAtMyAtMSAtMyA0IDUgNQoxMCAtMSAzIDMgNCAtMiAtMSAtNCAzIC0zIC0yIDAgMiAxIC0zIC0zIC01IDUgMCAtNCA0CjIgLTQgLTQgMyA0CjkgLTIgMSAyIDIgMCAtNCAzIC01IDMgMSAtNSAtMyAxIDUgLTIgLTQgLTQgNQo5IC0yIC0zIDQgMSAwIC0yIC0xIDAgNCAwIDEgMSAtMyAwIDUgLTEgNSAxCjcgMyAtNSA0IDQgLTIgLTMgMSAtNCAtNCAtNSAtNSAtMyAtMiAtMgoyIDIgMiA1IDAKMiAxIDIgLTEgNAo4IDAgMiAyIC00IC0yIC0zIDUgLTMgLTQgMCAxIDIgLTMgMyAtMSAtNAo2IC0zIC0xIDUgLTIgLTUgMiAtNSAwIDAgMCAtNSAtNQo5IDIgLTMgLTQgNSAwIC0xIDIgNSAtMiAtMyAtNSAtMiAtNSA0IC0yIDUgLTQgNAo4IDAgNSAtMSAtMyAyIDAgLTEgLTQgMiAtMyAtMiAtMyA1IC0yIC01IDUKOCAtMSAtNSAyIC01IDIgNSAxIC0zIC01IC01IDMgMyA0IDAgLTQgLTQKNSAyIC00IDIgLTUgNSAtMiA1IC01IDIgMQoyIC01IC0xIDEgMgo2IDUgLTUgLTUgLTIgLTMgMyA1IDEgLTIgMyAtMiAtNAo3IC00IC00IDMgLTMgNCAtNCAtMiA0IC01IDIgMyAwIDIgMQo4IDMgNSAtMyAtNSAtMSAyIC0yIDMgLTQgLTMgLTQgLTEgMSAwIC0zIC0zCjEwIC00IDIgLTIgLTQgMSA0IDUgLTIgLTEgLTUgMiAtMiAtNSAtMiAtMSAwIC00IC00IDEgMAo4IDMgLTUgLTEgNSAtMyA0IC01IDIgLTMgNSAtNSAtNSAxIC0xIC00IDEKMTAgLTMgMCAtNCAtMSAtMyAtMyAtMSAyIDUgMCAtMSAtNSAtNCAzIDMgLTUgLTUgLTQgLTUgLTMKNyAxIC0xIDQgLTUgLTIgMyAxIC0zIDAgLTIgMCAwIC01IC0zCjQgLTQgLTEgMyA0IC0zIDUgLTMgMgo4IDEgLTQgMCAwIDIgLTIgMiA0IDEgLTMgMiAtNSA1IC0zIDMgMgo5IDQgLTUgLTQgLTMgMSA1IC01IC0yIC01IC01IDQgLTIgMiAwIDEgMiAtMiAtMgo2IC0yIC0xIDMgMyAwIC0yIDAgMCAtMiA1IC01IDQKOSAxIDAgLTMgLTMgLTIgMiAxIC01IC00IC0zIDEgLTQgLTEgLTIgMCAyIDUgMgoxMCA0IC00IC0yIDEgLTUgLTQgMyA0IDQgNCAtMyAxIC01IDAgNSAtNSA0IC0xIDMgMQoxMCAyIDMgMSAyIDQgMyA1IC01IDMgLTMgMSAwIDQgLTUgNSAzIC0yIC0xIC0xIC00CjIgLTQgMCA0IDQKOCAtMSAtNSAtMSAtMiAtNCAtNSAtNCA1IC0xIDAgNCA0IDMgLTIgMiAzCjMgLTQgLTIgMCAzIDEgMQoyIDUgMSAxIC0zCjEwIC0yIDUgNCAwIDEgMCAtNCAzIDAgLTIgLTQgLTIgLTEgLTEgLTEgNCAtNCAtMyAtNCA1CjggLTMgLTIgNSAtNCA0IDQgNSA1IDMgMSAxIC0yIDEgMiAzIDQKNSAtMSAxIC0xIC00IDQgMCAtMyAxIC01IC0zCjMgLTEgNCAzIDQgNCAtMwo4IC0yIDAgNCAtNCA1IDAgLTQgLTUgMiAtMyAzIC00IDAgLTIgMiA1CjggLTIgNCAyIC00IC01IC0yIDUgLTMgLTQgLTQgMCAtMyAwIDQgLTQgNAo5IDAgLTMgMyA0IDUgLTQgMCAtNCAtNCAtMiAtMiAwIDIgNSAxIDMgMiAxCjIgLTQgNCAwIC0zCjIgMiAxIDAgLTQKOSAtMSAtNSAtNCAyIDAgLTIgMCAxIC0zIC0xIDIgNSAtNCA0IC00IC0xIDEgMAozIDMgMSAtNSAtMiAxIC0yCjYgNSAtNCAtMyAtMSAtNCAzIDIgMSAtMyAzIDEgLTIKMyAyIDAgLTMgLTQgLTUgLTMKNyAtNSAtNSAtMSAzIDIgMSAzIC00IC01IC0yIC0yIDIgLTUgMgoxMCAtNCAtMSA0IDUgMiAzIC01IC01IC0zIDEgMyAxIDUgNCAtNSAtMiAtNSAxIDEgMQo4IC01IC0zIDAgLTMgLTEgLTIgMSAtNCAwIC0yIDEgMSAtMiAtMSAtMyAyCjcgLTEgLTUgNSA1IC00IC0zIDIgLTQgLTQgNSAxIC01IDEgLTMKNSA1IC01IDIgMCA0IDIgLTMgMiAtMiAtNAo3IDQgLTUgLTUgLTMgMSAtNCAtMSAtMyA1IDAgMCAtMiAtMyAtMQoxMCAzIC00IC0xIC00IC00IDUgMSAwIC0xIDAgMiAzIDQgMiAtMiA0IC0zIC0xIC01IDMKMiAxIC0zIC0yIDIKMiAtNSAtNSAtMiAxCjUgMyA0IC0yIC0xIC01IDMgLTIgLTUgMSAwCjIgMyAtMSAyIDIKMyA1IC00IDAgNCAzIC0xCjggMSAtMSAtMyA0IDUgLTQgNSAyIC01IC0yIDMgNCAxIC00IDMgLTEKMTAgMCAxIC0zIDIgMCAtNCAyIC01IDQgNCAtMiA0IC01IC00IC00IDAgLTMgMCAwIC0xCjYgMiAtMyAtNCAtMSAtMSAtNSA1IC0xIDEgLTMgNSAtMgo2IDQgMiAtMyA0IC01IDAgMyAtNCAtNSAyIC01IDQKOSA1IDUgLTUgLTUgNCAxIC00IDEgMyA0IDIgNCA0IC0yIDMgLTQgLTQgNQozIC01IC00IDIgLTEgNSAzCjEwIC00IDUgNSAtMSAtNCAxIC0xIC0xIDIgLTEgLTEgLTUgMCAzIDUgMiAyIDEgLTMgNQo1IC0xIDAgNCAtMiAyIC00IC0zIDAgLTQgNQo4IDIgNCAzIDUgMSAtMyAtMiAtMyAyIC0zIDUgMyAwIC01IC0zIDEKNiA0IDIgMyAzIC0yIDMgLTIgNSAtMyAxIDMgLTEKNyAtMSA0IC00IC0yIC0xIDEgMCAzIC0xIDIgNSA1IC00IC00CjMgMCAtNCAtNCA1IC0xIDQKNCAtNSAzIDEgLTMgNSAwIC0xIC00CjggMCAtNCAtMSAtNCAtMSAwIDEgMSAtNCAtNCAtMyAtNSAtNCAzIDEgLTQKOSAxIC00IC0yIDEgLTUgNCAzIC0xIDUgLTQgNSAzIDAgMCAtNCAyIC0zIC0yCjMgNCAzIC00IC00IDUgNQoxMCAxIDEgLTQgLTEgLTEgMyAyIDEgLTUgLTIgMyAtMyAxIDUgNSAtMyAtMiAwIDMgLTQK"

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

func solveCase(points [][2]int) int64 {
	mX := make(map[int]int)
	mY := make(map[int]int)
	mD1 := make(map[int]int)
	mD2 := make(map[int]int)
	for _, p := range points {
		x, y := p[0], p[1]
		mX[x]++
		mY[y]++
		mD1[x-y]++
		mD2[x+y]++
	}
	var ans int64
	for _, v := range mX {
		ans += int64(v * (v - 1))
	}
	for _, v := range mY {
		ans += int64(v * (v - 1))
	}
	for _, v := range mD1 {
		ans += int64(v * (v - 1))
	}
	for _, v := range mD2 {
		ans += int64(v * (v - 1))
	}
	return ans
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
		n, errN := strconv.Atoi(fields[0])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != 1+2*n {
			fmt.Fprintf(os.Stderr, "line %d expected %d values, got %d\n", idx+1, 1+2*n, len(fields))
			os.Exit(1)
		}
		pts := make([][2]int, n)
		for i := 0; i < n; i++ {
			x, errX := strconv.Atoi(fields[1+2*i])
			y, errY := strconv.Atoi(fields[2+2*i])
			if errX != nil || errY != nil {
				fmt.Fprintf(os.Stderr, "invalid point on line %d\n", idx+1)
				os.Exit(1)
			}
			pts[i] = [2]int{x, y}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for _, p := range pts {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		inputs = append(inputs, sb.String())
		exps = append(exps, solveCase(pts))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
