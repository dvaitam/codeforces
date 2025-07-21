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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // Read floors: input line i is floor n-i+1
   floors := make([][]byte, n)
   for i := 0; i < n; i++ {
       // read line
       var line string
       fmt.Fscan(reader, &line)
       floors[n-1-i] = []byte(line)
   }

   // visited positions for current floor: use versioning to avoid O(m) resets
   visited := make([][2]int, m)
   curVer := 1

   // state
   floor := n - 1
   pos := 0
   dir := 1 // +1 right, -1 left
   var time int64

   // reset visited for a new floor or after breaking
   resetVisited := func() {
       curVer++
   }
   // initial version
   resetVisited()

   for {
       // reached first floor?
       if floor == 0 {
           fmt.Fprint(writer, time)
           return
       }
       // fall if below is empty
       if floors[floor-1][pos] == '.' {
           floor--
           time++
           resetVisited()
           continue
       }
       // horizontal move
       next := pos + dir
       var cell byte
       if next < 0 || next >= m {
           cell = '#'
       } else {
           cell = floors[floor][next]
       }
       if cell == '+' {
           // break brick
           floors[floor][next] = '.'
           dir = -dir
           time++
           resetVisited()
           continue
       }
       if cell == '.' {
           pos = next
           // dir unchanged
       } else if cell == '#' {
           // just turn
           dir = -dir
       }
       // time for horizontal non-breaking move
       time++
       // detect loop: state pos,dir
       di := 1
       if dir < 0 {
           di = 0
       }
       if visited[pos][di] == curVer {
           fmt.Fprint(writer, "Never")
           return
       }
       visited[pos][di] = curVer
   }
}
