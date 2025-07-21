package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &c[i])
   }
   if n == 1 {
       if c[0] == 1 {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
       return
   }
   // Exactly one root with size n
   cntN := 0
   for _, v := range c {
       if v == n {
           cntN++
       }
       if v == 2 {
           fmt.Println("NO")
           return
       }
   }
   if cntN != 1 {
       fmt.Println("NO")
       return
   }
   // initial built trees: leaves
   counts := make(map[int]int)
   var targets []int
   for _, v := range c {
       if v == 1 {
           counts[1]++
       } else {
           targets = append(targets, v)
       }
   }
   sort.Ints(targets)
   // build subtrees bottom-up
   for _, s := range targets {
       need := s - 1
       // collect available sizes < s
       var keys []int
       totalAvailable := 0
       for k, ct := range counts {
           if k < s && ct > 0 {
               keys = append(keys, k)
               totalAvailable += k * ct
           }
       }
       if totalAvailable < need {
           fmt.Println("NO")
           return
       }
       sort.Sort(sort.Reverse(sort.IntSlice(keys)))
       used := make(map[int]int)
       ok := dfsPick(need, 0, 0, keys, counts, used)
       if !ok {
           fmt.Println("NO")
           return
       }
       // consume used children
       childrenCount := 0
       for k, u := range used {
           counts[k] -= u
           childrenCount += u
       }
       if childrenCount < 2 {
           fmt.Println("NO")
           return
       }
       // add this subtree
       counts[s]++
   }
   fmt.Println("YES")
}

// dfsPick tries to pick sizes from keys[pos:] to sum to rem, with at least two items
// counts gives available counts, used accumulates picks
func dfsPick(rem, pos, items int, keys []int, counts map[int]int, used map[int]int) bool {
   if rem == 0 {
       return items >= 2
   }
   if pos >= len(keys) {
       return false
   }
   k := keys[pos]
   maxTake := counts[k]
   if maxTake*k > rem {
       maxTake = rem / k
   }
   for take := maxTake; take >= 0; take-- {
       newRem := rem - take*k
       newItems := items + take
       if newRem == 0 {
           if newItems >= 2 {
               if take > 0 {
                   used[k] = take
               }
               return true
           }
           continue
       }
       // prune: if no further keys, skip
       if pos+1 >= len(keys) {
           continue
       }
       // estimate possible sum from remaining keys
       // we skip heavy pruning for simplicity
       if take > 0 {
           used[k] = take
       }
       if dfsPick(newRem, pos+1, newItems, keys, counts, used) {
           return true
       }
       if take > 0 {
           delete(used, k)
       }
   }
   return false
}
