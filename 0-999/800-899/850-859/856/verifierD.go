package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const embeddedRefCpp = `
#include <bits/stdc++.h>
using namespace std;
#define ll long long

template<typename T, bool maximum_mode = false>
struct RMQ {
    static int highest_bit(unsigned x) {
        return x == 0 ? -1 : 31 - __builtin_clz(x);
    }

    int n = 0;
    vector<T> values;
    vector<vector<int>> range_low;

    RMQ(const vector<T> &_values = {}) {
        if (!_values.empty())
            build(_values);
    }

    int better_index(int a, int b) const {
        return (maximum_mode ? values[b] < values[a] : values[a] < values[b]) ? a : b;
    }

    void build(const vector<T> &_values) {
        values = _values;
        n = int(values.size());
        int levels = highest_bit(n) + 1;
        range_low.resize(levels);

        for (int k = 0; k < levels; k++)
            range_low[k].resize(n - (1 << k) + 1);

        for (int i = 0; i < n; i++)
            range_low[0][i] = i;

        for (int k = 1; k < levels; k++)
            for (int i = 0; i <= n - (1 << k); i++)
                range_low[k][i] = better_index(range_low[k - 1][i], range_low[k - 1][i + (1 << (k - 1))]);
    }

    int query_index(int a, int b) const {
        assert(0 <= a && a < b && b <= n);
        int level = highest_bit(b - a);
        return better_index(range_low[level][a], range_low[level][b - (1 << level)]);
    }

    T query_value(int a, int b) const {
        return values[query_index(a, b)];
    }
};

struct LCA {
    int n = 0;
    vector<vector<int>> adj;
    vector<int> parent, depth, subtree_size;
    vector<int> euler, first_occurrence;
    vector<int> tour_start, tour_end;
    vector<int> tour_list;
    vector<int> heavy_root;
    vector<int> heavy_root_depth, heavy_root_parent;
    RMQ<int> rmq;
    bool built = false;

    LCA(int _n = 0) { init(_n); }
    LCA(const vector<vector<int>> &_adj) { init(_adj); }

    void init(int _n) {
        n = _n;
        adj.assign(n, {});
        parent.resize(n);
        depth.resize(n);
        subtree_size.resize(n);
        first_occurrence.resize(n);
        tour_start.resize(n);
        tour_end.resize(n);
        tour_list.resize(n);
        heavy_root.resize(n);
        built = false;
    }

    void init(const vector<vector<int>> &_adj) {
        init(int(_adj.size()));
        adj = _adj;
    }

    void add_edge(int a, int b) {
        adj[a].push_back(b);
        adj[b].push_back(a);
    }

    void dfs(int node, int par) {
        parent[node] = par;
        depth[node] = par < 0 ? 0 : depth[par] + 1;
        subtree_size[node] = 1;
        adj[node].erase(remove(adj[node].begin(), adj[node].end(), par), adj[node].end());
        for (int child : adj[node]) {
            dfs(child, node);
            subtree_size[node] += subtree_size[child];
        }
        sort(adj[node].begin(), adj[node].end(), [&](int a, int b) -> bool {
            return subtree_size[a] > subtree_size[b];
        });
    }

    int tour;

    void tour_dfs(int node, bool heavy) {
        heavy_root[node] = heavy ? heavy_root[parent[node]] : node;
        first_occurrence[node] = int(euler.size());
        euler.push_back(node);
        tour_list[tour] = node;
        tour_start[node] = tour++;
        bool heavy_child = true;
        for (int child : adj[node]) {
            tour_dfs(child, heavy_child);
            euler.push_back(node);
            heavy_child = false;
        }
        tour_end[node] = tour;
    }

    void build(const vector<int> &roots = {}, bool build_rmq = true) {
        depth.assign(n, -1);
        for (int root : roots)
            if (depth[root] < 0) dfs(root, -1);
        for (int i = 0; i < n; i++)
            if (depth[i] < 0) dfs(i, -1);
        tour = 0;
        euler.clear();
        euler.reserve(2 * n);
        for (int i = 0; i < n; i++)
            if (parent[i] < 0) {
                tour_dfs(i, false);
                euler.push_back(-1);
            }
        assert(int(euler.size()) == 2 * n);
        vector<int> euler_depth;
        euler_depth.reserve(euler.size());
        for (int node : euler)
            euler_depth.push_back(node < 0 ? node : depth[node]);
        if (build_rmq)
            rmq.build(euler_depth);
        euler_depth.clear();
        heavy_root_depth.resize(n);
        heavy_root_parent.resize(n);
        for (int i = 0; i < n; i++) {
            heavy_root_depth[i] = depth[heavy_root[i]];
            heavy_root_parent[i] = parent[heavy_root[i]];
        }
        built = true;
    }

    int get_lca(int a, int b) const {
        a = first_occurrence[a];
        b = first_occurrence[b];
        if (a > b) swap(a, b);
        return euler[rmq.query_index(a, b + 1)];
    }

    bool is_ancestor(int a, int b) const {
        return tour_start[a] <= tour_start[b] && tour_start[b] < tour_end[a];
    }
};

template<typename T>
struct fenwick_tree {
    static int highest_bit(unsigned x) {
        return x == 0 ? -1 : 31 - __builtin_clz(x);
    }

    int tree_n = 0;
    T tree_sum = T();
    vector<T> tree;

    fenwick_tree(int n = -1) {
        if (n >= 0) init(n);
    }

    void init(int n) {
        tree_n = n;
        tree_sum = T();
        tree.assign(tree_n + 1, T());
    }

    void update(int index, const T &change) {
        assert(0 <= index && index < tree_n);
        tree_sum += change;
        for (int i = index + 1; i <= tree_n; i += i & -i)
            tree[i] += change;
    }

    T query(int count) const {
        count = min(count, tree_n);
        T sum = T();
        for (int i = count; i > 0; i -= i & -i)
            sum += tree[i];
        return sum;
    }

    T query(int a, int b) const {
        return query(b) - query(a);
    }
};

const int N = 200001;
vector<int> ch[N], idx[N], euler;
vector<array<int, 3>> go_edges[N];
ll dp[N];

int main(){
    ios::sync_with_stdio(false); cin.tie(0);
    int n, m;
    cin >> n >> m;
    LCA lca(n+1);
    for(int i=2; i<=n; i++){
        int j;
        cin >> j;
        ch[j].push_back(i);
        lca.add_edge(i,j);
    }
    lca.build();
    while(m--){
        int u, v, c;
        cin >> u >> v >> c;
        go_edges[lca.get_lca(u, v)].push_back({u, v, c});
    }

    // Iterative DFS to build euler tour with idx
    {
        stack<int> stk;
        stk.push(1);
        vector<int> childpos(n+1, 0);
        // phase: 0 = pre-visit, 1 = post-visit
        vector<int> phase(n+1, 0);
        while (!stk.empty()) {
            int u = stk.top();
            if (phase[u] == 0) {
                idx[u].push_back(euler.size());
                euler.push_back(u);
                phase[u] = 1;
            }
            if (childpos[u] < (int)ch[u].size()) {
                int v = ch[u][childpos[u]++];
                stk.push(v);
            } else {
                idx[u].push_back(euler.size());
                euler.push_back(u);
                stk.pop();
            }
        }
    }

    fenwick_tree<ll> tree(2*n);

    // Iterative DFS2 for dp computation
    {
        stack<int> stk;
        stk.push(1);
        vector<int> childpos(n+1, 0);
        vector<bool> visited(n+1, false);
        while (!stk.empty()) {
            int u = stk.top();
            if (!visited[u]) {
                visited[u] = true;
                // Push children in reverse so leftmost is processed first
                for (int i = (int)ch[u].size()-1; i >= 0; i--) {
                    stk.push(ch[u][i]);
                }
            } else {
                stk.pop();
                int idx1 = idx[u][0];
                ll tot = 0;
                for (int v : ch[u]) {
                    tot += dp[v];
                }
                dp[u] = tot;
                for (auto &e : go_edges[u]) {
                    int v = e[0], w = e[1], c = e[2];
                    ll res = 0;
                    res += tree.query(idx1, idx[v][0]+1);
                    res += tree.query(idx1, idx[w][0]+1);
                    res += tot;
                    res += c;
                    dp[u] = max(dp[u], res);
                }
                tree.update(idx[u][0], tot - dp[u]);
                tree.update(idx[u][1], dp[u] - tot);
            }
        }
    }

    cout << dp[1] << '\n';
}
`

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %v", err)
	}
	ref := filepath.Join(wd, "refD.bin")
	cppPath := filepath.Join(wd, "refD.cpp")
	if err := os.WriteFile(cppPath, []byte(embeddedRefCpp), 0644); err != nil {
		return "", fmt.Errorf("write cpp: %v", err)
	}
	cmd := exec.Command("g++", "-std=c++17", "-O2", "-o", ref, cppPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference cpp: %v: %s", err, string(out))
	}
	return ref, nil
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(6) + 3
	m := rng.Intn(5)
	parents := make([]int, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		parents[i] = p
		if i > 2 {
			sb.WriteByte(' ')
		}
		if i == 2 {
			fmt.Fprintf(&sb, "%d", p)
		} else {
			fmt.Fprintf(&sb, "%d", p)
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		w := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
	}
	_ = parents
	return sb.String()
}

func runCaseD(bin, ref, input string) error {
	expect, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		if err := runCaseD(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
