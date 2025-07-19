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
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       var pos [10]int
       for j, ch := range s {
           pos[ch-'0'] = j
       }
       found := false
       for _, p := range primes {
           if pos[p.a] < pos[p.b] {
               fmt.Fprintln(writer, p.val)
               found = true
               break
           }
       }
       if !found {
           fmt.Fprintln(writer, -1)
       }
   }
}

var primes = []struct{ a, b, val int }{
   {1, 3, 13},
   {1, 7, 17},
   {1, 9, 19},
   {2, 3, 23},
   {2, 9, 29},
   {3, 1, 31},
   {3, 7, 37},
   {4, 1, 41},
   {4, 3, 43},
   {4, 7, 47},
   {5, 3, 53},
   {5, 9, 59},
   {6, 1, 61},
   {6, 7, 67},
   {7, 1, 71},
   {7, 3, 73},
   {8, 3, 83},
   {8, 9, 89},
   {9, 7, 97},
}
