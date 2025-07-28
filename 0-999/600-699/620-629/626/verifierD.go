package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Test struct {
	N   int    `json:"n"`
	Arr []int  `json:"arr"`
	Out string `json:"out"`
}

func runCase(bin string, t Test) (string, error) {
	arrStr := strings.TrimSpace(strings.Trim(fmt.Sprint(t.Arr), "[]"))
	input := fmt.Sprintf("%d\n%s\n", t.N, arrStr)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testsD.json")
	if err != nil {
		fmt.Println("cannot read testsD.json:", err)
		return
	}
	var tests []Test
	if err := json.Unmarshal(data, &tests); err != nil {
		fmt.Println("bad tests json:", err)
		return
	}
	passed := 0
	for i, t := range tests {
		got, err := runCase(bin, t)
		got = strings.TrimSpace(got)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			continue
		}
		expectedVal, err1 := strconv.ParseFloat(strings.TrimSpace(t.Out), 64)
		gotVal, err2 := strconv.ParseFloat(got, 64)
		if err1 == nil && err2 == nil {
			diff := math.Abs(expectedVal - gotVal)
			if diff <= 1e-6*math.Max(1.0, math.Abs(expectedVal)) {
				passed++
				continue
			}
		} else if got == strings.TrimSpace(t.Out) {
			passed++
			continue
		}
		fmt.Printf("case %d failed: expected %s got %s\n", i+1, t.Out, got)
	}
	fmt.Printf("passed %d/%d\n", passed, len(tests))
}
