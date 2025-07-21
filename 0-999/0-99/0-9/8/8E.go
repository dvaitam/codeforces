package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k uint64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // prefix bits: -1 unknown, 0 or 1 assigned
   prefix := make([]int8, n)
   for i := range prefix {
       prefix[i] = -1
   }
   numPairs := n / 2
   // DP memo arrays
   var dpMemo [51][2][2]uint64
   var dpVis [51][2][2]bool

   // count completions given current prefix
   var countDP func(i int, revState, crState int) uint64
   countDP = func(i int, revState, crState int) uint64 {
       // comp constraint: first bit must be 0 if assigned
       if prefix[0] == 1 {
           return 0
       }
       if i == numPairs {
           // handle middle for odd n
           if n%2 == 1 {
               mid := numPairs
               bit := prefix[mid]
               if bit >= 0 {
                   // assigned
                   if crState == 0 && bit == 1 {
                       return 0
                   }
                   // valid
                   return 1
               }
               // unassigned
               if crState == 0 {
                   // only bit 0 allowed
                   return 1
               }
               // both bits allowed
               return 2
           }
           return 1
       }
       if dpVis[i][revState][crState] {
           return dpMemo[i][revState][crState]
       }
       var res uint64 = 0
       // determine indices
       j := n - 1 - i
       ai := prefix[i]
       bj := prefix[j]
       // case: both assigned
       if ai >= 0 && bj >= 0 {
           a := ai
           b := bj
           // rev transition
           ns := revState
           if ns == 0 {
               if a < b {
                   ns = 1
               } else if a > b {
                   // invalid
                   dpVis[i][revState][crState] = true
                   dpMemo[i][revState][crState] = 0
                   return 0
               }
           }
           // comp_rev transition
           nc := crState
           if nc == 0 {
               // compare a vs 1-b
               if a < (1 - b) {
                   nc = 1
               } else if a > (1 - b) {
                   dpVis[i][revState][crState] = true
                   dpMemo[i][revState][crState] = 0
                   return 0
               }
           }
           res = countDP(i+1, ns, nc)
       } else if ai >= 0 {
           // left assigned only, assign b
           a := ai
           for b := int8(0); b <= 1; b++ {
               // transitions
               // rev
               ns := revState
               if ns == 0 {
                   if a < b {
                       ns = 1
                   } else if a > b {
                       continue
                   }
               }
               // comp_rev
               nc := crState
               if nc == 0 {
                   if a < (1 - b) {
                       nc = 1
                   } else if a > (1 - b) {
                       continue
                   }
               }
               // assign b temporarily
               prefix[j] = b
               res += countDP(i+1, ns, nc)
               prefix[j] = -1
           }
       } else {
           // both unassigned, assign a,b
           for a := int8(0); a <= 1; a++ {
               if i == 0 && a != 0 {
                   continue
               }
               for b := int8(0); b <= 1; b++ {
                   // rev
                   ns := revState
                   if ns == 0 {
                       if a < b {
                           ns = 1
                       } else if a > b {
                           continue
                       }
                   }
                   // comp_rev
                   nc := crState
                   if nc == 0 {
                       if a < (1 - b) {
                           nc = 1
                       } else if a > (1 - b) {
                           continue
                       }
                   }
                   // assign both
                   prefix[i], prefix[j] = a, b
                   res += countDP(i+1, ns, nc)
                   prefix[i], prefix[j] = -1, -1
               }
           }
       }
       dpVis[i][revState][crState] = true
       dpMemo[i][revState][crState] = res
       return res
   }

   // total minimal reps including all-zero
   total := countDP(0, 0, 0)
   // skip all-zero rep => map user k to K = k+1
   if k+1 > total {
       fmt.Println(-1)
       return
   }
   K := k + 1
   // enumeration by bits
   // first bit must be 0
   prefix[0] = 0
   for pos := 1; pos < n; pos++ {
       // try 0
       prefix[pos] = 0
       // reset DP memo
       for i := 0; i <= numPairs; i++ {
           for rv := 0; rv < 2; rv++ {
               for cr := 0; cr < 2; cr++ {
                   dpVis[i][rv][cr] = false
               }
           }
       }
       cnt0 := countDP(0, 0, 0)
       if K <= cnt0 {
           // keep 0
           continue
       }
       // choose 1
       K -= cnt0
       prefix[pos] = 1
       // no need to check cnt1 here
   }
   // output result
   out := make([]byte, n)
   for i := 0; i < n; i++ {
       if prefix[i] == 0 {
           out[i] = '0'
       } else {
           out[i] = '1'
       }
   }
   fmt.Println(string(out))
}
