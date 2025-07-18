#include <cstdio>
#include <vector>
#include <algorithm>
#include <cstring>
using namespace std;

struct node {
	int a, b, c, d;
	node(int _t, int _a, int _b, int _c, int _d): a(_a), b(_b), c(_c), d(_d) {
	}
};
vector<node> ans;
void add(int n) {
	if (n & 1) {
		ans.emplace_back(3, n, 0, n + 1, -1);
		ans.emplace_back(3, n, 0, n + 1, -1);
		for (int i = 1; i < n; i += 2) {
			ans.emplace_back(4, n, i, n + 1, i + 1);
			ans.emplace_back(4, n, i, n + 1, i + 1);
		}
	} else {
		ans.emplace_back(3, n, 0, n + 1, -1);
		ans.emplace_back(3, n, 1, n + 1, -1);
		ans.emplace_back(4, n, 0, n + 1, 1);
		for (int i = 2; i < n; i += 2) {
			ans.emplace_back(4, n, i, n + 1, i + 1);
			ans.emplace_back(4, n, i, n + 1, i + 1);
		}
	}
}
int main() {
	int n;
	scanf("%d", &n);
	int now;
	if (n & 1) {
		ans.emplace_back(3, 0, 1, 2, -1);
		ans.emplace_back(3, 0, 1, 2, -1);
		now = 3;
	} else {
		ans.emplace_back(3, 0, 1, 2, -1);
		ans.emplace_back(3, 1, 2, 3, -1);
		ans.emplace_back(3, 2, 3, 0, -1);
		ans.emplace_back(3, 3, 0, 1, -1);
		now = 4;
	}
	while (now < n) {
		add(now);
		now += 2;
	}
	printf("%d\n", ans.size());
	for (auto i: ans) {
		if (i.d == -1)
			printf("3 %d %d %d\n", i.a + 1, i.b + 1, i.c + 1);
		else
			printf("4 %d %d %d %d\n", i.a + 1, i.b + 1, i.c + 1, i.d + 1);
	}
	return 0;
}