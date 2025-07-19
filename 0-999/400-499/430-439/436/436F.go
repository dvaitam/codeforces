package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   maxg = 303
   INF  = int(0x3f3f3f3f)
)

type Pair struct {
   first  int
   second int
}

// Block for sqrt decomposition
type Block struct {
   ans      int64
   pans     int
   rightcnt int
   nextrcnt int
   ld, ldx  int
}

var (
   ab      []Pair
   have    []bool
   blocks  []Block
   userByB [][]int
)

// getLeader finds best leader in [from, to)
func (b *Block) getLeader(from, to int) {
   ourcnt := 0
   b.ans = 0
   b.pans = 0
   for i := to - 1; i >= from; i-- {
       if !have[i] {
           continue
       }
       ourcnt++
       curans := int64(b.rightcnt+ourcnt) * int64(ab[i].first)
       if curans > b.ans {
           b.ans = curans
           b.pans = ab[i].first
           b.ld = i
           b.ldx = b.rightcnt + ourcnt
       }
   }
}

// recalc recomputes ans, pans and nextrcnt for block in [from, to)
func (b *Block) recalc(from, to int) {
   b.getLeader(from, to)
   ourcnt := 0
   b.nextrcnt = INF
   for i := to - 1; i >= from; i-- {
       if !have[i] {
           continue
       }
       ourcnt++
       if i == b.ld || ab[i].first <= ab[b.ld].first {
           continue
       }
       curx := b.rightcnt + ourcnt
       lp := int64(ab[b.ld].first)*int64(b.ldx) - int64(ab[i].first)*int64(curx)
       // lp >= 0
       rp := int64(ab[i].first) - int64(ab[b.ld].first)
       t := int(lp / rp)
       t++
       if t < INF {
           nxt := b.rightcnt + t
           if nxt < b.nextrcnt {
               b.nextrcnt = nxt
           }
       }
   }
}

// update increases rightcnt and ans, and possibly triggers recalc
func (b *Block) update(from, to int) {
   b.rightcnt++
   b.ans += int64(b.pans)
   if b.rightcnt >= b.nextrcnt {
       b.recalc(from, to)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var w int64
   fmt.Fscan(in, &n, &w)
   ab = make([]Pair, n)
   var mb int
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       ab[i] = Pair{a, b}
       if b > mb {
           mb = b
       }
   }
   sort.Slice(ab, func(i, j int) bool {
       return ab[i].first < ab[j].first
   })
   userByB = make([][]int, mb+1)
   for i := 0; i < n; i++ {
       b := ab[i].second
       userByB[b] = append(userByB[b], i)
   }
   have = make([]bool, n)
   // initialize blocks
   blockCount := (n + maxg - 1) / maxg
   blocks = make([]Block, blockCount)
   for i := range blocks {
       blocks[i].nextrcnt = INF
   }

   watchAds := n
   // process for c = 0..mb+1
   for c := 0; c <= mb+1; c++ {
       if c > 0 {
           for _, u := range userByB[c-1] {
               watchAds--
               if ab[u].first == 0 {
                   continue
               }
               have[u] = true
               gid := u / maxg
               from := gid * maxg
               to := from + maxg
               if to > n {
                   to = n
               }
               blocks[gid].recalc(from, to)
               for i := 0; i < gid; i++ {
                   f := i * maxg
                   t := f + maxg
                   if t > n {
                       t = n
                   }
                   blocks[i].update(f, t)
               }
           }
       }
       // find best block
       var bestAns int64
       var bestPans int
       for i := range blocks {
           if blocks[i].ans > bestAns {
               bestAns = blocks[i].ans
               bestPans = blocks[i].pans
           }
       }
       total := bestAns + w*int64(c)*int64(watchAds)
       fmt.Fprintf(out, "%d %d\n", total, bestPans)
   }
}
