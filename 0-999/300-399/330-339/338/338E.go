package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var h int64
   if _, err := fmt.Fscan(reader, &n, &m, &h); err != nil {
       return
   }
   b := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Special case: pattern length 1
   if m == 1 {
       var cnt int
       for i := 0; i < n; i++ {
           if abs64(a[i]-b[0]) <= h {
               cnt++
           }
       }
       fmt.Println(cnt)
       return
   }
   // If pattern longer than text
   if n < m {
       fmt.Println(0)
       return
   }
   // Build difference arrays
   diffB := make([]int64, m-1)
   for i := 0; i < m-1; i++ {
       diffB[i] = b[i+1] - b[i]
   }
   diffA := make([]int64, n-1)
   for i := 0; i < n-1; i++ {
       diffA[i] = a[i+1] - a[i]
   }
   // Prefix-function for KMP
   pi := make([]int, m-1)
   for i := 1; i < len(diffB); i++ {
       j := pi[i-1]
       for j > 0 && diffB[i] != diffB[j] {
           j = pi[j-1]
       }
       if diffB[i] == diffB[j] {
           j++
       }
       pi[i] = j
   }
   // KMP search
   var cnt int
   j := 0
   for i := 0; i < len(diffA); i++ {
       for j > 0 && diffA[i] != diffB[j] {
           j = pi[j-1]
       }
       if diffA[i] == diffB[j] {
           j++
       }
       if j == len(diffB) {
           // match at diffA[i - (m-2) ... i]
           start := i - (m-2)
           if start >= 0 && int(start) < len(a) {
               if abs64(a[start]-b[0]) <= h {
                   cnt++
               }
           }
           j = pi[j-1]
       }
   }
   fmt.Println(cnt)
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
