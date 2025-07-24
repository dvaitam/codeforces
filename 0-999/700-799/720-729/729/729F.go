package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n int
   a []int64
   pref []int64
   memo map[state]int64
)

type state struct {
   l, r, k int
   turn    int // 0 for Igor, 1 for Zhenya
}

func getSum(l, r int) int64 {
   if l > r {
       return 0
   }
   return pref[r+1] - pref[l]
}

func solve(l, r, k, turn int) int64 {
   if l > r {
       return 0
   }
   st := state{l, r, k, turn}
   if v, ok := memo[st]; ok {
       return v
   }
   var res int64
   rem := r - l + 1
   moves := []int{ }
   if k == 0 {
       moves = []int{1, 2}
   } else {
       moves = []int{k, k + 1}
   }
   if turn == 0 {
       // Igor: maximize
       res = -(1 << 60)
       for _, x := range moves {
           if x <= rem {
               sum := getSum(l, l+x-1)
               val := sum + solve(l+x, r, x, 1)
               if val > res {
                   res = val
               }
           }
       }
   } else {
       // Zhenya: minimize Igor's advantage
       res = (1 << 60)
       for _, x := range moves {
           if x <= rem {
               sum := getSum(r-x+1, r)
               val := solve(l, r-x, x, 0) - sum
               if val < res {
                   res = val
               }
           }
       }
   }
   // if no moves possible, res stays at inf; set to 0
   if res > (1<<59) || res < -(1<<59) {
       res = 0
   }
   memo[st] = res
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   a = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   pref = make([]int64, n+1)
   for i := 0; i < n; i++ {
       pref[i+1] = pref[i] + a[i]
   }
   memo = make(map[state]int64)
   ans := solve(0, n-1, 0, 0)
   fmt.Println(ans)
}
