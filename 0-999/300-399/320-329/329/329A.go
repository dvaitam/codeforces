package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       a := make([][]bool, n)
       for i := 0; i < n; i++ {
           var s string
           fmt.Fscan(reader, &s)
           row := make([]bool, n)
           for j := 0; j < n; j++ {
               row[j] = (s[j] != '.')
           }
           a[i] = row
       }
       badr, badc := false, false
       for i := 0; i < n; i++ {
           cnt := 0
           for j := 0; j < n; j++ {
               if a[i][j] {
                   cnt++
               }
           }
           if cnt == n {
               badr = true
           }
       }
       for i := 0; i < n; i++ {
           cnt := 0
           for j := 0; j < n; j++ {
               if a[j][i] {
                   cnt++
               }
           }
           if cnt == n {
               badc = true
           }
       }
       if badr && badc {
           fmt.Println(-1)
           continue
       }
       if badc {
           for i := 0; i < n; i++ {
               for j := 0; j < n; j++ {
                   if !a[i][j] {
                       fmt.Println(i+1, j+1)
                       break
                   }
               }
           }
       } else {
           for i := 0; i < n; i++ {
               for j := 0; j < n; j++ {
                   if !a[j][i] {
                       fmt.Println(j+1, i+1)
                       break
                   }
               }
           }
       }
   }
}
