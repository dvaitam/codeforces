package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(timeStr string, a int) string {
	var h, m int
	fmt.Sscanf(timeStr, "%d:%d", &h, &m)
	total := h*60 + m + a
	total %= 24 * 60
	h = total / 60
	m = total % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	// predefined tests
	times := []string{"00:00", "23:59", "12:34", "01:02"}
	adds := []int{0, 1, 30, 1439}
	tests := make([]struct {
		t string
		a int
	}, 0, 100)
	for i := range times {
		tests = append(tests, struct {
			t string
			a int
		}{times[i], adds[i]})
	}
	for len(tests) < 100 {
		h := rand.Intn(24)
		m := rand.Intn(60)
		a := rand.Intn(10001)
		tests = append(tests, struct {
			t string
			a int
		}{fmt.Sprintf("%02d:%02d", h, m), a})
	}
	for i, tc := range tests {
		inp := fmt.Sprintf("%s\n%d\n", tc.t, tc.a)
		expected := solveB(tc.t, tc.a)
		got, err := runBinary(bin, inp)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: input=%s %d expected %s got %s\n", i+1, tc.t, tc.a, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
