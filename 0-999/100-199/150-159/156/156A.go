package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, u string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &u)
   n := len(s)
   m := len(u)
   best := 0
   // Iterate over all alignments (diagonals) between s and u
   for d := -m + 1; d < n; d++ {
       i := d
       j := 0
       if i < 0 {
           j = -i
           i = 0
       }
       count := 0
       // Traverse diagonal and count matching characters
       for i < n && j < m {
           if s[i] == u[j] {
               count++
           }
           i++
           j++
       }
       if count > best {
           best = count
       }
   }
   // Minimum changes = insertions/deletions to match length + replacements = m - best
   fmt.Println(m - best)
}
