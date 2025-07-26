package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

// reverse returns the reversed string of s, handling Unicode correctly.
func reverse(s string) string {
   r := []rune(s)
   for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
       r[i], r[j] = r[j], r[i]
   }
   return string(r)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       fmt.Fprintln(os.Stderr, err)
       os.Exit(1)
   }
   // Trim newline and carriage return
   s = strings.TrimRight(s, "\r\n")
   if s == reverse(s) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
