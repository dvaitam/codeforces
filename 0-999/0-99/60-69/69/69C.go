package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // k allies, n basic, m composite, q purchases
   var k, n, m, q int
   if _, err := fmt.Fscan(reader, &k, &n, &m, &q); err != nil {
       return
   }
   // basic artifact names
   basicNames := make([]string, n)
   nameToBasic := make(map[string]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &basicNames[i])
       nameToBasic[basicNames[i]] = i
   }
   // consume end of line before reading composite descriptions
   reader.ReadString('\n')
   // composite artifacts: name and requirements of basics
   compNames := make([]string, m)
   compReq := make([][]int, m)
   for i := 0; i < m; i++ {
       line, _ := reader.ReadString('\n')
       line = strings.TrimSpace(line)
       parts := strings.SplitN(line, ": ", 2)
       compNames[i] = parts[0]
       rest := parts[1]
       tokens := strings.Split(rest, ", ")
       req := make([]int, n)
       for _, tok := range tokens {
           sp := strings.SplitN(tok, " ", 2)
           cnt, _ := strconv.Atoi(sp[1])
           idx := nameToBasic[sp[0]]
           req[idx] = cnt
       }
       compReq[i] = req
   }
   // inventory counts per hero
   basicCnt := make([][]int, k)
   compCnt := make([][]int, k)
   for i := 0; i < k; i++ {
       basicCnt[i] = make([]int, n)
       compCnt[i] = make([]int, m)
   }
   // process purchases
   for i := 0; i < q; i++ {
       var ai int
       var bname string
       fmt.Fscan(reader, &ai, &bname)
       h := ai - 1
       bi := nameToBasic[bname]
       basicCnt[h][bi]++
       // check for composite crafting (at most one)
       for j := 0; j < m; j++ {
           can := true
           for idx, need := range compReq[j] {
               if need > 0 && basicCnt[h][idx] < need {
                   can = false
                   break
               }
           }
           if can {
               // consume basics
               for idx, need := range compReq[j] {
                   if need > 0 {
                       basicCnt[h][idx] -= need
                   }
               }
               // add composite
               compCnt[h][j]++
               break
           }
       }
   }
   // output results
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   type pair struct{ name string; cnt int }
   for i := 0; i < k; i++ {
       var lst []pair
       for bi := 0; bi < n; bi++ {
           if basicCnt[i][bi] > 0 {
               lst = append(lst, pair{basicNames[bi], basicCnt[i][bi]})
           }
       }
       for cj := 0; cj < m; cj++ {
           if compCnt[i][cj] > 0 {
               lst = append(lst, pair{compNames[cj], compCnt[i][cj]})
           }
       }
       sort.Slice(lst, func(a, b int) bool { return lst[a].name < lst[b].name })
       fmt.Fprintln(writer, len(lst))
       for _, p := range lst {
           fmt.Fprintln(writer, p.name, p.cnt)
       }
   }
}
