package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read input string
   line, _ := reader.ReadString('\n')
   s := strings.TrimSpace(line)
   n := len(s)
   seen := make(map[string]bool)
   // Generate all cyclic shifts by moving last i chars to front
   for i := 0; i < n; i++ {
       // shift by i: take last i chars to front
       if i == 0 {
           seen[s] = true
           continue
       }
       // last i chars: s[n-i:]
       t := s[n-i:] + s[:n-i]
       seen[t] = true
   }
   // Output count of distinct strings
   fmt.Println(len(seen))
}
