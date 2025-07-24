package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   counts := make(map[int]int)
   for i := 0; i < t; i++ {
       var op string
       var s string
       fmt.Fscan(reader, &op, &s)
       switch op {
       case "+":
           key := numKey(s)
           counts[key]++
       case "-":
           key := numKey(s)
           counts[key]--
       case "?":
           key := patKey(s)
           fmt.Fprintln(writer, counts[key])
       }
   }
}

// numKey computes the parity pattern key for a decimal number string.
func numKey(s string) int {
   key := 0
   // least significant digit at position 0
   for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
       if (s[j]-'0')&1 == 1 {
           key |= 1 << i
       }
   }
   return key
}

// patKey computes the key from the pattern string of '0' and '1'.
func patKey(s string) int {
   key := 0
   for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
       if s[j] == '1' {
           key |= 1 << i
       }
   }
   return key
}
