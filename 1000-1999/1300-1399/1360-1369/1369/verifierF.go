package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// memoWin[a][starterTurn]: can starter force a WIN from state a?
// memoLose[a][starterTurn]: can starter force a LOSE from state a?
var memoWin, memoLose map[[2]int64]int8 // 0=unknown,1=true,2=false
var curE int64

func evalMove(m int64, starterTurn bool, forLose bool) bool {
	if m > curE {
		// current player wrote > e, current player loses
		if starterTurn {
			return forLose // starter loses: good for ForceLose (true), bad for ForceWin (false)
		}
		return !forLose // non-starter loses (starter wins): good for ForceWin (!false=true), bad for ForceLose (!true=false)
	}
	if forLose {
		return canForceLose(m, !starterTurn)
	}
	return canForceWin(m, !starterTurn)
}

func canForceWin(a int64, starterTurn bool) bool {
	key := [2]int64{a, 0}
	if !starterTurn {
		key[1] = 1
	}
	if v := memoWin[key]; v != 0 {
		return v == 1
	}
	r1 := evalMove(2*a, starterTurn, false)
	r2 := evalMove(a+1, starterTurn, false)
	var res bool
	if starterTurn {
		res = r1 || r2
	} else {
		res = r1 && r2
	}
	if res {
		memoWin[key] = 1
	} else {
		memoWin[key] = 2
	}
	return res
}

func canForceLose(a int64, starterTurn bool) bool {
	key := [2]int64{a, 0}
	if !starterTurn {
		key[1] = 1
	}
	if v := memoLose[key]; v != 0 {
		return v == 1
	}
	r1 := evalMove(2*a, starterTurn, true)
	r2 := evalMove(a+1, starterTurn, true)
	var res bool
	if starterTurn {
		res = r1 || r2
	} else {
		res = r1 && r2
	}
	if res {
		memoLose[key] = 1
	} else {
		memoLose[key] = 2
	}
	return res
}

func solveF(s, e int64) string {
	curE = e
	memoWin = make(map[[2]int64]int8)
	memoLose = make(map[[2]int64]int8)
	w, l := 0, 0
	if canForceWin(s, true) {
		w = 1
	}
	if canForceLose(s, true) {
		l = 1
	}
	return fmt.Sprintf("%d %d", w, l)
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	const tests = 100
	for t := 0; t < tests; t++ {
		s := int64(rand.Intn(20) + 1)
		e := s + int64(rand.Intn(20))
		input := fmt.Sprintf("1\n%d %d\n", s, e)
		expected := solveF(s, e) + "\n"
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
