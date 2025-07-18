#include <bits/stdc++.h>
using namespace std;
using ll = long long;

int f(ll x) {
    int res = 0;
    while (x) {
        int d = x % 10;
        x /= 10;
        res += (d == 4 || d == 7) ? 1 : 0;
    }
    return res;
}

ll solve(int add, ll x, ll y, ll a) {
	if (f(a) == x && f(a + add) == y) {
		return a;
	}
	vector<ll> pw(11);
	pw[0] = 1;
	for (int i = 1; i <= 10; i++) {
		pw[i] = pw[i - 1] * 10;
	}
	const ll lnf = 1e18;
	vector g(11, vector(x + 1, vector(y + 1, vector(2, lnf))));
	auto upd = [&](ll &v, ll val) {
		v = min(v, val);
	};
	for (int a = 0; a < 10; a++) {
		int b = (a + add) % 10, nd = (a + add) / 10;
		if (f(a) <= x && f(b) <= y) {
			upd(g[0][f(a)][f(b)][nd], a);
		}
	}
	for (int i = 0; i < 10; i++) {
		for (int ca = 0; ca <= x; ca++) {
			for (int cb = 0; cb <= y; cb++) {
				for (int d = 0; d < 2; d++) {
					if (g[i][ca][cb][d] < lnf) {
						for (int na = 0; na < 10; na++) {
							int nb = (na + d) % 10, nd = (na + d) / 10;
							int nca = ca + f(na), ncb = cb + f(nb);
							if (nca <= x && ncb <= y) {
								upd(g[i + 1][nca][ncb][nd], g[i][ca][cb][d] + na * pw[i + 1]);
							}
						}
					}
				}
			}
		}
	}
	ll res = lnf;
	for (int i = 0; i <= 10; i++) {
		int pre = a / pw[i];
		for (int a = 1; a <= 9 - pre % 10; a++) {
			if (i == 0) {
				if (f(pre + a) == x && f(pre + a + add) == y) {
					upd(res, pre + a);
				}
			} else {
				for (int d = 0; d < 2; d++) {
					int dx = x - f(pre + a), dy = y - f(pre + a + d);
					if (dx >= 0 && dy >= 0 && g[i - 1][dx][dy][d] < lnf) {
						upd(res, (pre + a) * pw[i] + g[i - 1][dx][dy][d]);
					}
				}
			}
		}
	}
	return res;
}

ll solve10(ll a, ll b) {
	ll ans = ll(1e10) + a;
	for (int d = 0; d <= 9; d++) {
		vector<int> fv(20, -1000);
		for (int i = 0; i <= b - a; i++) {
			fv[d + i] = f(a + i) - f((d + i) % 10);
		}
		set<int> pv, sv;
		for (int i = 0; i < 10; i++) {
			pv.insert(fv[i]);
			sv.insert(fv[i + 10]);
		}
		pv.erase(-1000), sv.erase(-1000);
		if ((int)pv.size() > 1 || (int)sv.size() > 1) {
			continue;
		}
		ll lim = a / 10;
		if (d <= a % 10) {
			++lim;
		}
		if (sv.empty()) {
			int x = *pv.begin();
			if (x >= 0) {
				ans = min(ans, solve(0, x, x, lim) * 10 + d);
			}
		} else {
			int x = *pv.begin(), y = *sv.begin();
			if (x >= 0 && y >= 0) {
				ans = min(ans, solve(1, x, y, lim) * 10 + d);
			}
		}
	}
	return ans;
}

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr), cout.tie(nullptr);
    ll a, l;
    cin >> a >> l;
    ll b = a + l - 1;
    ll base = 1, ax = 0;
    while (b >= a + 9) {
        ax += a % 10 * base;
        a /= 10, b /= 10, base *= 10;
    }
	ll ans = ax + solve10(a, b) * base;
    cout << ans << '\n';
    return 0;
}