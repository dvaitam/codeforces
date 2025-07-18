#include <iostream>
#include <string>
#include <vector>

using namespace std;

int cnt[101];

int main() {
	int n, k, x;
	cin >> n >> k;
	for (int i = 0; i < n; ++i) {
		cin >> x;
		cnt[x]++;
	}
	int ans = 0, mx = 0;
	for (int i = 0; i <= 100; ++i) {
		if (mx < cnt[i]) mx = cnt[i];
	}
	if (mx % k) {
		mx -= mx % k;
		mx += k;
	}
	for (int i = 0; i <= 100; ++i) {
		if (cnt[i]) ans += mx - cnt[i];
	}
	cout << ans;
		return 0;
}