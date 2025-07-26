package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       os.Exit(1)
   }
   // Remove trailing newline/carriage return
   s := strings.TrimRight(input, "\r\n")
   // Reverse runes to handle Unicode correctly
   runes := []rune(s)
   for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
       runes[i], runes[j] = runes[j], runes[i]
   }
   // Print reversed string
   fmt.Println(string(runes))
}
