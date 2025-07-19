package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // initialize containers
   ansFour := make([][4]int, n+2)
   ansThree := make([][3]int, n+2)
   freq := make([]int, n+2)
   nextGreater := make([]int, n+2)
   nextLess := make([]int, n+2)
   // BIT for available positions
   bit := make([]int, n+3)

   // BIT functions
   add := func(i, x int) {
       for j := i; j < len(bit); j += j & -j {
           bit[j] += x
       }
   }
   query := func(i int) int {
       s := 0
       for j := i; j > 0; j -= j & -j {
           s += bit[j]
       }
       return s
   }
   // find smallest idx such that prefix sum >= target
   getIndex := func(target int) int {
       if target == 0 {
           return 0
       }
       sum, pos := 0, 0
       // largest power of two <= n+2
       step := 1
       for step<<1 < len(bit) {
           step <<= 1
       }
       for k := step; k > 0; k >>= 1 {
           if pos+k < len(bit) && sum+bit[pos+k] < target {
               sum += bit[pos+k]
               pos += k
           }
       }
       return pos + 1
   }

   // initialize answer for position n+1
   inf := n + 1
   ansFour[n+1] = [4]int{inf, inf, inf, inf}
   ansThree[n+1] = [3]int{inf, inf, inf}
   // stacks for next greater/less
   maxstk := make([]int, 0, n)
   minstk := make([]int, 0, n)
   rpos := inf
   add(rpos, 1)
   nextGreater[inf] = inf
   nextLess[inf] = inf

   for i := n; i >= 1; i-- {
       // carry over previous results
       ansFour[i] = ansFour[i+1]
       ansThree[i] = ansThree[i+1]
       // maintain max stack
       for len(maxstk) > 0 && a[maxstk[len(maxstk)-1]] < a[i] {
           idx := maxstk[len(maxstk)-1]
           freq[idx]--
           if freq[idx] == 0 {
               add(idx, 1)
           }
           maxstk = maxstk[:len(maxstk)-1]
       }
       if len(maxstk) == 0 {
           nextGreater[i] = inf
       } else if a[maxstk[len(maxstk)-1]] > a[i] {
           nextGreater[i] = maxstk[len(maxstk)-1]
       } else {
           nextGreater[i] = nextGreater[maxstk[len(maxstk)-1]]
       }
       // maintain min stack
       for len(minstk) > 0 && a[minstk[len(minstk)-1]] > a[i] {
           idx := minstk[len(minstk)-1]
           freq[idx]--
           if freq[idx] == 0 {
               add(idx, 1)
           }
           minstk = minstk[:len(minstk)-1]
       }
       if len(minstk) == 0 {
           nextLess[i] = inf
       } else if a[minstk[len(minstk)-1]] < a[i] {
           nextLess[i] = minstk[len(minstk)-1]
       } else {
           nextLess[i] = nextLess[minstk[len(minstk)-1]]
       }
       // push current
       maxstk = append(maxstk, i)
       minstk = append(minstk, i)
       freq[i] = 2
       // try to find four-length sequence
       for {
           rp := getIndex(query(rpos) - 1)
           if rp == 0 {
               break
           }
           // find u in maxstk: first element <= rp
           u := 0
           for _, v := range maxstk {
               if v <= rp {
                   u = v
                   break
               }
           }
           // find v in minstk: first element <= rp
           v := 0
           for _, x := range minstk {
               if x <= rp {
                   v = x
                   break
               }
           }
           if !(a[u] > a[i] && a[u] > a[rp] && a[v] < a[i] && a[v] < a[rp]) {
               break
           }
           rpos = rp
           if u > v {
               u, v = v, u
           }
           ansFour[i] = [4]int{i, u, v, rpos}
       }
       // three-length sequence
       tg := nextGreater[i]
       if nextLess[tg] < ansThree[i][2] {
           ansThree[i] = [3]int{i, tg, nextLess[tg]}
       }
       tl := nextLess[i]
       if nextGreater[tl] < ansThree[i][2] {
           ansThree[i] = [3]int{i, tl, nextGreater[tl]}
       }
   }

   // answer queries
   for qq := 0; qq < q; qq++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       af := ansFour[l]
       if af[3] <= r {
           fmt.Fprint(writer, 4, '\n')
           fmt.Fprintln(writer, af[0], af[1], af[2], af[3])
       } else {
           at := ansThree[l]
           if at[2] <= r {
               fmt.Fprint(writer, 3, '\n')
               fmt.Fprintln(writer, at[0], at[1], at[2])
           } else {
               fmt.Fprint(writer, 0, '\n')
           }
       }
   }
}
