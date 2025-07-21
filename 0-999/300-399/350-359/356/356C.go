package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   cnt := make([]int64, 5)
   var ai int
   var S int64
   for i := int64(0); i < n; i++ {
       fmt.Fscan(in, &ai)
       cnt[ai]++
       S += int64(ai)
   }
   // cost to increase to 3 or 4 from ai
   cost3 := []int64{3, 2, 1, 0, 0}
   cost4 := []int64{4, 3, 2, 1, 0}
   // bounds for number of compartments with 4 students
   // c4 >= max(0, S-3n), c4 <= S/4
   var best int64 = -1
   c4min := int64(0)
   if S > 3*n {
       c4min = S - 3*n
   }
   c4max := S / 4
   if c4max > n {
       c4max = n
   }
   // iterate possible c4
   for c4 := c4min; c4 <= c4max; c4++ {
       remS := S - 4*c4
       if remS < 0 || remS%3 != 0 {
           continue
       }
       c3 := remS / 3
       if c3 < 0 || c4+c3 > n {
           continue
       }
       // assign cost4 greedily from compartments with highest ai
       need4 := c4
       used4 := make([]int64, 5)
       for ai := 4; ai >= 0 && need4 > 0; ai-- {
           take := cnt[ai]
           if take > need4 {
               take = need4
           }
           used4[ai] = take
           need4 -= take
       }
       if need4 > 0 {
           continue
       }
       // assign cost3 from remaining
       need3 := c3
       used3 := make([]int64, 5)
       // order cost3: ai with cost0: 3,4; then 2,1,0
       for _, ai := range []int{3, 4, 2, 1, 0} {
           if need3 <= 0 {
               break
           }
           avail := cnt[ai] - used4[ai]
           if avail <= 0 {
               continue
           }
           take := avail
           if take > need3 {
               take = need3
           }
           used3[ai] = take
           need3 -= take
       }
       if need3 > 0 {
           continue
       }
       // compute cost
       var cost int64
       for ai := 0; ai <= 4; ai++ {
           cost += used4[ai] * cost4[ai]
           cost += used3[ai] * cost3[ai]
       }
       if best < 0 || cost < best {
           best = cost
       }
   }
   fmt.Println(best)
}
