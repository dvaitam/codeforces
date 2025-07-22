package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// key identifies a project by name and version
// used as map key

type key struct {
	name string
	ver  int
}

type item struct {
	k    key
	dist int
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	deps := make(map[key][]key)
	var root key
	for i := 0; i < n; i++ {
		var name string
		var version int
		fmt.Fscan(reader, &name, &version)
		if i == 0 {
			root = key{name, version}
		}
		var m int
		fmt.Fscan(reader, &m)
		arr := make([]key, m)
		for j := 0; j < m; j++ {
			var dn string
			var dv int
			fmt.Fscan(reader, &dn, &dv)
			arr[j] = key{dn, dv}
		}
		deps[key{name, version}] = arr
	}

	// BFS with tie-breaking: smallest distance, then largest version
	best := make(map[string]struct{ dist, ver int })
	best[root.name] = struct{ dist, ver int }{0, root.ver}
	q := []item{{root, 0}}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		b := best[cur.k.name]
		if b.dist != cur.dist || b.ver != cur.k.ver {
			continue // outdated entry
		}
		for _, nxt := range deps[cur.k] {
			nd := cur.dist + 1
			nb, ok := best[nxt.name]
			if !ok || nd < nb.dist || (nd == nb.dist && nxt.ver > nb.ver) {
				best[nxt.name] = struct{ dist, ver int }{nd, nxt.ver}
				q = append(q, item{nxt, nd})
			}
		}
	}

	// collect result excluding root
	res := make([]key, 0, len(best)-1)
	for name, info := range best {
		if name == root.name && info.ver == root.ver && info.dist == 0 {
			continue
		}
		res = append(res, key{name, info.ver})
	}

	// sort lexicographically by name
	sort.Slice(res, func(i, j int) bool {
		return res[i].name < res[j].name
	})

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for _, k := range res {
		fmt.Fprintf(writer, "%s %d\n", k.name, k.ver)
	}
}
