#include <cstdio>
#include <cstring>
#include <algorithm>
using namespace std;

const int maxn = 1e5+10;

int n, ans = 0, cnt = 0;
char a[maxn];

void FindCnt() {
	for (int i = 1; i < n; i++)
		if (a[i] == 'S' && a[i-1] == 'G')
			cnt++;
}

void FindG(int &pos, int &l, int &r) {
	l = r = -1;
	int flag = 1;
	for (; pos < n; pos++) {
		if (flag) {
			if (a[pos] == 'G') {
				l = pos;
				flag = 0;
			}
		}
		else {
			if (a[pos] == 'S') {
				r = pos;
				break;
			}
		}
	}
}

int main() {
	scanf("%d%s", &n, a);
	a[n++] = 'S';
	FindCnt();
	int l1 = -1, r1 = -1, l2 = -1, r2 = -1, i = 0;
	FindG(i, l1, r1);
	ans = r1 - l1;
	if (cnt > 1) {
		ans++;
		for (; i < n; ) {
			FindG(i, l2, r2);
			if (l2 == -1)
				break;
			if (r1 + 1 == l2) {
				if (cnt == 2) {
					ans = r1 - l1 + r2 - l2;
					break;
				}
				ans = max(ans, r2 - l1);
			}
			else 
				ans = max(ans, r2 - l2 + 1);
			l1 = l2, r1 = r2;
		}
	}
	printf("%d", ans);
	return 0;
}