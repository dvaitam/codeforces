#include <iostream>
#include <algorithm>
#include <cstdio>
#include <cstring>
#include <cstdlib>
#include <cmath>
#include <vector>
#include <map>
#include <cmath>
#include <set>

using namespace std;

typedef pair<int, int> PII;
#define mp make_pair
const int V = 1024;
const int E = 200020;
int head[V], cur[E], nxt[E], ne;
void init(int n) {
    ne = 0, memset(head, -1, sizeof(head));
}
inline void addEdge(int u, int v) {
    cur[ne] = v, nxt[ne] = head[u], head[u] = ne++;
}
int dfn[V], low[V], idx[V];
set<PII> s;
void tarjan(int u, int fa, int d) {
    dfn[u] = low[u] = d;
    for (int i = head[u]; i != -1; i = nxt[i]) {
        int v = cur[i];
        if (dfn[v] == -1) {
            tarjan(v, u, d + 1);
            low[u] = min(low[u], low[v]);
            if (dfn[u] < low[v]) {
                s.insert(mp(u, v));
                s.insert(mp(v, u));
            }
        } else if (v != fa) {
            low[u] = min(low[u], dfn[v]);
        }
    }
}
void dfs(int u, int c) {
    if (idx[u] != -1) return;
    idx[u] = c;
    for (int i = head[u]; i != -1; i = nxt[i]) {
        int v = cur[i];
        if (s.count(mp(u, v))) continue;
        dfs(v, c);
    }
}

vector<vector<int> >g;
vector<vector<int> >no;
bool sg[V][V];

bool output(int a, int b) {
    for (int i = 0; i < no[a].size(); ++i) {
        for (int j = 0; j < no[b].size(); ++j) {
            if (!sg[no[a][i]][no[b][j]]) {
                printf("%d %d\n", no[a][i], no[b][j]);
                return true;
            }
        }
    }
    //puts("error");
    return false;
}

pair<int, int> gao(int u, int fa) {
    vector<int> a, b;
    for (int i = 0; i < g[u].size(); ++i) {
        int  v = g[u][i];
        if (v == fa) continue;
        pair<int, int> ret = gao(v, u);
        b.push_back(ret.first);
        if (ret.second != -1) {
            a.push_back(ret.second);
            swap(b[0], b[b.size() - 1]);
        }
    }
    if (b.size() == 0) return make_pair(u, -1);
    for (int i = b.size() - 1; i >= 0; --i) {
        a.push_back(b[i]);
    }
    for (int i = 2 - a.size() % 2; i < a.size(); i += 2) {
        output(a[i], a[i + 1]);
    }
    if (fa == -1) output(a[0], a[1]);
    return a.size() & 1 ? make_pair(a[0], -1) : make_pair(a[0], a[1]);
}

int main() {
    int n, m;
    scanf("%d %d", &n, &m);
    init(n + 1);
    for (int i = 0; i < m; ++i) {
        int u, v;
        scanf("%d %d", &u, &v);
        addEdge(u, v);
        addEdge(v, u);
        sg[u][v] = sg[v][u] = true;
    }
    if (n == 2) {
        puts("-1");
        return 0;
    }
    memset(dfn, -1, sizeof(dfn));
    memset(idx, -1, sizeof(idx));
    tarjan(1, 0, 0);
    int cid = 0;
    for (int i = 1; i <= n; ++i) {
        if (idx[i] == -1) {
            dfs(i, cid); ++cid;
        }
    }

    g.clear(); g.resize(cid);
    no.clear(); no.resize(cid);
    for (int u = 1; u <= n; ++u) {
        for (int i = head[u]; i != -1; i = nxt[i]) {
            if (idx[u] != idx[cur[i]]) {
                g[idx[u]].push_back(idx[cur[i]]);
                //g[idx[cur[i]]].push_back(idx[u]);
            }
        }
    }

    for (int i = 1; i <= n; ++i) {
        no[idx[i]].push_back(i);
    }
    if (cid == 1) {
        puts("0");
    } else if (cid == 2) {
        puts("1");
        for (int i = 1; i <= n; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (!sg[i][j] && idx[i] != idx[j]) {
                    printf("%d %d\n", i, j);
                    return 0;
                }
            }
        }
    } else {
        int leaf = 0, st = -1;
        for (int i = 0; i < cid; ++i) {
            if (g[i].size() == 1) ++leaf;
            else if (g[i].size() > 1 && st == -1) st = i;
        }
        printf("%d\n", (leaf + 1) / 2);
        gao(st, -1);
    }
    return 0;
}