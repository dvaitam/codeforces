package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   stars := make([][2]int, 0)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j := 0; j < m; j++ {
           if s[j] == '*' {
               stars = append(stars, [2]int{i, j})
           }
       }
   }
   ansI, ansJ := -1, -1
   for i := 0; i < n && ansI == -1; i++ {
       for j := 0; j < m && ansI == -1; j++ {
           ok := true
           for _, p := range stars {
               if p[0] != i && p[1] != j {
                   ok = false
                   break
               }
           }
           if ok {
               ansI, ansJ = i, j
           }
       }
   }
   if ansI != -1 {
       fmt.Println("YES")
       fmt.Println(ansI+1, ansJ+1)
   } else {
       fmt.Println("NO")
   }
}
