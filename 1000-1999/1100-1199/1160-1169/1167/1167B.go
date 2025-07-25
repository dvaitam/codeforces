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

   special := []int{4, 8, 15, 16, 23, 42}
   prods := make(map[int]int)
   // Query products a1 * ai for i = 2..5
   for i := 2; i <= 5; i++ {
       fmt.Fprintf(writer, "? 1 %d\n", i)
       writer.Flush()
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       prods[i] = x
   }

   var ans [7]int
   // Try each possible a1
   for _, a1 := range special {
       used := make(map[int]bool)
       ans[1] = a1
       used[a1] = true
       ok := true
       for i := 2; i <= 5; i++ {
           p := prods[i]
           if p%a1 != 0 {
               ok = false
               break
           }
           ai := p / a1
           // check ai is one of special and not used
           valid := false
           for _, v := range special {
               if v == ai {
                   valid = true
                   break
               }
           }
           if !valid || used[ai] {
               ok = false
               break
           }
           ans[i] = ai
           used[ai] = true
       }
       if !ok {
           continue
       }
       // Deduce a6 as the remaining special number
       for _, v := range special {
           if !used[v] {
               ans[6] = v
           }
       }
       break
   }

   // Output the guessed permutation
   fmt.Fprintf(writer, "! %d %d %d %d %d %d\n", ans[1], ans[2], ans[3], ans[4], ans[5], ans[6])
   writer.Flush()
}
