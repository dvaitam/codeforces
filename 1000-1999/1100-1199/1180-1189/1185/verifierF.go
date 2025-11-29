package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1185FSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   const nax = 1 << 9
   fr := make([]int, nax)
   for i := 0; i < n; i++ {
       var k, x int
       fmt.Fscan(reader, &k)
       mask := 0
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &x)
           mask |= 1 << (x - 1)
       }
       fr[mask]++
   }
   ile := make([]int, nax)
   for i := 0; i < nax; i++ {
       if fr[i] == 0 {
           continue
       }
       for x := 0; x < nax; x++ {
           if x&i == i {
               ile[x] += fr[i]
           }
       }
   }
   const INF = int(1e9 + 5)
   type pair struct{ first, second int }
   piz := make([]pair, nax)
   for i := range piz {
       piz[i] = pair{INF, INF}
   }
   bon := pair{INF, INF}
   for idx := 1; idx <= m; idx++ {
       var pr, k, x int
       fmt.Fscan(reader, &pr, &k)
       mask := 0
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &x)
           mask |= 1 << (x - 1)
       }
       if pr < piz[mask].first {
           if piz[mask].first != INF {
               old := piz[mask]
               if old.first < bon.first || (old.first == bon.first && old.second < bon.second) {
                   bon = old
               }
           }
           piz[mask] = pair{pr, idx}
       } else {
           cur := pair{pr, idx}
           if cur.first < bon.first || (cur.first == bon.first && cur.second < bon.second) {
               bon = cur
           }
       }
   }
   ans := pair{INF, INF}
   var ret pair
   for i := 0; i < nax; i++ {
       if piz[i].first == INF {
           continue
       }
       // combine with different masks
       for j := 0; j < nax; j++ {
           if i == j || piz[j].first == INF {
               continue
           }
           cm := i | j
           cov := ile[cm]
           sum := piz[i].first + piz[j].first
           nw := pair{-cov, sum}
           if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
               ans = nw
               ret = pair{piz[i].second, piz[j].second}
           }
       }
       // combine with second best
       if bon.first != INF {
           cov := ile[i]
           sum := piz[i].first + bon.first
           nw := pair{-cov, sum}
           if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
               ans = nw
               ret = pair{piz[i].second, bon.second}
           }
       }
   }
   fmt.Fprintln(writer, ret.first, ret.second)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1185FSource

type offer struct {
	price int
	goods []int
}

type testCase struct {
	n      int
	m      int
	items  [][]int
	offers []offer
}

var testcases = []testCase{
	{n: 5, m: 2, items: [][]int{[]int{5, 1}, []int{3}, []int{8, 6, 3}, []int{5}, []int{4, 7}}, offers: []offer{{price: 18, goods: []int{2, 4, 5}}, {price: 18, goods: []int{5, 2, 7}}}},
	{n: 4, m: 4, items: [][]int{[]int{6}, []int{5, 8}, []int{2, 4, 6}, []int{5, 2, 1}}, offers: []offer{{price: 19, goods: []int{6}}, {price: 16, goods: []int{9}}, {price: 19, goods: []int{9, 1, 6}}, {price: 12, goods: []int{7}}}},
	{n: 3, m: 4, items: [][]int{[]int{2, 9, 5}, []int{9, 4, 1}, []int{5, 9, 6}}, offers: []offer{{price: 7, goods: []int{8, 4}}, {price: 5, goods: []int{4, 1, 2}}, {price: 6, goods: []int{6}}, {price: 18, goods: []int{5, 6, 4}}}},
	{n: 5, m: 5, items: [][]int{[]int{3, 8}, []int{3}, []int{7, 2}, []int{4, 2}, []int{8, 9, 7}}, offers: []offer{{price: 13, goods: []int{9}}, {price: 14, goods: []int{5, 7}}, {price: 3, goods: []int{5}}, {price: 15, goods: []int{3, 1}}, {price: 1, goods: []int{2, 5, 8}}}},
	{n: 3, m: 3, items: [][]int{[]int{9, 8}, []int{9, 5}, []int{7, 8}}, offers: []offer{{price: 9, goods: []int{8, 9, 4}}, {price: 5, goods: []int{7, 8, 3}}, {price: 15, goods: []int{6, 3}}}},
	{n: 5, m: 5, items: [][]int{[]int{5, 6, 8}, []int{8, 3}, []int{9, 1}, []int{8, 1, 2}, []int{2}}, offers: []offer{{price: 12, goods: []int{8, 1}}, {price: 7, goods: []int{5}}, {price: 20, goods: []int{7}}, {price: 17, goods: []int{2, 8}}, {price: 8, goods: []int{6}}}},
	{n: 3, m: 3, items: [][]int{[]int{4, 6, 2}, []int{2}, []int{2, 1, 6}}, offers: []offer{{price: 16, goods: []int{2, 1}}, {price: 18, goods: []int{2}}, {price: 6, goods: []int{5, 8, 4}}}},
	{n: 1, m: 3, items: [][]int{[]int{8}}, offers: []offer{{price: 16, goods: []int{8, 5}}, {price: 16, goods: []int{2, 5, 3}}, {price: 10, goods: []int{1, 2}}}},
	{n: 5, m: 4, items: [][]int{[]int{5, 8}, []int{8, 6}, []int{2}, []int{4}, []int{1, 3, 8}}, offers: []offer{{price: 12, goods: []int{1}}, {price: 7, goods: []int{7, 6, 1}}, {price: 6, goods: []int{2, 7, 5}}, {price: 5, goods: []int{8, 2}}}},
	{n: 2, m: 2, items: [][]int{[]int{9, 1}, []int{9}}, offers: []offer{{price: 5, goods: []int{5}}, {price: 5, goods: []int{7, 4, 2}}}},
	{n: 5, m: 3, items: [][]int{[]int{3, 6}, []int{9}, []int{1, 7, 2}, []int{7, 8}, []int{7}}, offers: []offer{{price: 17, goods: []int{3, 4, 6}}, {price: 18, goods: []int{6, 5, 9}}, {price: 6, goods: []int{3}}}},
	{n: 2, m: 5, items: [][]int{[]int{1}, []int{1}}, offers: []offer{{price: 16, goods: []int{2}}, {price: 18, goods: []int{6}}, {price: 20, goods: []int{1}}, {price: 18, goods: []int{7}}, {price: 8, goods: []int{5, 2, 4}}}},
	{n: 1, m: 2, items: [][]int{[]int{8}}, offers: []offer{{price: 4, goods: []int{3}}, {price: 2, goods: []int{1, 7, 6}}}},
	{n: 3, m: 4, items: [][]int{[]int{4, 9}, []int{9, 4, 3}, []int{1, 5, 8}}, offers: []offer{{price: 11, goods: []int{2, 1, 3}}, {price: 12, goods: []int{2, 4, 7}}, {price: 1, goods: []int{6, 7}}, {price: 2, goods: []int{3, 1}}}},
	{n: 1, m: 4, items: [][]int{[]int{8}}, offers: []offer{{price: 1, goods: []int{5, 2}}, {price: 19, goods: []int{5, 4, 2}}, {price: 5, goods: []int{9, 3}}, {price: 8, goods: []int{7, 9, 5}}}},
	{n: 5, m: 5, items: [][]int{[]int{5, 3}, []int{1, 4}, []int{2, 3}, []int{2, 1, 5}, []int{2}}, offers: []offer{{price: 5, goods: []int{7, 8, 4}}, {price: 17, goods: []int{9, 8, 2}}, {price: 13, goods: []int{8}}, {price: 1, goods: []int{7, 5}}, {price: 11, goods: []int{5}}}},
	{n: 3, m: 4, items: [][]int{[]int{1}, []int{3, 2}, []int{4, 2}}, offers: []offer{{price: 7, goods: []int{7, 3}}, {price: 5, goods: []int{1, 8, 3}}, {price: 11, goods: []int{7, 5, 6}}, {price: 3, goods: []int{3, 1}}}},
	{n: 1, m: 3, items: [][]int{[]int{9, 2, 7}}, offers: []offer{{price: 3, goods: []int{2}}, {price: 7, goods: []int{6, 3}}, {price: 8, goods: []int{3}}}},
	{n: 2, m: 3, items: [][]int{[]int{8, 2}, []int{3, 8, 7}}, offers: []offer{{price: 7, goods: []int{6, 3, 8}}, {price: 9, goods: []int{5}}, {price: 1, goods: []int{1}}}},
	{n: 4, m: 5, items: [][]int{[]int{1, 9}, []int{2, 7}, []int{6, 2, 9}, []int{3, 9}}, offers: []offer{{price: 12, goods: []int{8, 7, 9}}, {price: 1, goods: []int{3, 2, 4}}, {price: 12, goods: []int{3}}, {price: 18, goods: []int{5, 2}}, {price: 16, goods: []int{7}}}},
	{n: 2, m: 3, items: [][]int{[]int{6, 5}, []int{6, 5, 8}}, offers: []offer{{price: 9, goods: []int{9, 1}}, {price: 17, goods: []int{4, 1, 5}}, {price: 16, goods: []int{9, 5, 2}}}},
	{n: 1, m: 4, items: [][]int{[]int{4}}, offers: []offer{{price: 19, goods: []int{4, 3}}, {price: 9, goods: []int{5}}, {price: 10, goods: []int{6, 4, 9}}, {price: 6, goods: []int{1, 2, 7}}}},
	{n: 3, m: 3, items: [][]int{[]int{3, 8}, []int{4, 8, 9}, []int{6, 2, 9}}, offers: []offer{{price: 3, goods: []int{1, 4}}, {price: 20, goods: []int{5, 4}}, {price: 2, goods: []int{9, 5}}}},
	{n: 5, m: 4, items: [][]int{[]int{6, 5, 7}, []int{6, 2}, []int{7, 3, 8}, []int{9}, []int{9, 1}}, offers: []offer{{price: 3, goods: []int{2, 1, 4}}, {price: 14, goods: []int{2, 3, 6}}, {price: 9, goods: []int{1, 2}}, {price: 9, goods: []int{5, 6, 4}}}},
	{n: 4, m: 2, items: [][]int{[]int{2, 5, 7}, []int{8, 5}, []int{4, 1}, []int{9}}, offers: []offer{{price: 18, goods: []int{3, 4}}, {price: 10, goods: []int{1, 2, 3}}}},
	{n: 2, m: 2, items: [][]int{[]int{7}, []int{8, 4}}, offers: []offer{{price: 3, goods: []int{5, 3, 2}}, {price: 5, goods: []int{1}}}},
	{n: 2, m: 4, items: [][]int{[]int{8, 9}, []int{6, 8}}, offers: []offer{{price: 7, goods: []int{9}}, {price: 4, goods: []int{3}}, {price: 13, goods: []int{6}}, {price: 15, goods: []int{5, 4, 9}}}},
	{n: 2, m: 2, items: [][]int{[]int{4, 7}, []int{2, 1, 3}}, offers: []offer{{price: 12, goods: []int{5, 7}}, {price: 13, goods: []int{6, 9}}}},
	{n: 5, m: 2, items: [][]int{[]int{3, 5}, []int{8, 4}, []int{8}, []int{3, 6, 5}, []int{8, 6}}, offers: []offer{{price: 3, goods: []int{7, 8}}, {price: 1, goods: []int{2}}}},
	{n: 4, m: 5, items: [][]int{[]int{1, 9, 7}, []int{9, 6, 3}, []int{3}, []int{9, 2}}, offers: []offer{{price: 13, goods: []int{2, 1, 3}}, {price: 8, goods: []int{9, 5, 7}}, {price: 20, goods: []int{5, 7, 3}}, {price: 1, goods: []int{6, 7, 3}}, {price: 6, goods: []int{6, 1}}}},
	{n: 1, m: 3, items: [][]int{[]int{1}}, offers: []offer{{price: 1, goods: []int{2, 7}}, {price: 16, goods: []int{8, 4}}, {price: 7, goods: []int{3, 2}}}},
	{n: 2, m: 2, items: [][]int{[]int{5}, []int{4, 6}}, offers: []offer{{price: 16, goods: []int{3, 6}}, {price: 9, goods: []int{9}}}},
	{n: 1, m: 3, items: [][]int{[]int{7}}, offers: []offer{{price: 17, goods: []int{7}}, {price: 1, goods: []int{2, 9, 1}}, {price: 20, goods: []int{4, 2, 5}}}},
	{n: 5, m: 5, items: [][]int{[]int{3, 8}, []int{1}, []int{6}, []int{9}, []int{5, 3, 9}}, offers: []offer{{price: 1, goods: []int{5}}, {price: 4, goods: []int{9, 4}}, {price: 7, goods: []int{7, 8, 1}}, {price: 14, goods: []int{9, 3, 2}}, {price: 20, goods: []int{5}}}},
	{n: 4, m: 4, items: [][]int{[]int{7, 1, 4}, []int{4}, []int{4, 7}, []int{1, 6}}, offers: []offer{{price: 12, goods: []int{6, 7, 4}}, {price: 12, goods: []int{7, 5, 9}}, {price: 5, goods: []int{5, 4, 9}}, {price: 15, goods: []int{3, 1, 2}}}},
	{n: 1, m: 3, items: [][]int{[]int{5, 8, 9}}, offers: []offer{{price: 1, goods: []int{9}}, {price: 10, goods: []int{9}}, {price: 2, goods: []int{7, 2}}}},
	{n: 3, m: 4, items: [][]int{[]int{5}, []int{1, 5, 7}, []int{5}}, offers: []offer{{price: 11, goods: []int{9}}, {price: 20, goods: []int{2, 1}}, {price: 8, goods: []int{4, 9}}, {price: 3, goods: []int{1, 6}}}},
	{n: 3, m: 3, items: [][]int{[]int{7, 9, 8}, []int{1, 9, 2}, []int{9, 4, 5}}, offers: []offer{{price: 8, goods: []int{3, 5}}, {price: 6, goods: []int{6, 9, 3}}, {price: 19, goods: []int{3}}}},
	{n: 1, m: 3, items: [][]int{[]int{9, 8, 1}}, offers: []offer{{price: 4, goods: []int{9}}, {price: 2, goods: []int{6, 8}}, {price: 13, goods: []int{9, 4}}}},
	{n: 4, m: 5, items: [][]int{[]int{8, 7}, []int{9, 4, 1}, []int{2, 7, 6}, []int{7, 2, 5}}, offers: []offer{{price: 17, goods: []int{6, 3}}, {price: 4, goods: []int{2, 9, 8}}, {price: 7, goods: []int{4, 1}}, {price: 12, goods: []int{3, 1, 8}}, {price: 13, goods: []int{5}}}},
	{n: 4, m: 2, items: [][]int{[]int{8}, []int{1}, []int{8}, []int{6, 4, 2}}, offers: []offer{{price: 15, goods: []int{7, 4}}, {price: 11, goods: []int{4, 2, 7}}}},
	{n: 4, m: 5, items: [][]int{[]int{1, 4, 3}, []int{5, 2, 9}, []int{4, 6}, []int{6, 5, 3}}, offers: []offer{{price: 13, goods: []int{6, 8}}, {price: 12, goods: []int{2}}, {price: 19, goods: []int{5, 2}}, {price: 8, goods: []int{9, 1}}, {price: 3, goods: []int{2, 3, 7}}}},
	{n: 1, m: 4, items: [][]int{[]int{7}}, offers: []offer{{price: 2, goods: []int{4}}, {price: 1, goods: []int{5}}, {price: 6, goods: []int{4}}, {price: 4, goods: []int{4, 1, 6}}}},
	{n: 3, m: 5, items: [][]int{[]int{1}, []int{6, 3, 7}, []int{3, 7}}, offers: []offer{{price: 9, goods: []int{7}}, {price: 10, goods: []int{8, 1, 6}}, {price: 10, goods: []int{3, 1}}, {price: 7, goods: []int{2, 3}}, {price: 1, goods: []int{6, 9, 1}}}},
	{n: 5, m: 3, items: [][]int{[]int{7, 6}, []int{2, 4, 3}, []int{8}, []int{7, 2, 3}, []int{4, 1}}, offers: []offer{{price: 14, goods: []int{2, 3}}, {price: 1, goods: []int{1}}, {price: 1, goods: []int{2, 8}}}},
	{n: 1, m: 3, items: [][]int{[]int{5, 9}}, offers: []offer{{price: 8, goods: []int{4}}, {price: 11, goods: []int{5, 3, 1}}, {price: 3, goods: []int{8}}}},
	{n: 1, m: 5, items: [][]int{[]int{8, 7, 3}}, offers: []offer{{price: 5, goods: []int{2, 3, 8}}, {price: 8, goods: []int{4}}, {price: 11, goods: []int{7}}, {price: 12, goods: []int{6, 3, 4}}, {price: 7, goods: []int{6}}}},
	{n: 2, m: 3, items: [][]int{[]int{3}, []int{3}}, offers: []offer{{price: 16, goods: []int{4, 7, 6}}, {price: 17, goods: []int{8}}, {price: 15, goods: []int{9, 1, 7}}}},
	{n: 1, m: 4, items: [][]int{[]int{9}}, offers: []offer{{price: 11, goods: []int{8, 9, 7}}, {price: 8, goods: []int{7}}, {price: 13, goods: []int{8}}, {price: 5, goods: []int{8}}}},
	{n: 1, m: 2, items: [][]int{[]int{2}}, offers: []offer{{price: 3, goods: []int{4}}, {price: 7, goods: []int{4, 8, 6}}}},
	{n: 2, m: 2, items: [][]int{[]int{9}, []int{6, 2, 1}}, offers: []offer{{price: 20, goods: []int{1, 6}}, {price: 8, goods: []int{3, 2, 7}}}},
	{n: 5, m: 5, items: [][]int{[]int{4, 6}, []int{4}, []int{1, 6}, []int{2, 9, 1}, []int{2, 9}}, offers: []offer{{price: 13, goods: []int{2, 7, 5}}, {price: 5, goods: []int{8, 3}}, {price: 1, goods: []int{1}}, {price: 6, goods: []int{3}}, {price: 10, goods: []int{7, 9, 1}}}},
	{n: 3, m: 2, items: [][]int{[]int{9, 4}, []int{7, 6}, []int{9}}, offers: []offer{{price: 11, goods: []int{6, 2}}, {price: 13, goods: []int{5, 3, 7}}}},
	{n: 3, m: 5, items: [][]int{[]int{3}, []int{7, 5}, []int{7}}, offers: []offer{{price: 10, goods: []int{4, 5, 2}}, {price: 8, goods: []int{1, 9}}, {price: 14, goods: []int{8}}, {price: 8, goods: []int{9, 6}}, {price: 16, goods: []int{4, 9}}}},
	{n: 5, m: 2, items: [][]int{[]int{3}, []int{8, 4, 2}, []int{6, 5}, []int{8}, []int{4}}, offers: []offer{{price: 3, goods: []int{1}}, {price: 10, goods: []int{5, 4}}}},
	{n: 2, m: 4, items: [][]int{[]int{9}, []int{8, 5, 9}}, offers: []offer{{price: 7, goods: []int{4, 1}}, {price: 17, goods: []int{2}}, {price: 20, goods: []int{7}}, {price: 14, goods: []int{5, 4}}}},
	{n: 5, m: 3, items: [][]int{[]int{8}, []int{4, 8}, []int{8, 1, 7}, []int{9, 7, 3}, []int{4, 5, 2}}, offers: []offer{{price: 16, goods: []int{5}}, {price: 10, goods: []int{9, 3}}, {price: 15, goods: []int{4, 5, 6}}}},
	{n: 1, m: 4, items: [][]int{[]int{6, 9}}, offers: []offer{{price: 5, goods: []int{5}}, {price: 16, goods: []int{2, 9, 7}}, {price: 3, goods: []int{2, 7, 3}}, {price: 16, goods: []int{3}}}},
	{n: 3, m: 2, items: [][]int{[]int{4, 1, 9}, []int{1, 9, 6}, []int{5, 6}}, offers: []offer{{price: 2, goods: []int{5, 7, 1}}, {price: 18, goods: []int{3, 4, 9}}}},
	{n: 1, m: 2, items: [][]int{[]int{4, 3}}, offers: []offer{{price: 15, goods: []int{5}}, {price: 6, goods: []int{1, 8, 4}}}},
	{n: 2, m: 2, items: [][]int{[]int{4}, []int{2, 7, 9}}, offers: []offer{{price: 4, goods: []int{4}}, {price: 13, goods: []int{9}}}},
	{n: 2, m: 5, items: [][]int{[]int{1}, []int{2, 4}}, offers: []offer{{price: 17, goods: []int{4, 5}}, {price: 5, goods: []int{6, 2, 3}}, {price: 16, goods: []int{4}}, {price: 4, goods: []int{5, 6, 8}}, {price: 1, goods: []int{8, 3}}}},
	{n: 4, m: 5, items: [][]int{[]int{7}, []int{9}, []int{7, 1, 4}, []int{6}}, offers: []offer{{price: 2, goods: []int{9}}, {price: 16, goods: []int{8, 9, 3}}, {price: 15, goods: []int{2}}, {price: 20, goods: []int{8}}, {price: 15, goods: []int{6, 9, 2}}}},
	{n: 4, m: 5, items: [][]int{[]int{7, 4, 9}, []int{3, 6}, []int{7, 8, 1}, []int{1}}, offers: []offer{{price: 5, goods: []int{4, 3}}, {price: 13, goods: []int{5}}, {price: 9, goods: []int{6, 1, 5}}, {price: 13, goods: []int{3, 5, 9}}, {price: 12, goods: []int{2, 5, 9}}}},
	{n: 4, m: 2, items: [][]int{[]int{4}, []int{7, 4, 2}, []int{1, 9, 8}, []int{9, 4}}, offers: []offer{{price: 10, goods: []int{5, 1}}, {price: 12, goods: []int{3}}}},
	{n: 1, m: 5, items: [][]int{[]int{1, 4}}, offers: []offer{{price: 6, goods: []int{7, 1, 5}}, {price: 13, goods: []int{4}}, {price: 3, goods: []int{7, 3, 1}}, {price: 10, goods: []int{6}}, {price: 15, goods: []int{4, 3, 2}}}},
	{n: 1, m: 5, items: [][]int{[]int{2}}, offers: []offer{{price: 6, goods: []int{1, 2}}, {price: 17, goods: []int{4}}, {price: 14, goods: []int{8}}, {price: 11, goods: []int{6}}, {price: 20, goods: []int{3, 8}}}},
	{n: 5, m: 3, items: [][]int{[]int{2, 6}, []int{5, 6, 8}, []int{2, 7}, []int{8}, []int{8, 4}}, offers: []offer{{price: 13, goods: []int{2}}, {price: 7, goods: []int{7, 3, 2}}, {price: 13, goods: []int{4, 1, 9}}}},
	{n: 1, m: 2, items: [][]int{[]int{8}}, offers: []offer{{price: 11, goods: []int{2, 9, 8}}, {price: 2, goods: []int{2, 5, 6}}}},
	{n: 1, m: 2, items: [][]int{[]int{9, 4, 3}}, offers: []offer{{price: 6, goods: []int{2, 1, 6}}, {price: 20, goods: []int{3, 8, 9}}}},
	{n: 1, m: 4, items: [][]int{[]int{4}}, offers: []offer{{price: 1, goods: []int{7, 8}}, {price: 18, goods: []int{2}}, {price: 12, goods: []int{1, 2, 3}}, {price: 15, goods: []int{4, 5, 8}}}},
	{n: 3, m: 5, items: [][]int{[]int{6, 9, 5}, []int{1, 3, 2}, []int{4, 2}}, offers: []offer{{price: 14, goods: []int{5, 4}}, {price: 12, goods: []int{1}}, {price: 9, goods: []int{4, 8}}, {price: 10, goods: []int{3, 2}}, {price: 18, goods: []int{1}}}},
	{n: 3, m: 3, items: [][]int{[]int{9, 8, 6}, []int{3, 9}, []int{8, 4}}, offers: []offer{{price: 17, goods: []int{4, 7, 5}}, {price: 11, goods: []int{1, 9, 4}}, {price: 19, goods: []int{9, 5}}}},
	{n: 1, m: 3, items: [][]int{[]int{7, 2}}, offers: []offer{{price: 12, goods: []int{6, 3}}, {price: 3, goods: []int{2, 1}}, {price: 12, goods: []int{6}}}},
	{n: 4, m: 3, items: [][]int{[]int{1, 7}, []int{9}, []int{1, 2}, []int{3, 6}}, offers: []offer{{price: 5, goods: []int{6}}, {price: 1, goods: []int{2}}, {price: 11, goods: []int{9, 7}}}},
	{n: 1, m: 3, items: [][]int{[]int{2, 3}}, offers: []offer{{price: 18, goods: []int{1}}, {price: 14, goods: []int{7, 2}}, {price: 13, goods: []int{4}}}},
	{n: 4, m: 2, items: [][]int{[]int{4, 6}, []int{5}, []int{6, 7, 2}, []int{8, 7, 3}}, offers: []offer{{price: 19, goods: []int{8}}, {price: 15, goods: []int{2, 6, 3}}}},
	{n: 2, m: 3, items: [][]int{[]int{9, 5}, []int{9, 6, 4}}, offers: []offer{{price: 3, goods: []int{8, 9}}, {price: 16, goods: []int{3, 8, 6}}, {price: 9, goods: []int{1, 2, 3}}}},
	{n: 3, m: 5, items: [][]int{[]int{5}, []int{5}, []int{5, 1, 3}}, offers: []offer{{price: 4, goods: []int{2, 9, 3}}, {price: 20, goods: []int{7}}, {price: 15, goods: []int{8, 2, 6}}, {price: 12, goods: []int{9, 4}}, {price: 12, goods: []int{6}}}},
	{n: 1, m: 5, items: [][]int{[]int{6, 9}}, offers: []offer{{price: 12, goods: []int{7, 3}}, {price: 2, goods: []int{3, 9}}, {price: 17, goods: []int{4}}, {price: 17, goods: []int{2, 4, 8}}, {price: 10, goods: []int{6}}}},
	{n: 4, m: 2, items: [][]int{[]int{5, 4, 7}, []int{9, 3, 8}, []int{9}, []int{1, 3, 2}}, offers: []offer{{price: 16, goods: []int{2}}, {price: 20, goods: []int{9, 4}}}},
	{n: 2, m: 5, items: [][]int{[]int{5, 3}, []int{5, 3, 4}}, offers: []offer{{price: 16, goods: []int{2, 4}}, {price: 6, goods: []int{1}}, {price: 10, goods: []int{6, 1, 8}}, {price: 15, goods: []int{6, 5, 7}}, {price: 16, goods: []int{5, 4, 7}}}},
	{n: 1, m: 3, items: [][]int{[]int{4}}, offers: []offer{{price: 9, goods: []int{1}}, {price: 5, goods: []int{5, 7}}, {price: 10, goods: []int{1}}}},
	{n: 1, m: 3, items: [][]int{[]int{1, 8}}, offers: []offer{{price: 4, goods: []int{1}}, {price: 8, goods: []int{4}}, {price: 5, goods: []int{5}}}},
	{n: 5, m: 3, items: [][]int{[]int{1}, []int{3, 4}, []int{8, 3}, []int{9, 6}, []int{7}}, offers: []offer{{price: 5, goods: []int{2, 6}}, {price: 4, goods: []int{2}}, {price: 8, goods: []int{8, 3}}}},
	{n: 3, m: 2, items: [][]int{[]int{3, 6, 8}, []int{5, 3}, []int{7}}, offers: []offer{{price: 15, goods: []int{5}}, {price: 2, goods: []int{8, 2, 1}}}},
	{n: 4, m: 2, items: [][]int{[]int{4, 1}, []int{1, 4, 7}, []int{5, 3}, []int{2, 7, 9}}, offers: []offer{{price: 12, goods: []int{9, 5}}, {price: 20, goods: []int{3, 4, 8}}}},
	{n: 4, m: 2, items: [][]int{[]int{1}, []int{4, 9, 7}, []int{7, 2}, []int{6, 1, 3}}, offers: []offer{{price: 2, goods: []int{6, 2, 3}}, {price: 5, goods: []int{5, 4}}}},
	{n: 4, m: 4, items: [][]int{[]int{1}, []int{9, 6}, []int{8, 7}, []int{7, 9, 5}}, offers: []offer{{price: 18, goods: []int{6, 2, 7}}, {price: 16, goods: []int{2}}, {price: 6, goods: []int{2}}, {price: 2, goods: []int{5}}}},
	{n: 3, m: 3, items: [][]int{[]int{2, 4, 1}, []int{5, 1, 2}, []int{5, 6, 2}}, offers: []offer{{price: 18, goods: []int{4, 1, 5}}, {price: 14, goods: []int{9, 1, 6}}, {price: 5, goods: []int{2, 1, 7}}}},
	{n: 1, m: 4, items: [][]int{[]int{3, 1}}, offers: []offer{{price: 19, goods: []int{8, 4, 6}}, {price: 20, goods: []int{1, 2}}, {price: 4, goods: []int{7, 6, 8}}, {price: 9, goods: []int{2, 6}}}},
	{n: 3, m: 3, items: [][]int{[]int{2, 5}, []int{1}, []int{3, 9}}, offers: []offer{{price: 8, goods: []int{5, 1}}, {price: 9, goods: []int{5, 8}}, {price: 20, goods: []int{9, 7}}}},
	{n: 3, m: 5, items: [][]int{[]int{2, 7}, []int{2}, []int{2}}, offers: []offer{{price: 10, goods: []int{8, 5, 3}}, {price: 3, goods: []int{1, 3, 2}}, {price: 4, goods: []int{4, 9, 2}}, {price: 11, goods: []int{1, 7}}, {price: 8, goods: []int{3, 6, 4}}}},
	{n: 4, m: 2, items: [][]int{[]int{1}, []int{2}, []int{9}, []int{4, 2}}, offers: []offer{{price: 8, goods: []int{7, 5, 2}}, {price: 13, goods: []int{3}}}},
	{n: 4, m: 4, items: [][]int{[]int{7, 9, 8}, []int{7, 5, 2}, []int{1}, []int{2, 4}}, offers: []offer{{price: 7, goods: []int{3, 6}}, {price: 20, goods: []int{8, 4}}, {price: 12, goods: []int{7}}, {price: 3, goods: []int{7, 4, 2}}}},
	{n: 5, m: 3, items: [][]int{[]int{8, 3, 4}, []int{4, 2}, []int{4, 6}, []int{5}, []int{1}}, offers: []offer{{price: 15, goods: []int{7, 1, 4}}, {price: 7, goods: []int{6}}, {price: 19, goods: []int{3, 4}}}},
	{n: 1, m: 3, items: [][]int{[]int{2}}, offers: []offer{{price: 7, goods: []int{9, 6}}, {price: 15, goods: []int{7, 1}}, {price: 3, goods: []int{2, 8, 7}}}},
	{n: 2, m: 4, items: [][]int{[]int{7, 9, 3}, []int{5, 6}}, offers: []offer{{price: 3, goods: []int{2, 3}}, {price: 1, goods: []int{7, 3, 1}}, {price: 3, goods: []int{9, 1, 4}}, {price: 14, goods: []int{2, 9, 4}}}},
	{n: 2, m: 4, items: [][]int{[]int{7, 4, 8}, []int{7, 9}}, offers: []offer{{price: 2, goods: []int{8}}, {price: 17, goods: []int{8, 6}}, {price: 12, goods: []int{1, 2}}, {price: 1, goods: []int{4, 7, 8}}}},
	{n: 3, m: 3, items: [][]int{[]int{2, 5}, []int{2, 5, 4}, []int{9}}, offers: []offer{{price: 19, goods: []int{7}}, {price: 7, goods: []int{8}}, {price: 8, goods: []int{6}}}},
}

func solveCase(tc testCase) string {
	const nax = 1 << 9
	fr := make([]int, nax)
	for _, arr := range tc.items {
		mask := 0
		for _, x := range arr {
			mask |= 1 << (x - 1)
		}
		fr[mask]++
	}
	ile := make([]int, nax)
	for i := 0; i < nax; i++ {
		if fr[i] == 0 {
			continue
		}
		for x := 0; x < nax; x++ {
			if x&i == i {
				ile[x] += fr[i]
			}
		}
	}
	const INF = int(1e9 + 5)
	type pair struct{ first, second int }
	piz := make([]pair, nax)
	for i := range piz {
		piz[i] = pair{INF, INF}
	}
	bon := pair{INF, INF}
	for idx, off := range tc.offers {
		pr := off.price
		mask := 0
		for _, x := range off.goods {
			mask |= 1 << (x - 1)
		}
		if pr < piz[mask].first {
			if piz[mask].first != INF {
				old := piz[mask]
				if old.first < bon.first || (old.first == bon.first && old.second < bon.second) {
					bon = old
				}
			}
			piz[mask] = pair{pr, idx + 1}
		} else {
			cur := pair{pr, idx + 1}
			if cur.first < bon.first || (cur.first == bon.first && cur.second < bon.second) {
				bon = cur
			}
		}
	}
	ans := pair{INF, INF}
	var ret pair
	for i := 0; i < nax; i++ {
		if piz[i].first == INF {
			continue
		}
		for j := 0; j < nax; j++ {
			if i == j || piz[j].first == INF {
				continue
			}
			cm := i | j
			cov := ile[cm]
			sum := piz[i].first + piz[j].first
			nw := pair{-cov, sum}
			if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
				ans = nw
				ret = pair{piz[i].second, piz[j].second}
			}
		}
		if bon.first != INF {
			cov := ile[i]
			sum := piz[i].first + bon.first
			nw := pair{-cov, sum}
			if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
				ans = nw
				ret = pair{piz[i].second, bon.second}
			}
		}
	}
	return fmt.Sprintf("%d %d", ret.first, ret.second)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, arr := range tc.items {
			fmt.Fprintf(&sb, "%d", len(arr))
			for _, v := range arr {
				fmt.Fprintf(&sb, " %d", v)
			}
			sb.WriteByte('\n')
		}
		for _, off := range tc.offers {
			fmt.Fprintf(&sb, "%d %d", off.price, len(off.goods))
			for _, v := range off.goods {
				fmt.Fprintf(&sb, " %d", v)
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != want {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
