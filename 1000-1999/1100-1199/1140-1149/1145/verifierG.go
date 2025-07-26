package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		var cmd *exec.Cmd
		if strings.HasSuffix(bin, ".go") {
			cmd = exec.Command("go", "run", bin)
		} else {
			cmd = exec.Command(bin)
		}
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Printf("Run %d: runtime error: %v\n", i+1, err)
			fmt.Printf("Output: %s\n", result)
			os.Exit(1)
		}
		if result != "" {
			fmt.Printf("Run %d failed: expected no output got %s\n", i+1, result)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
