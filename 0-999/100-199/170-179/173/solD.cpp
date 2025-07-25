#include <iostream>
#include <algorithm>
#include <cstdio>
#include <cstring>
#include <cstdlib>
#include <cmath>
#include <vector>
#include <queue>

using namespace std;

const int V = 100010;
const int E = 200020;
int head[V], cur[E], nxt[E], ne;
void init(int n) {
    ne = 0, memset(head, -1, sizeof(head));
}
inline void addEdge(int u, int v) {
    cur[ne] = v, nxt[ne] = head[u], head[u] = ne++;
}
int col[V], du[V];
bool vist[V];
int idx[V], tot;
int n, m;

void dfs(int u, int c) {
    if (col[u] == c) return;
    col[u] = c;
    for (int i = head[u]; i != -1; i = nxt[i]) {
        dfs(cur[i], 1 - c);
    }
}

void output() {
    puts("YES");
    for (int i = 1; i <= n; ++i) {
        if (!idx[i] && !col[i]) idx[i] = ++tot;
    }
    for (int i = 1; i <= n; ++i) {
        if (!idx[i] && col[i]) idx[i] = ++tot;
    }
    for (int i = 1; i <= n; ++i) {
        printf("%d ", (idx[i] - 1) / 3 + 1);
    }
}

void paint(int x) {
    idx[x] = ++tot;
    memset(vist, false, sizeof(vist));
    for (int i = head[x]; i != -1; i = nxt[i]) vist[cur[i]] = true;
    int t = 0;
    for (int i = 1; i <= n && t < 2; ++i) {
        if (!idx[i] && col[i] != col[x] && !vist[i]) {
            idx[i] = ++tot; ++t;
        }
    }
}

int main() {

    scanf("%d %d", &n, &m);
    init(n + 1);
    for (int i = 0; i < m; ++i) {
        int a, b;
        scanf("%d %d", &a, &b);
        addEdge(a, b);
        addEdge(b, a);
    }
    memset(col, -1, sizeof(col));
    for (int i = 1; i <= n; ++i) {
        if (col[i] == -1) {
            dfs(i, 0);
        }
    }
    int cnt[2] = {0, 0};
    for (int i = 1; i <= n; ++i) {
        ++cnt[col[i]];
    }
    if (cnt[0] % 3 == 0) {
        output();
        return 0;
    }
    if (cnt[0] % 3 == 2) {
        for (int i = 1; i <= n; ++i) {
            col[i] = 1 - col[i];
        }
        swap(cnt[0], cnt[1]);
    }

    for (int u = 1; u <= n; ++u) {
        for (int i = head[u]; i != -1; i = nxt[i]) {
            ++du[u];
        }
    }

    for (int u = 1; u <= n; ++u) {
        if (col[u] != 0) continue;
        if (du[u] + 2 <= cnt[1]) {
            paint(u); output(); return 0;
        }
    }
    int t = 0;
    for (int u = 1; u <= n; ++u) {
        if (col[u] == 1 && du[u] + 2 <= cnt[0]) {
            paint(u); ++t;
            if (t == 2) {
                output(); return 0;
            }
        }
    }
    puts("NO");
    return 0;
}