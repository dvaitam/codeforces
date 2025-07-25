package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func season(month string) string {
	switch month {
	case "December", "January", "February":
		return "winter"
	case "March", "April", "May":
		return "spring"
	case "June", "July", "August":
		return "summer"
	default:
		return "autumn"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	for i, m := range months {
		in := fmt.Sprintf("%s\n", m)
		want := season(m)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != want {
			fmt.Printf("test %d failed: input %q expected %q got %q\n", i+1, strings.TrimSpace(in), want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
