package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   lens := make([]int, n+1)
   myMap := make(map[string]int)
   for i := 1; i <= n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       if myMap[s] == 0 {
           myMap[s] = i
           a[i] = i
           lens[i] = len(s)
       } else {
           a[i] = myMap[s]
       }
   }
   sum := n - 1
   for i := 1; i <= n; i++ {
       sum += lens[a[i]]
   }
   ans := 0
   for i := 1; i <= n; i++ {
       for j := i; j <= n; j++ {
           cnt := 1
           x := -1
           for k := i; k <= j; k++ {
               x += lens[a[k]]
           }
           poi := j + 1
           length := j - i + 1
           for poi+length-1 <= n {
               ok := true
               for k := 0; k < length; k++ {
                   if a[i+k] != a[poi+k] {
                       ok = false
                       break
                   }
               }
               if ok {
                   cnt++
                   poi += length
               } else {
                   poi++
               }
           }
           if cnt > 1 {
               if v := x * cnt; v > ans {
                   ans = v
               }
           }
       }
   }
   fmt.Println(sum - ans)
}
