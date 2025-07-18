#include <cmath>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <iostream>
#include <vector>
#include <queue>
#include <map>
#include <set>

inline int read() {
    char c = getchar(); int x = 0;
    while (c < '0' || c > '9') { c = getchar(); }
    while (c >= '0' && c <= '9') { x = (x << 1) + (x << 3) + (c & 15); c = getchar(); }
    return x;
}

const int maxN = 2000005;

int n, q, vol, a[maxN], l[maxN], r[maxN], stack[maxN];
long long ans[maxN], sm[maxN];
std::pair<std::pair<int, int>, int> s[maxN];
std::vector<int> tag[maxN];

struct TreeArray {
    long long f[maxN];

    inline void modify(int u, int x) { for (u++; u <= n + 5; u += u & -u) { f[u] += x; } }
    inline long long query(int u) { long long res = 0; for (u++; u; u ^= u & -u) { res += f[u]; } return res; }
} f, g, h, cnt, sum;

int main() {
    n = read(); q = read(); a[0] = a[n + 1] = n + 1;
    for (int i = 1; i <= n; i++) { a[i] = read(); }
    stack[vol = 1] = 0;
    for (int i = 1; i <= n; i++) {
        while (a[stack[vol]] <= a[i]) { vol--; }
        l[i] = stack[vol]; stack[++vol] = i;
    }
    stack[vol = 1] = n + 1;
    for (int i = n; i; i--) {
        while (a[stack[vol]] <= a[i]) { vol--; }
        r[i] = stack[vol]; stack[++vol] = i;
    }
    for (int i = 1; i <= q; i++) { s[i].first.second = read(); }
    for (int i = 1; i <= q; i++) { s[i].first.first = read(); s[i].second = i; }
    std::sort(s + 1, s + q + 1);
    for (int i = 1, j = 1; i <= n; i++) {
        for (int k = 0, u; k < (int) tag[i].size(); k++) {
            u = tag[i][k];
            f.modify(u, l[u] + 1); g.modify(u, -1); h.modify(u, i - l[u] - 1);
        }
        f.modify(i, -l[i] - 1); g.modify(i, 1); tag[r[i]].push_back(i); sm[i] += sm[i - 1] + l[i];
        cnt.modify(l[i], 1); sum.modify(l[i], l[i]);
        for (; j <= q && s[j].first.first == i; j++) {
            ans[s[j].second] += h.query(i) - h.query(s[j].first.second - 1) + f.query(i) - f.query(s[j].first.second - 1);
            ans[s[j].second] += (g.query(i) - g.query(s[j].first.second - 1)) * (i + 1);
            ans[s[j].second] += sum.query(s[j].first.second - 1) - sm[s[j].first.second - 1];
            ans[s[j].second] -= (cnt.query(s[j].first.second - 1) - (s[j].first.second - 1)) * (s[j].first.second - 1);
        }
    }
    for (int i = 1; i <= q; i++) { printf("%I64d ", ans[i]); }
    return 0;
}