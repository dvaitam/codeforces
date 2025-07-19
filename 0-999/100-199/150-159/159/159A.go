package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   // read n and d
   if !scanner.Scan() {
       return
   }
   n, _ := strconv.Atoi(scanner.Text())
   scanner.Scan()
   d, _ := strconv.Atoi(scanner.Text())

   nameToID := make(map[string]int)
   var idToName []string
   // map for directed pair times
   pairTimes := make(map[int64][]int)
   // set for unordered pairs with mutual events
   outSet := make(map[int64]struct{})

   // helper to create key for pair
   key := func(u, v int) int64 {
       return (int64(u) << 32) | int64(uint32(v))
   }

   for i := 0; i < n; i++ {
       scanner.Scan()
       n1 := scanner.Text()
       scanner.Scan()
       n2 := scanner.Text()
       scanner.Scan()
       t, _ := strconv.Atoi(scanner.Text())
       // assign IDs
       p1, ok := nameToID[n1]
       if !ok {
           p1 = len(idToName)
           nameToID[n1] = p1
           idToName = append(idToName, n1)
       }
       p2, ok2 := nameToID[n2]
       if !ok2 {
           p2 = len(idToName)
           nameToID[n2] = p2
           idToName = append(idToName, n2)
       }
       // insert time into p1->p2
       k12 := key(p1, p2)
       times12 := pairTimes[k12]
       // find insertion index
       idx12 := sort.Search(len(times12), func(i int) bool { return times12[i] >= t })
       if idx12 == len(times12) {
           times12 = append(times12, t)
       } else if times12[idx12] != t {
           times12 = append(times12, 0)
           copy(times12[idx12+1:], times12[idx12:])
           times12[idx12] = t
       }
       pairTimes[k12] = times12
       // check in p2->p1 for any time >= t-d and != t
       k21 := key(p2, p1)
       times21 := pairTimes[k21]
       // lower_bound of t-d
       idx21 := sort.Search(len(times21), func(i int) bool { return times21[i] >= t-d })
       if idx21 < len(times21) && times21[idx21] != t {
           // record unordered pair
           u, v := p1, p2
           if u > v {
               u, v = v, u
           }
           outSet[key(u, v)] = struct{}{}
       }
   }

   // prepare output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // number of pairs
   fmt.Fprintln(writer, len(outSet))
   // collect and sort pairs
   type pair struct{ u, v int }
   pairs := make([]pair, 0, len(outSet))
   for kk := range outSet {
       u := int(kk >> 32)
       v := int(uint32(kk & 0xffffffff))
       pairs = append(pairs, pair{u, v})
   }
   sort.Slice(pairs, func(i, j int) bool {
       if pairs[i].u != pairs[j].u {
           return pairs[i].u < pairs[j].u
       }
       return pairs[i].v < pairs[j].v
   })
   for _, p := range pairs {
       fmt.Fprintln(writer, idToName[p.u], idToName[p.v])
   }
}
