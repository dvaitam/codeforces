#include <bits/stdc++.h>
using namespace std;

const int MAX_N = 300;
int a[MAX_N + 10], len[MAX_N + 10];
map<string, int> myMap;

int main() {
	int n;
	cin >> n;
	for (int i = 1; i <= n; i++) {
		string s;
		cin >> s;
		if (myMap[s] == 0) {
			a[i] = myMap[s] = i;
			len[i] = s.size();
		}
		else
			a[i] = myMap[s];
	}
	int sum = n - 1;
	for (int i = 1; i <= n; i++)
		sum += len[a[i]];
	int ans = 0;
	for (int i = 1; i <= n; i++)
		for (int j = i; j <= n; j++) {
			int cnt = 1, x = -1;
			for (int k = i; k <= j; k++)
				x += len[a[k]];
			int poi = j + 1;
			while (poi + j - i <= n) {
				bool ok = true;
				for (int k = 0; k < j - i + 1; k++)
					if (a[i + k] != a[poi + k]) {
						ok = false;
						break;
					}
				if (ok) {
					cnt++;
					poi += j - i + 1;
				}
				else
					poi++;
			}
			if (cnt > 1)
				ans = max(ans, x * cnt);
		}
	cout << sum - ans;
}