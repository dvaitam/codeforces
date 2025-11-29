package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1191ASource = `package main

import "fmt"

func main() {
    var x int
    if _, err := fmt.Scan(&x); err != nil {
        return
    }
    bestRank := -1
    bestA := 0
    var bestCat byte
    for a := 0; a <= 2; a++ {
        r := (x + a) % 4
        rank := 0
        var cat byte
        switch r {
        case 1:
            rank = 4
            cat = 'A'
        case 3:
            rank = 3
            cat = 'B'
        case 2:
            rank = 2
            cat = 'C'
        case 0:
            rank = 1
            cat = 'D'
        }
        if rank > bestRank {
            bestRank = rank
            bestA = a
            bestCat = cat
        }
    }
    fmt.Printf("%d %c\n", bestA, bestCat)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1191ASource

var testcases = []int{
	30,
	31,
	32,
	33,
	34,
	35,
	36,
	37,
	38,
	39,
	40,
	41,
	42,
	43,
	44,
	45,
	46,
	47,
	48,
	49,
	50,
	51,
	52,
	53,
	54,
	55,
	56,
	57,
	58,
	59,
	60,
	61,
	62,
	63,
	64,
	65,
	66,
	67,
	68,
	69,
	70,
	71,
	72,
	73,
	74,
	75,
	76,
	77,
	78,
	79,
	80,
	81,
	82,
	83,
	84,
	85,
	86,
	87,
	88,
	89,
	90,
	91,
	92,
	93,
	94,
	95,
	96,
	97,
	98,
	99,
	100,
	30,
	31,
	32,
	33,
	34,
	35,
	36,
	37,
	38,
	39,
	40,
	41,
	42,
	43,
	44,
	45,
	46,
	47,
	48,
	49,
	50,
	51,
	52,
	53,
	54,
	55,
	56,
	57,
	58,
}

func solveCase(x int) string {
	bestRank := -1
	bestA := 0
	var bestCat byte
	for a := 0; a <= 2; a++ {
		r := (x + a) % 4
		rank := 0
		var cat byte
		switch r {
		case 1:
			rank = 4
			cat = 'A'
		case 3:
			rank = 3
			cat = 'B'
		case 2:
			rank = 2
			cat = 'C'
		case 0:
			rank = 1
			cat = 'D'
		}
		if rank > bestRank {
			bestRank = rank
			bestA = a
			bestCat = cat
		}
	}
	return fmt.Sprintf("%d %c", bestA, bestCat)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, x := range testcases {
		input := fmt.Sprintf("%d\n", x)
		expect := solveCase(x)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
