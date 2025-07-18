#include <map>
#include <set>
#include <cmath>
#include <ctime> 
#include <queue>
#include <stack>
#include <bitset> 
#include <cstdio>
#include <cstdlib>
#include <climits>
#include <cstring>
#include <iostream>
#include <algorithm>

#define fi first
#define se second
#define pb push_back
#define MP std::make_pair
#define PII std::pair<int, int>
#define all(x) (x).begin(), (x).end()
#define CL(a, b) memset(a, b, sizeof a)
#define rep(i, l, r) for (int i = (l); i <= (r); ++ i)
#define per(i, r, l) for (int i = (r); i >= (l); -- i)
#define PE(x, a) for (int x = head[a]; x;x = edge[x].next)

typedef long long ll;

template <class T>
inline void rd(T &x) {
    char c = getchar(), f = 0; x = 0;
    while (!isdigit(c)) f = (c == '-'), c = getchar();
    while (isdigit(c)) x = x * 10 + c - '0', c = getchar();
    x = f ? -x : x;
}

const int MAXN = 4e5 + 7;

struct Edge {
    int t, next;
} edge[MAXN << 1];

int head[MAXN], cnt, f[MAXN], deg[MAXN];

int n, m, ord[MAXN];
PII E[MAXN];
bool ty[MAXN];

void add(int u, int v) {
    edge[++cnt] = (Edge){v, head[u]}; head[u] = cnt;
}

void solve() {
    cnt = 0;
    rd(n), rd(m);
    rep(i, 0, n) head[i] = deg[i] = ord[i] = 0, f[i] = i;
    bool fl = 1;
    rep(i, 1, m) {
        rd(ty[i]);
        rd(E[i].fi), rd(E[i].se);
        if (ty[i]) add(E[i].fi, E[i].se), deg[E[i].se]++;
    }
    std::queue<int> q;
    rep(i, 1, n) if (!deg[i]) q.push(i);
    int cnt = 0;
    while (!q.empty()) {
        int u = q.front(); q.pop();
        ord[u] = ++cnt;
        PE(e, u) {
            int v = edge[e].t;
            deg[v]--;
            if (!deg[v]) q.push(v);
        }
    }
    // rep(i, 1, n) printf("i = %d, ord = %d\n", i, ord[i]);
    if (cnt != n) puts("NO");
    else {
        puts("YES");
        rep(i, 1, m) if (!ty[i]) {
            if (ord[E[i].fi] < ord[E[i].se]) printf("%d %d\n", E[i].fi, E[i].se);
            else printf("%d %d\n", E[i].se, E[i].fi);
        } else printf("%d %d\n", E[i].fi, E[i].se);
    }
}

int main() {
    int t; rd(t);
    while (t--) solve();
    return 0;
}