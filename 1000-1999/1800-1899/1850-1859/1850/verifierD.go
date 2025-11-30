package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesB64 = "MTUgMCA1IDUgMjMgMTAgNDcgNDIgMTkgMTYgMzggMTMgMzggMiAzNyA0MyAxMAo4IDUgMjUgNDYgMzIgMjMgMzQgMjggMzIgMTcKMiAwIDIzIDI5CjcgMyAyNyAzMyAxMCAzNSAxMSAxNSAxNAoyIDEgMjAgMTEKNCA0IDMyIDIzIDMyIDQzCjEwIDEgMjggNTAgMjYgNDcgMzMgNDggMjMgNTAgMzcgMjIKNyAzIDEwIDQ4IDI1IDQ1IDQ3IDI5IDQxCjEwIDEgMzEgMTcgMzEgMzIgMzIgNTAgMjIgNDIgMjkgMjkKNyA0IDQ2IDM1IDQ2IDI5IDMxIDQyIDE0CjcgNSAxMCAzOSAxNyA0OSAzMCAxOSAxOQoxNCA1IDMyIDM1IDMzIDMyIDQxIDM5IDM3IDI2IDE5IDQ2IDEzIDMxIDMyIDIzCjEyIDQgNCA1MCAyMSA0NiAwIDEyIDQ3IDYgMyAzNiA0MSAzCjYgNCAxNCA0MyA2IDQ4IDMzIDgKMTUgMiAxNSAxMyAzIDI3IDQ1IDQ4IDIgMyAyMyAyMyAxMSAxNSA0MyAxIDUKMyAwIDEgMiA0NgoyIDIgMTYgOAoxNSAxIDQ3IDExIDMzIDQ0IDAgMjQgMzcgMiA1MCAxNSA5IDIgMCAyMiAzOQoxMiA1IDQ3IDcgMTggMjEgMzEgMSAxOSAyOCAzNSA0OSAzOCA0NwoyIDIgNDggMjUKMTUgNCA0NSA5IDMwIDE0IDUgNDIgNDMgMjAgNiAxIDI4IDUwIDggMzMgMzcKMTQgMyAzMSAzMiAyMCA5IDIxIDE2IDE2IDM4IDI2IDQxIDEgNDQgMzUgOAoxMiAwIDE2IDIgOCAxMCAxMCA2IDI5IDQwIDE0IDMyIDQ1IDIKNSAxIDQ1IDI4IDQgMTYgNQoxMSAxIDM5IDUwIDM5IDQ1IDIzIDE2IDQzIDI3IDE3IDMzIDQ4CjIgMSAyIDI0CjggMSA3IDMyIDQ2IDUgMTUgNiA2IDEKNCAxIDYgMTMgMSAzMwoxMiAzIDI5IDE5IDM0IDQxIDI0IDEzIDQzIDQ4IDEzIDQ2IDI3IDI3CjEwIDAgMzcgMzcgMyAyNiAzMyAzNyAxMSA2IDQyIDMwCjcgMCAzMyA3IDM5IDIzIDE4IDQ0IDIzCjYgMCA0MyAyNiA2IDYgMTkgMTIKMTUgNSAxIDI4IDMgMjYgNDAgMzEgMjkgMTMgMzcgMzkgNCAwIDE4IDEgMjMKNiA1IDQgMTQgNDggMzEgMTIgNwoxMSAyIDI1IDQ1IDI5IDggNDggMjIgMjUgNyAxNiA3IDcKMyA0IDIxIDQxIDI1CjUgNSA2IDEgMzkgNDIgMzAKMTQgMCA0NiA0NSAzMSAxOCAyMiAyOSA5IDIzIDE3IDMwIDMzIDMwIDQ2IDQ2CjE0IDMgMzEgNDMgMTggMjUgMTQgMTAgMzEgMzggMTYgMzUgMjcgNDQgNDMgNDQKMyA0IDQ2IDM2IDYKMyAyIDExIDM0IDkKMTQgMyA0IDUgNDMgNDEgMiA4IDE4IDI0IDE0IDQ1IDQyIDQzIDIxIDI4CjQgNCAxOCA3IDkgMzQKMTQgMyA2IDIxIDMzIDE1IDQ1IDMyIDE2IDEwIDEwIDI5IDQ1IDE1IDI1IDIyCjE0IDQgNDYgOSAyOSAyOCA0NiAxIDM4IDI0IDQ3IDExIDI1IDMyIDMgMzAKNiAzIDE2IDQ1IDQ2IDI2IDQ1IDQxCjkgMiAzNSAyMSA0NSA0NyA0MiA1IDQ4IDQ2IDE0CjEwIDQgMTIgMjUgNDIgMjQgNDAgMCAyMCAyOSAzMyA0NQo5IDUgMTEgNiAxIDI1IDEzIDQ2IDM2IDM4IDI0CjUgMCAyNCAzNSA0OSAxMiAxNwoxMyA0IDM3IDEyIDMxIDM5IDggMCAzOSA0MyAyNyAzMCAxNiAzMiAzNgo0IDMgNDUgMTMgNDggNAo3IDAgMzEgMzQgNDIgNDIgNCA0OCAzNwo5IDUgMjEgMjkgMTcgMzIgMjkgMSA1IDM5IDQ4CjcgMSA0OCA0OCA1MCAyNSAxNiA0MyA0MAoxNCA1IDggMyAxMCAzMSAyNCAyOSA0MyAxOCA5IDAgMTggMzUgMjkgMAo3IDAgMzQgMjQgMzYgMjggMTMgNDMgMTkKOSA1IDggMzAgNDQgMzQgNDUgMTkgNCAxNiAyMAo2IDIgNDEgNTAgMTkgNDEgNDEgMjUKMTAgMCAzMiA0MCAxMyAyNSAzOCAzMyA5IDMyIDQwIDUKNiAwIDE0IDI5IDM1IDE0IDMzIDE3CjIgMCA3IDQzCjE1IDMgMjMgMTMgMjAgMjIgNCAyMSAyOSAyMyAxMCAzMSAyOCAxOCAyOSA4IDQ1CjkgNSAxMyAxNyAyMCAxMCA2IDE1IDMwIDEyIDQ4CjEyIDIgMTEgMjIgOCA1MCA4IDE0IDE3IDM1IDQwIDI0IDI1IDQ3CjcgMiA0NiAzOCAzMiAzNyA0NCA0NiAyMAoxMyAzIDQ4IDQ1IDQ1IDQwIDQ4IDQ1IDE4IDM0IDM5IDQwIDQyIDQgMjMKNiAzIDMwIDExIDE2IDIyIDI4IDMwCjMgMSAyMCAyNCA4CjIgMCAyMiAxMAo3IDAgNDYgNDggNDEgMjcgMCAzNCAyMAo1IDQgMjQgMzQgMTggMzAgNDAKNCAyIDIwIDEyIDMxIDYKNCAxIDIxIDE2IDkgMjYKNyAyIDUgMjEgMTIgMTUgNDUgMTUgNDYKMTEgMCAyMSAyMyA0MSA0OSAzOSAzIDkgMTEgNCAyNyAyOAoxNCAyIDggMjAgMzMgMzYgNyAyMSA0MSA0OSA0NSAzOSAyNSAxNCAzIDI1CjE0IDMgMzEgMzkgMjAgMzQgMzkgMzggNSAzNyAzMiAzNCA0MiAzMSAyNSA0NAo5IDEgMjYgMjQgMzMgMjggMiA2IDI4IDM3IDgKMyA1IDMyIDExIDQKOCAyIDI5IDQ1IDAgMTYgNiA0MiAyMiAxNAo0IDAgOSAyNyA0MiA1CjcgNSAyOSAzIDMwIDE1IDQgMzAgOAoxMCAwIDggNDQgMzIgMzQgMyAzIDEyIDM0IDAgMzMKNyA1IDMzIDE1IDggMjMgMzEgMCA4CjEwIDAgMTUgNiAyOSAxMyAzIDM5IDEzIDQwIDI0IDIxCjExIDUgMjUgNDUgMzMgMzIgNDkgNDMgMTAgMzIgNiA5IDQwCjUgMSAyNCAxMiAxOSAyMSAyNwo0IDMgOCAyNSAyMCAxOQoxNCAwIDM1IDYgMzAgMTcgMTggMzMgNDggMzEgMTcgMTQgMjYgNDQgOCA0NAoxMCA1IDYgMSAzOCAzNSA0OCAxMiAxMyAxMiAyNSAzNwoyIDUgOCA0MAoyIDUgMTYgNDQKMTMgMyAzNCAzIDQ3IDQ5IDE0IDkgMzggMjAgMiA0NCAxMiA2IDgKMTIgNSAzNCAxMSA0OCA1IDQzIDI5IDQwIDE4IDEzIDEwIDIwIDQ0CjE1IDIgMzMgMzYgNCAyNiAyNiA0MiA0NiAyIDI5IDE5IDQyIDcgNDAgNDQgNDYKNiAwIDEzIDI2IDIxIDE2IDM0IDQ2CjggNCAzMyA0NyAxMiAyNyA0OSA4IDQ0IDEwCjE0IDMgMjkgMjIgMjQgMzAgMzkgMTYgMzkgMTIgMzcgMzAgMjggMTIgNDggMzAKMTUgNCAyMSAxOSA0IDEwIDIzIDM4IDQwIDMwIDE0IDQ5IDM5IDQxIDQyIDM2IDgKMTIgMiAxMyAzNCAxOSA2IDAgNTAgMSAxMiAyMCAzIDIwIDM0CjYgNSA0MiAyMSAyOCA0IDI2IDMwCjE1IDUgMSAxOCAzNyAzNiA4IDEzIDkgMTAgMzggNDkgMjQgNDYgNCA0MCAzNwo5IDIgNDEgNSAzMSAzMCA1MCAxNSA5IDM2IDE5CjE1IDEgMTIgMzkgNDUgMjEgMzcgMzkgNDQgMjUgMzMgMjYgMTUgNDEgMTMgMzUgMwo2IDUgMTUgOCAzOSA0NiAyNSAyNwozIDMgMjUgMjUgMzAKOCAyIDEzIDE1IDE0IDMgMzQgMzMgNSAzOAoxMCA1IDAgMyAyNCA0NSAyNyAyNSAxNCAzMiAxNyA2CjcgNCAyMyAzMyA1MCAzMSAzNyA0IDQ1Cg=="

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

func solveCase(n, k int, arr []int) int {
	sort.Ints(arr)
	maxLen, cur := 1, 1
	for i := 1; i < n; i++ {
		if arr[i]-arr[i-1] <= k {
			cur++
		} else {
			if cur > maxLen {
				maxLen = cur
			}
			cur = 1
		}
	}
	if cur > maxLen {
		maxLen = cur
	}
	return n - maxLen
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
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "invalid line %d\n", idx+1)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(fields[0])
		k, errK := strconv.Atoi(fields[1])
		if errN != nil || errK != nil {
			fmt.Fprintf(os.Stderr, "invalid header on line %d\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != n+2 {
			fmt.Fprintf(os.Stderr, "line %d expected %d numbers, got %d\n", idx+1, n+2, len(fields))
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value on line %d\n", idx+1)
				os.Exit(1)
			}
			arr[i] = val
		}
		input := fmt.Sprintf("1\n%d %d\n", n, k)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		inputs = append(inputs, input)
		exps = append(exps, solveCase(n, k, arr))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		var got int
		fmt.Fscan(strings.NewReader(out), &got)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %d\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
