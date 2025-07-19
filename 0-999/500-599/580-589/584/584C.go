package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// note returns a lowercase letter different from x and y
func note(x, y byte) byte {
   for z := byte('a'); z <= 'z'; z++ {
       if z != x && z != y {
           return z
       }
   }
   return 'a'
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, T int
   if _, err := fmt.Fscan(reader, &N, &T); err != nil {
       return
   }
   s1 := make([]byte, N)
   s2 := make([]byte, N)
   var str1, str2 string
   fmt.Fscan(reader, &str1, &str2)
   copy(s1, str1)
   copy(s2, str2)

   eq, dif := 0, 0
   for i := 0; i < N; i++ {
       if s1[i] == s2[i] {
           eq++
       } else {
           dif++
       }
   }
   // determine number of equal positions to change
   // sd in [max(0, T-dif) .. min(eq, T-ceil(dif/2))]
   sdMin := max(0, T-dif)
   // ceil(dif/2)
   half := (dif + 1) / 2
   sdMax := min(eq, T-half)
   if sdMin > sdMax {
       fmt.Fprintln(writer, -1)
       return
   }
   sd := sdMin
   Tprime := T - sd
   // number of mismatches to replace with new letters
   db := 2*Tprime - dif
   ans := make([]byte, N)

   // handle mismatched positions
   c := 0
   for i := 0; i < N; i++ {
       if s1[i] != s2[i] {
           if c < db {
               ans[i] = note(s1[i], s2[i])
           } else {
               if c%2 == 1 {
                   ans[i] = s1[i]
               } else {
                   ans[i] = s2[i]
               }
           }
           c++
       }
   }
   // handle equal positions
   c = 0
   for i := 0; i < N; i++ {
       if s1[i] == s2[i] {
           if c < sd {
               // change to a different letter
               // next letter cyclically
               ans[i] = byte('a' + (s1[i]-'a'+1)%26)
           } else {
               ans[i] = s1[i]
           }
           c++
       }
   }
   // output
   writer.Write(ans)
   writer.WriteByte('\n')
