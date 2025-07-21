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
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       var results []string
       // Try all possible a such that a*b = 12
       for a := 1; a <= 12; a++ {
           if 12%a != 0 {
               continue
           }
           b := 12 / a
           found := false
           // Check each column j
           for j := 0; j < b; j++ {
               allX := true
               for i := 0; i < a; i++ {
                   idx := i*b + j
                   if s[idx] != 'X' {
                       allX = false
                       break
                   }
               }
               if allX {
                   found = true
                   break
               }
           }
           if found {
               results = append(results, fmt.Sprintf("%dx%d", a, b))
           }
       }
       // Output
       fmt.Fprint(writer, len(results))
       for _, pair := range results {
           fmt.Fprint(writer, " ", pair)
       }
       fmt.Fprintln(writer)
   }
}
