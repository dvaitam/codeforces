package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscanf(reader, "%d\n", &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscanf(reader, "%s", &s); err != nil {
       return
   }
   // number of points must be even
   if n%2 != 0 {
       fmt.Println(0)
       return
   }
   const mask = 0xFFFFFFFF
   var ways uint64 = 1
   // for each pair (i, i+1): i is small letter, i+1 is uppercase
   for i := 0; i < n; i += 2 {
       c1 := s[i]
       c2 := s[i+1]
       // uppercase letters were erased, so position i+1 must be '?' now
       if c2 != '?' {
           fmt.Println(0)
           return
       }
       // choose small letter at position i if unknown
       if c1 == '?' {
           ways = (ways * 25) & mask
       }
   }
   fmt.Println(ways)
