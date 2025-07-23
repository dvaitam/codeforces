package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       return
   }
   // Remove trailing newline/carriage return
   s = strings.TrimRight(s, "\r\n")
   // Reverse and output
   for i := len(s) - 1; i >= 0; i-- {
       fmt.Printf("%c", s[i])
   }
   fmt.Println()
}
