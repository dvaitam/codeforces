package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   cards := make([][]int, n)
   for i := 0; i < n; i++ {
       var m int
       fmt.Fscan(in, &m)
       cards[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &cards[i][j])
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < n; i++ {
       // Check if any other player's card is subset of player i's card
       has := make([]bool, 101)
       for _, x := range cards[i] {
           if x >= 1 && x <= 100 {
               has[x] = true
           }
       }
       win := true
       for j := 0; j < n; j++ {
           if j == i {
               continue
           }
           subset := true
           for _, x := range cards[j] {
               if x < 1 || x > 100 || !has[x] {
                   subset = false
                   break
               }
           }
           if subset {
               win = false
               break
           }
       }
       if win {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
