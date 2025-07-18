#include <bits/stdc++.h>
using namespace std;

char str[500005];
int n;
map<int, int> mp;
int main() {
	scanf("%d", &n);
	for (int i = 1; i <= n; i++) {
		scanf("%s", str);
		int len = strlen(str);
		int cnt0 = 0, cnt1 = 0;
		for (int j = 0; j < len; j++) {
			if (str[j] == '(') ++cnt0;
			else {
				if (cnt0) --cnt0;
				else ++cnt1; 
			}
		}
		if (!cnt0 && !cnt1) ++mp[0];
		else if (!cnt0) ++mp[cnt1];
		else if (!cnt1) ++mp[-cnt0];
	}
	int res = 0; 
	for (pair<int, int> p : mp) if (p.first < 0) {
		int t = min(p.second, mp[-p.first]);
		res += t;
	}
	if (mp[0]) res += mp[0] / 2;
	printf("%d\n", res);
	return 0;
}