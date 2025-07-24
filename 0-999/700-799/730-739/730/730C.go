package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// BIT maintains prefix sums of shovel counts and costs by price index
type BIT struct {
   n       int
   bitKi   []int64
   bitCost []int64
   prices  []int64
}

// NewBIT initializes a BIT with given 1-based prices slice (prices[0] unused)
func NewBIT(prices []int64) *BIT {
   n := len(prices) - 1
   return &BIT{
       n:       n,
       bitKi:   make([]int64, n+1),
       bitCost: make([]int64, n+1),
       prices:  prices,
   }
}

// Add inserts deltaKi and deltaCost at index idx
func (b *BIT) Add(idx int, deltaKi, deltaCost int64) {
   for i := idx; i <= b.n; i += i & -i {
       b.bitKi[i] += deltaKi
       b.bitCost[i] += deltaCost
   }
}

// sumKi returns prefix sum of ki up to idx
func (b *BIT) sumKi(idx int) int64 {
   var s int64
   for i := idx; i > 0; i -= i & -i {
       s += b.bitKi[i]
   }
   return s
}

// sumCost returns prefix sum of cost up to idx
func (b *BIT) sumCost(idx int) int64 {
   var s int64
   for i := idx; i > 0; i -= i & -i {
       s += b.bitCost[i]
   }
   return s
}

// FindPriceIdxByCumulativeKi finds smallest idx such that sumKi(idx) >= target
func (b *BIT) FindPriceIdxByCumulativeKi(target int64) int {
   idx := 0
   var acc int64
   // compute highest power of two <= n
   bitMask := 1
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for mask := bitMask; mask > 0; mask >>= 1 {
       nxt := idx + mask
       if nxt <= b.n && acc+b.bitKi[nxt] < target {
           idx = nxt
           acc += b.bitKi[nxt]
       }
   }
   return idx + 1
}

// CalcCost returns minimum cost to buy target shovels
func (b *BIT) CalcCost(target int64) int64 {
   idx := b.FindPriceIdxByCumulativeKi(target)
   kiBefore := b.sumKi(idx - 1)
   costBefore := b.sumCost(idx - 1)
   need := target - kiBefore
   return costBefore + need*b.prices[idx]
}

// Store holds information about a Bulmart store
type Store struct {
   city     int
   ki       int64
   price    int64
   priceIdx int
}

func main() {
   defer writer.Flush()
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for e := 0; e < m; e++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   var w int
   fmt.Fscan(reader, &w)
   stores := make([]Store, w)
   prices := make([]int64, w)
   for i := 0; i < w; i++ {
       var c, ki, p int64
       fmt.Fscan(reader, &c, &ki, &p)
       stores[i].city = int(c)
       stores[i].ki = ki
       stores[i].price = p
       prices[i] = p
   }
   // compress prices
   sort.Slice(prices, func(i, j int) bool { return prices[i] < prices[j] })
   uniq := prices[:0]
   for _, v := range prices {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   // build 1-based price list
   priceList := make([]int64, len(uniq)+1)
   for i, v := range uniq {
       priceList[i+1] = v
   }
   // map each store to price index
   priceIndex := make(map[int64]int, len(uniq))
   for i, v := range uniq {
       priceIndex[v] = i + 1
   }
   for i := range stores {
       stores[i].priceIdx = priceIndex[stores[i].price]
   }
   // read queries
   var q int
   fmt.Fscan(reader, &q)
   gj := make([]int, q)
   rj := make([]int64, q)
   aj := make([]int64, q)
   queriesByCity := make([][]int, n+1)
   for j := 0; j < q; j++ {
       fmt.Fscan(reader, &gj[j], &rj[j], &aj[j])
       queriesByCity[gj[j]] = append(queriesByCity[gj[j]], j)
   }
   // answers
   ans := make([]int, q)
   for i := range ans {
       ans[i] = -1
   }
   // process queries per starting city
   for city := 1; city <= n; city++ {
       qs := queriesByCity[city]
       if len(qs) == 0 {
           continue
       }
       // BFS for distances
       dist := make([]int, n+1)
       for i := 1; i <= n; i++ {
           dist[i] = -1
       }
       queue := make([]int, 0, n)
       dist[city] = 0
       queue = append(queue, city)
       for bi := 0; bi < len(queue); bi++ {
           u := queue[bi]
           for _, v := range adj[u] {
               if dist[v] == -1 {
                   dist[v] = dist[u] + 1
                   queue = append(queue, v)
               }
           }
       }
       // bucket stores by distance
       buckets := make([][]int, n+1)
       var totalReachableKi int64
       for i, st := range stores {
           d := dist[st.city]
           if d >= 0 {
               buckets[d] = append(buckets[d], i)
               totalReachableKi += st.ki
           }
       }
       // track done queries
       doneCount := 0
       totalQ := len(qs)
       done := make([]bool, len(qs))
       // mark impossible queries
       for idx, j := range qs {
           if totalReachableKi < rj[j] {
               done[idx] = true
               // ans[j] remains -1
               doneCount++
           }
       }
       // initialize BIT
       bit := NewBIT(priceList)
       var totalKi int64
       // sweep distances
       for d := 0; d <= n && doneCount < totalQ; d++ {
           for _, si := range buckets[d] {
               st := stores[si]
               bit.Add(st.priceIdx, st.ki, st.ki*st.price)
               totalKi += st.ki
           }
           // check queries
           for idx, j := range qs {
               if done[idx] {
                   continue
               }
               if totalKi < rj[j] {
                   continue
               }
               cost := bit.CalcCost(rj[j])
               if cost <= aj[j] {
                   ans[j] = d
                   done[idx] = true
                   doneCount++
               }
           }
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
