package main

import (
   "bufio"
   "fmt"
   "os"
)

var n, m int
var s [][]byte

func work() bool {
   e := [5]int{-1, -1, -1, -1, -1}
   for i := 0; i < n; i++ {
       g := [5]int{-1, -1, -1, -1, -1}
       for j := 0; j < m; j++ {
           if s[i][j] != '0' {
               c := int(s[i][j] - '0')
               pr := (i + 1) & 1
               pc := (j + 1) & 1
               if g[c] != -1 && g[c] != pc {
                   return false
               }
               if e[c] != -1 && e[c] != pr {
                   return false
               }
               g[c] = pc
               e[c] = pr
           }
       }
   }
   cnt0, cnt1 := 0, 0
   for c := 1; c <= 4; c++ {
       if e[c] == 0 {
           cnt0++
       }
       if e[c] == 1 {
           cnt1++
       }
   }
   if cnt0 > 2 || cnt1 > 2 {
       return false
   }
   s1, s2, s3, s4 := 0, 0, 0, 0
   for c := 1; c <= 4; c++ {
       if e[c] != -1 {
           if e[c] == 1 {
               if s1 != 0 {
                   s2 = c
               } else {
                   s1 = c
               }
           } else {
               if s3 != 0 {
                   s4 = c
               } else {
                   s3 = c
               }
           }
       }
   }
   for c := 1; c <= 4; c++ {
       if e[c] == -1 {
           if s1 == 0 {
               s1 = c
           } else if s2 == 0 {
               s2 = c
           } else if s3 == 0 {
               s3 = c
           } else {
               s4 = c
           }
       }
   }
   for i := 0; i < n; i++ {
       // adjust pattern if needed
       for j := 0; j < m; j++ {
           if s[i][j] != '0' {
               c := int(s[i][j] - '0')
               pr := (i + 1) & 1
               pc := (j + 1) & 1
               if pr == 1 && pc == 1 && s1 != c {
                   s1, s2 = s2, s1
               } else if pr == 1 && pc == 0 && s2 != c {
                   s1, s2 = s2, s1
               } else if pr == 0 && pc == 1 && s3 != c {
                   s3, s4 = s4, s3
               } else if pr == 0 && pc == 0 && s4 != c {
                   s3, s4 = s4, s3
               }
               break
           }
       }
       // fill row
       for j := 0; j < m; j++ {
           pr := (i + 1) & 1
           pc := (j + 1) & 1
           var nc int
           if pr == 1 && pc == 1 {
               nc = s1
           } else if pr == 1 && pc == 0 {
               nc = s2
           } else if pr == 0 && pc == 1 {
               nc = s3
           } else {
               nc = s4
           }
           s[i][j] = byte('0' + nc)
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   s = make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       s[i] = []byte(line)
   }
   if work() {
       for i := 0; i < n; i++ {
           writer.Write(s[i])
           writer.WriteByte('\n')
       }
       return
   }
   // transpose
   t := make([][]byte, m)
   for i := 0; i < m; i++ {
       t[i] = make([]byte, n)
       for j := 0; j < n; j++ {
           t[i][j] = s[j][i]
       }
   }
   s = t
   n, m = m, n
   if work() {
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               writer.WriteByte(s[j][i])
           }
           writer.WriteByte('\n')
       }
       return
   }
   writer.WriteString("0\n")
}
