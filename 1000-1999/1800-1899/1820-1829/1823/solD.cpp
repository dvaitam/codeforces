#include <bits/stdc++.h>

using namespace std;

using ll = int64_t;
using vll = vector<ll>;
using pll = pair<ll, ll>;
using vpll = vector<pll>;

#define forcin(n, a) for (ll i = 0; i < n; ++i) cin >> a[i];
#define TS ll ts; cin >> ts; while (ts--)
#define all(x) ::begin(x), ::end(x)
constexpr ll INF = 1e18;

ll n, k;
vpll a;

bool check() {
	ll lp = 0;
	ll lc = 0;
	for (ll i = 0; i < k; ++i) {
		if (a[i].first - lp < a[i].second - lc) return false;
		if (a[i].first <= 3) {
			if (a[i].second != a[i].first) return false;
		}
		else {
			if (a[i].second < 3) return false;
		}

		lp = a[i].first;
		lc = a[i].second;
	}

	return true;
}

int main() {
	ios_base::sync_with_stdio(false);
	cin.tie(0);
	cout.tie(0);
	
	TS {
		cin >> n >> k;
		a.resize(k);
		for (ll i = 0; i < k; ++i) cin >> a[i].first;
		for (ll i = 0; i < k; ++i) cin >> a[i].second;

		if (!check()) {
			cout << "NO\n";
			continue;
		}
		for (auto& [x, _] : a) {
			--x;
		}

		cout << "YES\n";
		sort(all(a));

		string res(n, '-');
		res[0] = 'a';
		res[1] = 'b';
		res[2] = 'c';

		ll pal = 3;
		char c = 'd';
		for (auto [pos, cnt] : a) {
			// cerr << pos << ", " << cnt << ": " << res << '\n';
			if (pos < 3) {
				continue;
			}
			for (ll i = 0; i < cnt - pal; ++i) {
				res[pos - i] = c;
			}

			pal = cnt;
			++c;
		}

		c = 'a';
		for (ll i = 3; i < n; ++i) {
			if (res[i] == '-') {
				res[i] = c;
				++c;
				if (c == 'd') c = 'a';
			}
		}

		cout << res << '\n';
	}
	
	return 0;
}