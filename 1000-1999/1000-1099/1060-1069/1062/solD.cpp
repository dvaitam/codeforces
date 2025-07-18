#include<iostream>
#include<algorithm>
#include<cmath>
#include<map>
#include<vector>
#include<set>
#include<iomanip>
#include<fstream>
#include<string>
#include<random>
using namespace std;
#define x first
#define y second
#define all(a) a.begin(), a.end()
#define pb push_back
typedef long long ll;
mt19937 rng(time(NULL));
int n;
int get_mul(int x) {
	int l = 1, r = n / 2 + 2;
	while (r - l > 1) {
		int m = (l + r) / 2;
		if (x * m > n)
			r = m;
		else
			l = m;
	}
	return l;
}
int main() {
	cin >> n;
	ll ans = 0;
	int m = get_mul(2);
	for (int i = 2; i <= m; i++) {
		int d = get_mul(i), sum;
		sum = (2 + d) * ((d - 2 + 1) / 2) + (!(d % 2)) * ((2 + d) / 2);
		ans += (1ll * sum) * 4;
	}
	cout << ans;
	return 0;
}