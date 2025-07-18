#include <bits/stdc++.h>
using namespace std;
#define ll long long
#define ld long double
#define uint unsigned int
#define ull unsigned long long
#define pii pair<int, int>
#define pll pair<ll, ll>
#define fi first
#define se second

const int N = 2e5 + 5;
const int inf = 1e9;

int a[N], dp[N], prv[N], ans[N];
priority_queue<pii> pq;

int main() { 
    ios_base::sync_with_stdio(false); cin.tie(NULL);
    int n;
    cin >> n;
    for (int i = 1; i <= n; i++) {
    	cin >> a[i];
    	dp[i] = -1;
    }
    a[n + 1] = inf;
    dp[1] = inf;
    pq.push({1, 0});
    int l = 1;
    for (int i = 2; i <= n + 1; i++) {
    	pii mx = {-1, 0};
    	if (a[i] > a[i - 1]) mx = {dp[i - 1], i - 1};
    	while (!pq.empty() and pq.top().se < l - 1) pq.pop();
    	if (!pq.empty() and -pq.top().fi < a[i]) {
    		mx = max(mx, {a[i - 1], pq.top().se});
    	}
    	dp[i] = mx.fi;
    	prv[i] = mx.se;
    	if (dp[i - 1] > a[i]) pq.push({-a[i - 1], i - 1});
    	if (a[i] >= a[i - 1]) l = i;
    }
    int cur;
    if (dp[n] != -1) {
    	cout << "YES\n";
    	cur = n;
    }
    else if (dp[n + 1] != -1) {
    	cout << "YES\n";
    	cur = n + 1;
    }
    else {
    	cout << "NO\n";
    	return 0;
    }
    while (cur > 0) {
    	ans[cur] = 1;
    	cur = prv[cur];
    }
    for (int i = 1; i <= n; i++) cout << 1 - ans[i] << " ";
    cout << "\n";
}