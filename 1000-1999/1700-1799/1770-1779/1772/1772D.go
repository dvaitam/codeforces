package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       solve(reader)
   }
}

func solve(reader *bufio.Reader) {
   var n int
   fmt.Fscan(reader, &n)
   v := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }
   const INF = int(1e18)
   mn := INF
   mx := -INF
   for i := 0; i+1 < n; i++ {
       if v[i] < v[i+1] {
           // floor((v[i] + v[i+1]) / 2)
           x := (v[i] + v[i+1]) / 2
           if x < mn {
               mn = x
           }
       }
       if v[i] > v[i+1] {
           // ceil((v[i] + v[i+1]) / 2)
           x := (v[i] + v[i+1] + 1) / 2
           if x > mx {
               mx = x
           }
       }
   }
   // Output according to conditions
   if mx == -INF {
       fmt.Println(0)
   } else if mn == INF {
       // no increasing pair
       fmt.Println(v[0])
   } else {
       if mx <= mn {
           fmt.Println(mx)
       } else {
           fmt.Println(-1)
       }
   }
}
