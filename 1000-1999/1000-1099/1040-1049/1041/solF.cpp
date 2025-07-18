#include <bits/stdc++.h>

using namespace std;
typedef long long ll;
typedef long double ld;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;

int a[100000], b[100000];
int ta[100000], tb[100000];
int solve(int as, int ae, int bs, int be) {
    int cnt[2] = {};
    for (int i = as; i <= ae; ++i) ++cnt[a[i] & 1];
    for (int i = bs; i <= be; ++i) ++cnt[(b[i] & 1) ^ 1];
    int ret = max(cnt[0], cnt[1]);
    if (as > ae || bs > be) return ret;
    if (as == ae && bs == be) return max(ret, 2);
    int j, k;
    int A, B;
    j = as;
    for (int i = as; i <= ae; ++i) if (a[i] & 1) ta[j++] = a[i] >> 1;
    A = j;
    for (int i = as; i <= ae; ++i) if (!(a[i] & 1)) ta[j++] = a[i] >> 1;
    for (int i = as; i <= ae; ++i) a[i] = ta[i];
    k = bs;
    for (int i = bs; i <= be; ++i) if (b[i] & 1) tb[k++] = b[i] >> 1;
    B = k;
    for (int i = bs; i <= be; ++i) if (!(b[i] & 1)) tb[k++] = b[i] >> 1;
    for (int i = bs; i <= be; ++i) b[i] = tb[i];
    return max({ ret, solve(as, A - 1, bs, B - 1), solve(A, ae, B, be) });
}

int main() {
	ios_base::sync_with_stdio(false); cin.tie(0);
	int n, m, y1, y2, x;
	cin >> n >> y1;
	vector<int> as;
	for (int i = 0; i < n; ++i) cin >> a[i];
	cin >> m >> y2;
	vector<int> bs;
	for (int i = 0; i < m; ++i) cin >> b[i];
    cout << solve(0, n - 1, 0, m - 1) << endl;
	return 0;
}