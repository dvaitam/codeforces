#include <iostream>
#include <fstream>
#include <sstream>
#include <cstdlib>
#include <cstdio>
#include <cmath>
#include <string>
#include <cstring>
#include <algorithm>
#include <queue>
#include <stack>
#include <vector>
#include <set>
#include <map>
#include <list>
#include <iomanip>
#include <cctype>
#include <cassert>
#include <bitset>
#include <ctime>

using namespace std;

#define pau system("pause")
#define ll long long
#define pii pair<int, int>
#define pb push_back
#define pli pair<ll, int>
#define pil pair<int, ll>
#define clr(a, x) memset(a, x, sizeof(a))

const double pi = acos(-1.0);
const int INF = 0x3f3f3f3f;
const int MOD = 1e9 + 7;
const double EPS = 1e-9;

/*
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;
#define TREE tree<pli, null_type, greater<pli>, rb_tree_tag, tree_order_statistics_node_update>
TREE T;
*/

int n, ps[10015], cnt[10015];
int cal(int x, int f) {
	if (x == (x & -x)) return log2(x + 0.5) + f;
	return log2(x + 0.5) + 2;
}
int main() {
	scanf("%d", &n);
	int id = 0;
	for (int i = 2; i <= n; ++i) {
		if (n % i == 0) {
			ps[++id] = i;
			do {
				++cnt[id];
				n /= i;
			} while (n % i == 0);
		}
	}
	int ma = 1, res1 = 1, f = 0;
	for (int i = 1; i <= id; ++i) {
		res1 *= ps[i];
		ma = max(ma, cnt[i]);
		if (i < id && cnt[i] != cnt[i + 1]) f = 1;
	}
	printf("%d %d\n", res1, cal(ma, f));
	return 0;
}