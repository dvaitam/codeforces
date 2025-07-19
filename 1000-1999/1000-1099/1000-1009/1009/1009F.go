package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n       int
   reader  *bufio.Reader
   writer  *bufio.Writer
   last    []int
   eNext   []int
   eTo     []int
   ecnt    int
   mem     []int32
   pos     []int
   ln      []int
   mx      []int32
   mxid    []int
   ans     []int
)

func in() int {
   var x int
   c, err := reader.ReadByte()
   for err == nil && (c < '0' || c > '9') {
       c, err = reader.ReadByte()
   }
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   return x
}

func addEdge(a, b int) {
   ecnt++
   eNext[ecnt] = last[a]
   eTo[ecnt] = b
   last[a] = ecnt
}

func merge(x, y int) {
   if ln[x] < ln[y] {
       pos[x], pos[y] = pos[y], pos[x]
       ln[x], ln[y] = ln[y], ln[x]
       mx[x], mx[y] = mx[y], mx[x]
       mxid[x], mxid[y] = mxid[y], mxid[x]
   }
   baseX := pos[x]
   baseY := pos[y]
   for i := 0; i < ln[y]; i++ {
       mem[baseX+i] += mem[baseY+i]
       if mem[baseX+i] > mx[x] {
           mx[x] = mem[baseX+i]
           mxid[x] = i
       } else if mem[baseX+i] == mx[x] && i < mxid[x] {
           mxid[x] = i
       }
   }
}

func dfs(x, fa int) {
   mx[x] = 1
   mxid[x] = 0
   ln[x] = 1
   mem[pos[x]] = 1
   for i := last[x]; i != 0; i = eNext[i] {
       y := eTo[i]
       if y == fa {
           continue
       }
       pos[y] = pos[x] + ln[x]
       dfs(y, x)
       ln[y]++
       mem[pos[y]] = 0
       mxid[y]++
       merge(x, y)
   }
   ans[x] = mxid[x]
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   n = in()
   // allocate based on n
   last = make([]int, n+2)
   eNext = make([]int, (n*2)+2)
   eTo = make([]int, (n*2)+2)
   mem = make([]int32, (n*23)+2)
   pos = make([]int, n+2)
   ln = make([]int, n+2)
   mx = make([]int32, n+2)
   mxid = make([]int, n+2)
   ans = make([]int, n+2)
   for i := 1; i < n; i++ {
       a := in()
       b := in()
       addEdge(a, b)
       addEdge(b, a)
   }
   pos[1] = 0
   dfs(1, 0)
   for i := 1; i <= n; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
