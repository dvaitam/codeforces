package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   sum := 0
   idxMod := make([][]int, 3)
   for i := 0; i < n; i++ {
       d := int(s[i] - '0')
       sum += d
       idxMod[d%3] = append(idxMod[d%3], i)
   }
   p := sum % 3
   remove := make([]bool, n)
   switch p {
   case 1:
       if len(idxMod[1]) >= 1 {
           pos := idxMod[1][len(idxMod[1])-1]
           remove[pos] = true
       } else if len(idxMod[2]) >= 2 {
           l := len(idxMod[2])
           remove[idxMod[2][l-1]] = true
           remove[idxMod[2][l-2]] = true
       } else {
           fmt.Println(-1)
           return
       }
   case 2:
       if len(idxMod[2]) >= 1 {
           pos := idxMod[2][len(idxMod[2])-1]
           remove[pos] = true
       } else if len(idxMod[1]) >= 2 {
           l := len(idxMod[1])
           remove[idxMod[1][l-1]] = true
           remove[idxMod[1][l-2]] = true
       } else {
           fmt.Println(-1)
           return
       }
   }
   // build result
   var sb strings.Builder
   for i := 0; i < n; i++ {
       if !remove[i] {
           sb.WriteByte(s[i])
       }
   }
   t := sb.String()
   if len(t) == 0 {
       fmt.Println(-1)
       return
   }
   // strip leading zeros
   i := 0
   for i < len(t) && t[i] == '0' {
       i++
   }
   if i == len(t) {
       // all zeros
       fmt.Println("0")
       return
   }
   fmt.Println(t[i:])
}
