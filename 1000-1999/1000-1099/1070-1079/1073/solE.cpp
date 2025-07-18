#include <bits/stdc++.h>

using namespace std;

const int N = 25;
const int K = 10;
const int MOD = 998244353;

inline int add(int u, int v) {
    u += v;
    if (u >= MOD) u -= MOD;
    return u;
}

inline int sub(int u, int v) {
    u -= v;
    if (u < 0) u += MOD;
    return u;
}

inline int mul(int u, int v) {
    return (long long)u * v % MOD;
}

int k;
int f[N][1 << K];
int g[N][1 << K];
int p[N];

void init() {
    p[0] = 1;
    for (int i = 1; i < N; i++) {
        p[i] = mul(p[i - 1], 10);
    }
    g[0][0] = 1;
    for (int i = 1; i < N; i++) {
        for (int j = 0; j <= 9; j++) {
            for (int mask = 0; mask < (1 << K); mask++) {
                int newMask = mask | (1 << j);
                g[i][newMask] = add(g[i][newMask], g[i - 1][mask]);
                int foo = mul(j, p[i - 1]);
                foo = mul(foo, g[i - 1][mask]);
                foo = add(foo, f[i - 1][mask]);
                f[i][newMask] = add(f[i][newMask], foo);
            }
        }
    }
}

int get(long long val) {
    vector<int> v;
    while (val > 0) {
        v.push_back(val % 10);
        val /= 10;
    }
    // for (int u : v) {
    //     cout << u << " ";
    // }
    // cout << endl;
    int res = 0;
    int curMask = 0;
    int tot = 0;
    int sz = (int)v.size();

    for (int i = 1; i < sz; i++) {
        for (int j = 1; j <= 9; j++) {
            for (int mask = 0; mask < (1 << K); mask++) {
                int newMask = mask | (1 << j);
                int cnt = __builtin_popcount(newMask);
                if (cnt > k) continue;
                int foo = mul(j, p[i - 1]);
                foo = mul(foo, g[i - 1][mask]);
                foo = add(foo, f[i - 1][mask]);
                res = add(res, foo);
            }
        }
    }

    for (int i = sz; i > 0; i--) {
        int u = v[i - 1];
        for (int j = (i == sz ? 1 : 0); j < u; j++) {
            int newMask = curMask | (1 << j);
            int newTot = add(tot, mul(j, p[i - 1]));
            // cout << i << " " << j << " " << newMask << " " << newTot << endl;
            for (int mask = 0; mask < (1 << K); mask++) {
                int foo = newMask | mask;
                int cnt = __builtin_popcount(foo);
                if (cnt > k) {
                    continue;
                }
                // if (g[i - 1][mask])
                //     // cout << mask << " " << g[i - 1][mask] << " " << f[i - 1][mask] << endl;
                int bar = mul(newTot, g[i - 1][mask]);
                bar = add(bar, f[i - 1][mask]);
                res = add(res, bar);
            }
        }
        curMask |= (1 << u);
        tot = add(tot, mul(u, p[i - 1]));
    }
    return res;
}

int main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    init();
    // k = 1;
    // cout << get(30) << endl;
    long long l, r;
    cin >> l >> r >> k;
    // cout << get(l) << endl;
    cout << sub(get(r + 1), (get(l))) << endl;
    return 0;
}