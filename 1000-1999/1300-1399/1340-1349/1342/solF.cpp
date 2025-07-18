// LUOGU_RID: 102524366
/*+Rainybunny+*/

#include <bits/stdc++.h>

#define rep(i, l, r) for (int i = l, rep##i = r; i <= rep##i; ++i)
#define per(i, r, l) for (int i = r, per##i = l; i >= per##i; --i)

typedef std::pair<int, int> PII;
#define fi first
#define se second

template <typename Tp>
inline void chkmin(Tp& u, const Tp& v) { v < u && (u = v, 0); }
template <typename Tp>
inline void chkmax(Tp& u, const Tp& v) { u < v && (u = v, 0); }
template <typename Tp>
inline Tp imin(const Tp& u, const Tp& v) { return u < v ? u : v; }
template <typename Tp>
inline Tp imax(const Tp& u, const Tp& v) { return u < v ? v : u; }

const int MAXN = 15, IINF = 0x3f3f3f3f;
int n, a[MAXN], sum[1 << MAXN], f[MAXN + 1][MAXN + 1][1 << MAXN];
PII pre[MAXN + 1][MAXN + 1][1 << MAXN];

int main() {
    int cas; scanf("%d", &cas);
    while (cas--) {
        scanf("%d", &n);
        rep (i, 0, n - 1) scanf("%d", &a[i]), sum[1 << i] = a[i];
        rep (S, 1, (1 << n) - 1) sum[S] = sum[S & -S] + sum[S ^ (S & -S)];

        rep (i, 0, n) rep (j, 0, n) memset(f[i][j], 0x3f, 1 << n << 2);
        f[0][0][0] = 0;
        rep (i, 0, n - 1) rep (j, i, n - 1) {
            rep (S, 0, (1 << n) - 1) if (int cur = f[i][j][S]; cur < IINF) {
                for (int T = repS ^ S, z = 1; z; T = (T - 1) & (repS ^ S)) {
                    z = T;
                    if (!(T >> j) || cur >= sum[T]) continue;
                    int p = __builtin_ctz(T >> j) + j + 1;
                    if (f[i + 1][p][S | T] <= sum[T]) continue;
                    f[i + 1][p][S | T] = sum[T];
                    pre[i + 1][p][S | T] = { j, S };
                }
            }
        }

        static int idx[MAXN]; std::iota(idx, idx + n, 1);
        std::function<void(int, int, int)> trace =
        [&](const int i, const int j, const int S) {
            if (!i) return ;
            int T = S ^ pre[i][j][S].se;
            rep (k, 0, n - 1) if (T >> k & 1 && k != j - 1) {
                printf("%d %d\n", idx[k], idx[j - 1]), idx[k] = -1;
                rep (t, k + 1, n - 1) --idx[t];
            }
            trace(i - 1, pre[i][j][S].fi, pre[i][j][S].se);
        };

        per (i, n, 1) rep (j, i, n) {
            if (f[i][j][(1 << n) - 1] < IINF) {
                printf("%d\n", n - i), trace(i, j, (1 << n) - 1);
                goto FIN;
            }
        }
        FIN: ;
    }
    return 0;
}