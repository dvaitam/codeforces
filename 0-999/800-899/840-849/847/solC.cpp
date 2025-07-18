#define _CRT_SECURE_NO_WARNINGS
#include <iostream>
#include <string>
#include <cstring>
#include <cmath>
#include <algorithm>
#include <vector>
#include <map>
#include <unordered_map>
#include <set>
#include <queue>
#include <stack>
#include <iomanip>
using namespace std;

inline void boostIO() {
	ios_base::sync_with_stdio(0);
	cin.tie(0);
	cout.tie(0);
}

long long n, k, i, d, sum;
string ans;

int main() {
	boostIO();
	cin >> n >> k;

	sum = (n - 1) * n / 2;

	if (sum < k) {
		return cout << "Impossible" << endl, 0;
	}

	if (k == 0) {
		for (int i = 0; i < n; ++i) {
			cout << "()";
		}
		return 0;
	}

	i = n - 1;

	while (sum >= k) {
		sum -= i--;

		if (sum < k) {
			d = k - sum;

			n = n - ans.size() / 2;

			for (i = 0; i < n; ++i) {
				ans += '(';

				if (i == d - 1) {
					ans += "()";
					n--;
				}
			}
			for (i = 0; i < n; ++i) {
				ans += ')';
			}

			break;
		}
		else {
			ans += "()";
		}
	}

	cout << ans << endl;
	return 0;
}