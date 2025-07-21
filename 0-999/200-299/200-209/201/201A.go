package main

import (
   "bufio"
   "fmt"
   "os"
)

// bitmask for dynamic number of bits
type mask []uint64

// set bit i
func (m mask) set(i int) {
   m[i>>6] |= 1 << (uint(i) & 63)
}

// test bit i
func (m mask) test(i int) bool {
   return (m[i>>6] & (1 << (uint(i) & 63))) != 0
}

// or with other, return new mask
func (m mask) or(o mask) mask {
   n := make(mask, len(m))
   for i := range m {
       n[i] = m[i] | o[i]
   }
   return n
}

// solve for given n and x, return true if possible
func possible(n, x int) bool {
   half := (n + 1) / 2
   // build reps
   // rep at (i,j) for i,j in [1..half]
   type rep struct{ i, j int }
   var reps []rep
   for i := 1; i <= half; i++ {
       for j := 1; j <= half; j++ {
           reps = append(reps, rep{i, j})
       }
   }
   K := len(reps)
   // compute orbit cells and sizes
   // old rep sizes and conflict masks
   oldSize := make([]int, K)
   oldMask := make([]mask, K)
   words := (K + 63) / 64
   // prepare cell lists
   cells := make([][][2]int, K)
   for idx, r := range reps {
       seen := map[[2]int]bool{}
       var cl [][2]int
       for _, pi := range []int{r.i, n - r.i + 1} {
           for _, pj := range []int{r.j, n - r.j + 1} {
               c := [2]int{pi, pj}
               if !seen[c] {
                   seen[c] = true
                   cl = append(cl, c)
               }
           }
       }
       oldSize[idx] = len(cl)
       cells[idx] = cl
   }
   // build conflicts
   for a := 0; a < K; a++ {
       oldMask[a] = make(mask, words)
       for b := 0; b < K; b++ {
           if a == b {
               oldMask[a].set(b)
               continue
           }
           conflict := false
           for _, c1 := range cells[a] {
               for _, c2 := range cells[b] {
                   if abs(c1[0]-c2[0])+abs(c1[1]-c2[1]) == 1 {
                       conflict = true
                       break
                   }
               }
               if conflict {
                   break
               }
           }
           if conflict {
               oldMask[a].set(b)
           }
       }
   }
   // sort reps by size descending
   idxs := make([]int, K)
   for i := range idxs {
       idxs[i] = i
   }
   // simple stable sort by oldSize desc
   for i := 0; i < K; i++ {
       for j := i + 1; j < K; j++ {
           if oldSize[idxs[j]] > oldSize[idxs[i]] {
               idxs[i], idxs[j] = idxs[j], idxs[i]
           }
       }
   }
   // build new sizes and masks
   newSize := make([]int, K)
   newMask := make([]mask, K)
   for ni, oi := range idxs {
       newSize[ni] = oldSize[oi]
       // remap mask
       m := make(mask, words)
       for w := 0; w < words; w++ {
           m[w] = 0
       }
       for nj, oj := range idxs {
           if oldMask[oi].test(oj) {
               m.set(nj)
           }
       }
       newMask[ni] = m
   }
   // find center index in new order
   centerOld := -1
   if n%2 == 1 {
       mid := (n + 1) / 2
       for idx, r := range reps {
           if r.i == mid && r.j == mid {
               centerOld = idx
               break
           }
       }
       if centerOld < 0 {
           return false
       }
   }
   center := -1
   if centerOld >= 0 {
       for ni, oi := range idxs {
           if oi == centerOld {
               center = ni
               break
           }
       }
   }
   // if x odd, must have center
   if x%2 == 1 {
       if center < 0 {
           return false
       }
   }
   // prefix sums of sizes
   prefix := make([]int, K+1)
   for i := K - 1; i >= 0; i-- {
       prefix[i] = prefix[i+1] + newSize[i]
   }
   target := x
   banned0 := make(mask, words)
   currSum := 0
   // include center if needed
   if x%2 == 1 {
       currSum = 1
       target = x - 1
       // ban center and its conflicts
       banned0 = banned0.or(newMask[center])
   }
   // dfs
   var dfs func(int, int, mask) bool
   dfs = func(cur, sum int, banned mask) bool {
       if sum == target {
           return true
       }
       if sum > target || cur >= K {
           return false
       }
       if sum+prefix[cur] < target {
           return false
       }
       for i := cur; i < K; i++ {
           sz := newSize[i]
           if sum+sz <= target && !banned.test(i) {
               // include
               nb := banned.or(newMask[i])
               if dfs(i+1, sum+sz, nb) {
                   return true
               }
           }
       }
       return false
   }
   return dfs(0, currSum, banned0)
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var x int
   if _, err := fmt.Fscan(in, &x); err != nil {
       return
   }
   maxN := 2*x + 1
   for n := 1; n <= maxN; n++ {
       if possible(n, x) {
           fmt.Println(n)
           return
       }
   }
}
