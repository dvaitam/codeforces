#include <iostream>
#include <algorithm>
#include <cassert>

using namespace std;

const int K = 50 + 7;
long long a, b, m, p[K], r[K];

void build() {
	p[1] = 1; p[2] = 1;
	for (int i = 3; i < K; ++i) p[i] = p[i - 1] *(long long)2;
}

bool try_make(long long k, long long b) {
	b -= p[k] * a;
	if (p[k] > b || p[k] * m < b) return false;
	for (int i = 1; i < k; ++i) ++r[i], b -= p[i];
	for (int i = 1; i < k; ++i) {
		r[i] += min(b / p[k - i], m - 1);
		b -= p[k - i] * (r[i] - 1);
		assert(b >= 0);
	}
	assert(b == 0);
	return true;
}

void solve() {
	fill(r, r + K, 0);
	cin >> a >> b >> m;
	if (a == b) {
		cout << 1 << " " << a << '\n';
		return;
	}
	for (int k = 2; k <= 50; ++k) {
		if (p[k] * a > b) break;
		if (try_make(k, b)) {
			cout << k << " ";
			long long s = 0;
			for (int i = 0; i < k; ++i) cout << a << " ", s += a, a = s + r[i + 1];
			cout << '\n';
			return;
		}
	}
	cout << -1 << endl;
}

int main() {
	ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);
	build();
	int t;
	cin >> t;
	while (t--) solve();
	//system("pause");
	return 0;
}