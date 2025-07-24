package main

import (
   "bufio"
   "fmt"
   "os"
)

// Occurrence of x in sequence at position pos (0-based for pos in sequence)
type Occ struct {
   pos int // sequence index, 0-based
   k   int // sequence length
   r   int // position in sequence, 0-based
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   // Prepare occurrences
   occs := make([][]Occ, m+1)
   seqKs := make([]int, n)
   for i := 0; i < n; i++ {
       var k int
       fmt.Fscan(in, &k)
       seqKs[i] = k
       for j := 0; j < k; j++ {
           var x int
           fmt.Fscan(in, &x)
           occs[x] = append(occs[x], Occ{pos: i, k: k, r: j})
       }
   }
   // Precompute divisors for 1..40
   maxK := 40
   divisors := make([][]int, maxK+1)
   for d := 1; d <= maxK; d++ {
       for k := d; k <= maxK; k += d {
           divisors[k] = append(divisors[k], d)
       }
   }
   // Answer array
   ans := make([]int, m+1)
   // Process each x
   for x := 1; x <= m; x++ {
       list := occs[x]
       L := len(list)
       if L == 0 {
           ans[x] = 0
           continue
       }
       // counts for CRT constraints
       var cnt [41][]int
       tot := make([]int, 41)
       for g := 1; g <= maxK; g++ {
           cnt[g] = make([]int, g)
       }
       best := 0
       runStart := 0
       // split into contiguous position runs
       for idx := 0; idx <= L; idx++ {
           if idx == L || (idx > runStart && list[idx].pos != list[idx-1].pos+1) {
               // process run [runStart, idx)
               if runStart < idx {
                   l := runStart
                   // sliding window
                   for r := runStart; r < idx; r++ {
                       cur := list[r]
                       ds := divisors[cur.k]
                       // shrink until consistent
                       for {
                           conflict := false
                           for _, g := range ds {
                               res := cur.r % g
                               if tot[g] > 0 && cnt[g][res] == 0 {
                                   conflict = true
                                   break
                               }
                           }
                           if !conflict {
                               break
                           }
                           old := list[l]
                           for _, g := range divisors[old.k] {
                               cnt[g][old.r%g]--
                               tot[g]--
                           }
                           l++
                       }
                       // add cur
                       for _, g := range ds {
                           cnt[g][cur.r%g]++
                           tot[g]++
                       }
                       length := r - l + 1
                       if length > best {
                           best = length
                       }
                   }
                   // clear counts for next run
                   for g := 1; g <= maxK; g++ {
                       for i := range cnt[g] {
                           cnt[g][i] = 0
                       }
                       tot[g] = 0
                   }
               }
               runStart = idx
           }
       }
       ans[x] = best
   }
   // output
   for i := 1; i <= m; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%d", ans[i])
   }
   out.WriteByte('\n')
}
