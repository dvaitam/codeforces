package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s1, s2 string
   if _, err := fmt.Fscan(reader, &s1, &s2); err != nil {
       return
   }
   // Count zeros and ones in s1
   z1, o1 := 0, 0
   for i := 0; i < len(s1); i++ {
       if s1[i] == '0' {
           z1++
       } else {
           o1++
       }
   }
   m := len(s2)
   // Build KMP prefix function for s2
   pi := make([]int, m)
   for i := 1; i < m; i++ {
       j := pi[i-1]
       for j > 0 && s2[i] != s2[j] {
           j = pi[j-1]
       }
       if s2[i] == s2[j] {
           j++
       }
       pi[i] = j
   }
   // Count zeros and ones in s2
   z2, o2 := 0, 0
   for i := 0; i < m; i++ {
       if s2[i] == '0' {
           z2++
       } else {
           o2++
       }
   }
   // If not enough chars to use s2 even once, output s1
   if z1 < z2 || o1 < o2 {
       writer := bufio.NewWriter(os.Stdout)
       writer.WriteString(s1)
       writer.Flush()
       return
   }
   // Prepare output builder
   var builder strings.Builder
   builder.Grow(len(s1))
   // Append first occurrence of s2
   builder.WriteString(s2)
   z1 -= z2
   o1 -= o2
   // Overlap length from KMP
   overlap := pi[m-1]
   // Count zeros and ones in the overlapped prefix
   zo, oo := 0, 0
   for i := 0; i < overlap; i++ {
       if s2[i] == '0' {
           zo++
       } else {
           oo++
       }
   }
   // Counts for each additional segment
   segZ, segO := z2-zo, o2-oo
   // Append as many overlapped segments as possible
   for z1 >= segZ && o1 >= segO {
       builder.WriteString(s2[overlap:])
       z1 -= segZ
       o1 -= segO
   }
   // Append remaining zeros then ones
   for z1 > 0 {
       builder.WriteByte('0')
       z1--
   }
   for o1 > 0 {
       builder.WriteByte('1')
       o1--
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   writer.WriteString(builder.String())
   writer.Flush()
}
