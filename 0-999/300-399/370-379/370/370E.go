package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair represents (bookCount, daysOnBook)
type Pair struct {
   first, second int
}

// less returns true if a < b lexicographically
func (a Pair) less(b Pair) bool {
   if a.first != b.first {
       return a.first < b.first
   }
   return a.second < b.second
}

func pairMin(a, b Pair) Pair {
   if a.less(b) {
       return a
   }
   return b
}

func pairMax(a, b Pair) Pair {
   if a.less(b) {
       return b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   if a[1] > 1 {
       fmt.Fprintln(out, -1)
       return
   }
   Max := make([]Pair, n+1)
   Min := make([]Pair, n+1)
   Max[1], Min[1] = Pair{1, 1}, Pair{1, 1}
   flag := true
   for i := 2; i <= n && flag; i++ {
       prevMax, prevMin := Max[i-1], Min[i-1]
       var curMax Pair
       if prevMax.second >= 2 {
           curMax = Pair{prevMax.first + 1, 1}
       } else {
           curMax = Pair{prevMax.first, prevMax.second + 1}
       }
       var curMin Pair
       if prevMin.second == 5 {
           curMin = Pair{prevMin.first + 1, 1}
       } else {
           curMin = Pair{prevMin.first, prevMin.second + 1}
       }
       if a[i] != 0 {
           if a[i] < curMin.first || a[i] > curMax.first {
               flag = false
           }
       } else {
           // clamp days on book
           // as in original: Max[i] = min(Max[i], Pair{a[i],5}); Min[i] = max(Min[i], Pair{a[i],1})
           pMax := Pair{a[i], 5}
           pMin := Pair{a[i], 1}
           curMax = pairMin(curMax, pMax)
           curMin = pairMax(curMin, pMin)
       }
       Max[i], Min[i] = curMax, curMin
   }
   if !flag {
       fmt.Fprintln(out, -1)
       return
   }
   ans := Max[n]
   if ans.second == 1 {
       ans.first--
       ans.second = 5
   }
   if ans.less(Min[n]) {
       fmt.Fprintln(out, -1)
       return
   }
   // print result
   fmt.Fprintln(out, ans.first)
   // reconstruct
   cnt := make([]int, ans.first+3)
   a[n] = ans.first
   cnt[a[n]] = 1
   for i := n - 1; i >= 1; i-- {
       // take min of Max[i].first and next day's book
       bv := Max[i].first
       if a[i+1] < bv {
           bv = a[i+1]
       }
       a[i] = bv
       if cnt[a[i]] == 5 {
           a[i]--
       }
       cnt[a[i]]++
   }
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, a[i])
   }
   out.WriteByte('\n')
}
