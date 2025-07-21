package main

import (
   "bufio"
   "fmt"
   "os"
)

// Occurrence positions of tag in queries
type occ struct {
   qid      int
   positions []int // positions in descending order
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read document line
   doc, err := reader.ReadString('\n')
   if err != nil && len(doc) == 0 {
       return
   }
   // remove newline
   if len(doc) > 0 && doc[len(doc)-1] == '\n' {
       doc = doc[:len(doc)-1]
   }
   // read m
   var m int
   fmt.Fscan(reader, &m)
   // read queries
   queries := make([][]string, m)
   lengths := make([]int, m)
   // mapping tag-> per-query positions
   tagMap := make(map[string]map[int][]int)
   scanner := bufio.NewScanner(reader)
   for i := 0; i < m; i++ {
       scanner.Scan()
       line := scanner.Text()
       parts := bufio.NewScanner(bufio.NewReaderString(line))
       parts.Split(bufio.ScanWords)
       for parts.Scan() {
           token := parts.Text()
           queries[i] = append(queries[i], token)
           if tagMap[token] == nil {
               tagMap[token] = make(map[int][]int)
           }
           tagMap[token][i] = append(tagMap[token][i], len(queries[i]))
       }
       lengths[i] = len(queries[i])
   }
   // build occurrences list per tag
   occMap := make(map[string][]occ)
   for tag, m1 := range tagMap {
       tmp := make([]occ, 0, len(m1))
       for qi, ps := range m1 {
           // sort positions descending
           for l, r := 0, len(ps)-1; l < r; l, r = l+1, r-1 {
               ps[l], ps[r] = ps[r], ps[l]
           }
           tmp = append(tmp, occ{qid: qi, positions: ps})
       }
       occMap[tag] = tmp
   }
   // dp for each query: dp[q][k] means prefix k matched
   dp := make([][]bool, m)
   for i := 0; i < m; i++ {
       dp[i] = make([]bool, lengths[i]+1)
       dp[i][0] = true
   }
   ans := make([]int, m)
   // stack for changes count per element
   var changeStack []int
   // change records
   type change struct{ qid, pos int }
   var changes []change
   // process doc tags
   n := len(doc)
   for i := 0; i < n; {
       if doc[i] != '<' {
           i++
           continue
       }
       j := i + 1
       isClose := false
       if j < n && doc[j] == '/' {
           isClose = true
           j++
       }
       // find '>'
       k := j
       selfClose := false
       for k < n && doc[k] != '>' {
           k++
       }
       content := doc[j:k]
       // detect self-closing
       if !isClose && len(content) > 0 && content[len(content)-1] == '/' {
           selfClose = true
           content = content[:len(content)-1]
       }
       tag := string(content)
       if isClose {
           // closing tag: revert
           last := changeStack[len(changeStack)-1]
           changeStack = changeStack[:len(changeStack)-1]
           for len(changes) > last {
               c := changes[len(changes)-1]
               dp[c.qid][c.pos] = false
               changes = changes[:len(changes)-1]
           }
       } else {
           // opening or self-closing: enter
           prev := len(changes)
           changeStack = append(changeStack, prev)
           // process dp updates
           if occs, ok := occMap[tag]; ok {
               for _, oc := range occs {
                   q := oc.qid
                   // for each position in descending order
                   for _, posi := range oc.positions {
                       if !dp[q][posi] && dp[q][posi-1] {
                           dp[q][posi] = true
                           changes = append(changes, change{qid: q, pos: posi})
                           // count only when complete and matching here
                           if posi == lengths[q] {
                               ans[q]++
                           }
                       }
                   }
               }
           }
           // exit if self-closing
           if selfClose {
               last := changeStack[len(changeStack)-1]
               changeStack = changeStack[:len(changeStack)-1]
               for len(changes) > last {
                   c := changes[len(changes)-1]
                   dp[c.qid][c.pos] = false
                   changes = changes[:len(changes)-1]
               }
           }
       }
       i = k + 1
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
