package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // remove trailing newline if any
   if len(s) > 0 && s[len(s)-1] == '\n' {
       s = s[:len(s)-1]
   }
   n := len(s) / 2
   first := []rune(s[:n])
   second := []rune(s[n:])
   reverse(first)
   reverse(second)
   fmt.Println(string(first) + string(second))
}

// reverse reverses a slice of runes in place
func reverse(r []rune) {
   for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
       r[i], r[j] = r[j], r[i]
   }
}
