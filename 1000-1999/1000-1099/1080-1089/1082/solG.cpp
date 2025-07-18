#define _CRT_SECURE_NO_DEPRECATE
#define _USE_MATH_DEFINES
#include <iostream>
#include <cmath>
#include <string>
#include <vector>
#include <algorithm>
#include <cstdio>
#include <memory.h>
#include <set>
#include <ctime>
#include <map>
#include <cstring>
#include <iterator>
#include <queue>
#include <assert.h>
#include <bitset>
#include <complex>
#include <unordered_map>

//#pragma comment(linker, "/STACK:512000000")

using namespace std;

#define pb emplace_back
#define mp make_pair
#define all(a) (a).begin(), (a).end()
#define forn(i, n) for (int i = 0; i < (n); ++i)
#define forab(i, a, b) for (int i = (a); i < (b); ++i)

typedef long long ll;
typedef long double ld;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;

const int infi = 1e9 + 7;
const ll infl = (ll)1e18 + (ll)7;

struct edge {
    int to, flow, max_flow;
    edge() {}
    edge(int _to, int _max_flow) : to(_to), flow(0), max_flow(_max_flow) {}
};
vector<edge> edges;
vector<vector<int> > g;
void add_edge(int u, int v, int max_flow, int inv = 0) {
    g[u].pb(edges.size());
    edges.pb(v, max_flow);
    g[v].pb(edges.size());
    edges.pb(u, inv);
}
int dist[3010];
int que[3010];
int S, T;
int bfs(int n) {
    forn(i, n)
        dist[i] = -1;
    dist[S] = 0;
    que[0] = S;
    int l = 0, r = 1;
    while (l < r) {
        int v = que[l++];
        int d = dist[v] + 1;
        if (v == T)
            break;
        for (int i : g[v]) {
            edge& e = edges[i];
            if (e.flow == e.max_flow)
                continue;
            if (dist[e.to] != -1)
                continue;
            dist[e.to] = d;
            que[r++] = e.to;
        }
    }
    return dist[T];
}
int start[3010];
int dfs(int v, int flow) {
    if (v == T)
        return flow;
    int d = dist[v] + 1;
    for (; start[v] < g[v].size(); ++start[v]) {
        int i = g[v][start[v]];
        edge& e = edges[i];
        if (e.flow == e.max_flow || dist[e.to] != d)
            continue;
        int a = dfs(e.to, min(flow, e.max_flow - e.flow));
        if (!a)
            continue;
        e.flow += a;
        edges[i ^ 1].flow -= a;
        return a;
    }
    return 0;
}
ll dinic(int n) {
    ll ans = 0;
    while (bfs(n) != -1) {
        
        memset(start, 0, sizeof(start));
        int a = dfs(S, infi);

        while (a) {
            ans += a;
            a = dfs(S, infi);

        }
    }
    return ans;
}

int main() {
    cin.sync_with_stdio(false);
    cin.tie(0);
  //  freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);
    //freopen("customs.in", "r", stdin); freopen("customs.out", "w", stdout);
    int n, m;
    cin >> n >> m;
    S = n + m;
    T = S + 1;
    g.resize(T + 1);
    forn(i, n) {
        int a;
        cin >> a;
        add_edge(i, T, a);
    }
    ll ans = 0;
    forn(i, m) {
        int u, v, c;
        cin >> u >> v >> c;
        --u, --v;
        ans += c;
        add_edge(u, n + i, c, c);
        add_edge(v, n + i, c, c);
        add_edge(S, n + i, c);
    }
    ll q = dinic(T + 1);
    cout << ans - q << '\n';

    return 0;
}