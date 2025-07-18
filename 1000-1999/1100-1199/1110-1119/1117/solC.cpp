#include <cstdio>
#include <memory.h>
#include <cstring>
#include <vector>
#include <deque>
#include <cstdlib>
#include <queue>
#include <algorithm>
#include <cmath>
#include <cassert>
#include <functional>
#include <iostream>
#include <set>
#include <list>
#include <map>
#include <time.h>
#include <unordered_map>
#include <unordered_set>
#include <bitset>
#define sz(x) (int)(x).size()
#define all(x) (x).begin(), (x).end()
using namespace std;

typedef unsigned long long llu;
typedef long long ll;
typedef pair<int, int> pii;
typedef pair<int, pii> piii;
typedef pair<ll, ll> pll;
typedef pair<ll, int> pli;
typedef pair<int, ll> pil;
typedef pair<string, int> psi;
const ll MOD = 1e9 + 7;
const long double PI = 3.141592653589793238462643383279502884197;

priority_queue<int, vector<int>, greater<int> > pq;
vector<int> v;

char s[200003];

int main() {
	ll x1, y1, x2, y2;
	int n;
	scanf("%lld %lld %lld %lld %d", &x1, &y1, &x2, &y2, &n);
	scanf("%s", s + 1);

	x2 -= x1;
	y2 -= y1;

	int xx = 0, yy = 0;
	for (int i = 1; i <= n; i++) {
		if (s[i] == 'U') yy--;
		else if (s[i] == 'D') yy++;
		else if (s[i] == 'L') xx++;
		else xx--;
	}

	ll l = 1, r = 2000000000, mid;
	while (l <= r) {
		mid = (l + r) / 2;
		ll tmp = abs(x2 + xx * mid) + abs(y2 + yy * mid) - mid * n;
		if (tmp <= 0) r = mid - 1;
		else l = mid + 1;
	}
	if (r >= 1e9 + 1e8) return !printf("-1");
	ll ans = r * n;
	
	x2 += xx * r;
	y2 += yy * r;

	for (int i = 1; i <= n; i++) {
		if (s[i] == 'U') y2--;
		else if (s[i] == 'D') y2++;
		else if (s[i] == 'L') x2++;
		else x2--;

		if (abs(y2) + abs(x2) <= i + r * n) return !printf("%lld", ans + i);
	}

	// r만큼 가능.
	// if (abs(xx) + abs(yy) == n && ((x2 > 0 && xx < 0) || (x2 < 0 && xx > 0) || (y2 < 0 && yy > 0) || (y2 > 0 && yy < 0))) return !printf("-1");
}