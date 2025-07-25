package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   freq := make([]int, k+1)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= k {
           freq[x]++
       }
   }
   // total sets to pick
   S := (n + 1) / 2
   // full sets per type: each gives two portions fully used
   totFull := 0
   oddCnt := 0
   for _, cnt := range freq {
       totFull += cnt / 2
       if cnt%2 != 0 {
           oddCnt++
       }
   }
   var ans int
   if totFull >= S {
       ans = 2 * S
   } else {
       // use all full sets, remaining sets yield at most one each for odd counts
       rem := S - totFull
       if rem > oddCnt {
           rem = oddCnt
       }
       ans = 2*totFull + rem
   }
   // cannot exceed n
   if ans > n {
       ans = n
   }
   fmt.Println(ans)
}
