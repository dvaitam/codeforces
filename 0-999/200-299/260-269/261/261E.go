package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   l, r int64
   p int
   seen map[int64]struct{}
)

func dfs(lastB, remOps int, curProd int64) {
   // try next multiplication: choose t op1 before op2
   // need at least 1 op for multiplication
   for t := 0; t <= remOps-1; t++ {
       // for first multiplication, need at least one increment to have b>=1
       if lastB == 0 && t < 1 {
           continue
       }
       newB := lastB + t
       if newB <= 0 {
           continue
       }
       // cost: t increments + 1 multiplication
       nextRem := remOps - t - 1
       // compute new product
       // check overflow / bound
       if curProd > r/int64(newB) {
           continue
       }
       newProd := curProd * int64(newB)
       if newProd > r {
           continue
       }
       if newProd >= l {
           seen[newProd] = struct{}{}
       }
       if nextRem > 0 {
           dfs(newB, nextRem, newProd)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   _, err := fmt.Fscan(reader, &l, &r, &p)
   if err != nil {
       return
   }
   seen = make(map[int64]struct{})
   // start DFS from initial state: a=1, b=0, remOps = p
   dfs(0, p, 1)
   // count
   cnt := 0
   for x := range seen {
       if x >= l && x <= r {
           cnt++
       }
   }
   // output
   fmt.Println(cnt)
}
