package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []int{
	990085, 584054, 13706, 5143, 389994, 299330, 872337, 872618, 225038, 57062,
	978309, 326955, 688419, 711381, 155253, 133621, 70643, 752812, 163996, 618060,
	438342, 281471, 802465, 849724, 137088, 76496, 944571, 202729, 967038, 801531,
	168470, 631275, 436994, 227115, 819471, 411650, 567354, 889506, 514036, 181528,
	908447, 187494, 199513, 921436, 856910, 691451, 167533, 910938, 495302, 65880,
	524576, 195396, 265880, 38917, 407593, 754130, 516235, 781849, 63710, 90187,
	248336, 499516, 15126, 28723, 659648, 995663, 51914, 132955, 99248, 28404,
	129089, 253826, 931925, 957724, 131256, 297559, 396342, 995049, 895486, 301274,
	950126, 135342, 926885, 3014, 365784, 190212, 404578, 955689, 628816, 218830,
	247258, 745819, 349290, 384048, 38786, 536168, 750302, 467880, 8803, 97086,
}

const nodes = 95

func solve(k int) string {
	edge := make([][]bool, nodes)
	for i := range edge {
		edge[i] = make([]bool, nodes)
	}
	draw := func(a, b int) {
		edge[a][b] = true
		edge[b][a] = true
	}
	for i := 0; i < 30; i++ {
		for j := 2*i + 2; j < 2*i+4; j++ {
			for l := 2*i + 4; l < 2*i+6; l++ {
				draw(j, l)
			}
		}
	}
	for i := 64; i < 94; i++ {
		draw(i, i+1)
	}
	draw(0, 2)
	draw(0, 3)
	draw(94, 1)
	for i := 0; i < 30; i++ {
		if k&(1<<i) != 0 {
			draw(2*i+2, 64+i+1)
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(nodes))
	sb.WriteByte('\n')
	for i := 0; i < nodes; i++ {
		for j := 0; j < nodes; j++ {
			if edge[i][j] {
				sb.WriteByte('Y')
			} else {
				sb.WriteByte('N')
			}
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, k := range rawTestcases {
		expected := solve(k)
		input := fmt.Sprintf("%d\n", k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
