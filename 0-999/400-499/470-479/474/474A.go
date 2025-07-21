package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Read direction of shift
	dirLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	direction := strings.TrimSpace(dirLine)
	// Read typed sequence
	seqLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	typed := strings.TrimSpace(seqLine)

	// Keyboard layout
	const keyboard = "1234567890-=QWERTYUIOP[]\\ASDFGHJKL;'ZXCVBNM,./"
	pos := make(map[rune]int)
	for i, r := range keyboard {
		pos[r] = i
	}

	var result strings.Builder
	for _, r := range typed {
		if idx, ok := pos[r]; ok {
			if direction == "R" {
				// Hands moved right: typed key is to the right of intended
				result.WriteByte(keyboard[idx-1])
			} else {
				// Hands moved left: typed key is to the left of intended
				result.WriteByte(keyboard[idx+1])
			}
		} else {
			// For any other character (e.g., space), output as is
			result.WriteRune(r)
		}
	}

	fmt.Println(result.String())
}
