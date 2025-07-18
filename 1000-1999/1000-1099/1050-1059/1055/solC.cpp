//TMWCTW...
#include <bits/stdc++.h>

using namespace std;

#define int long long
const int maxn = 3e5 + 20, inf = 1e17;
int ans, l1, r1, t1, l2, r2, t2, g, k, f1, f2;

int32_t main(){
	ios::sync_with_stdio(false);cin.tie(0);cout.tie(0);
    cin >> l1 >> r1 >> t1;
    cin >> l2 >> r2 >> t2;
    ans = max(1ll * 0, min(r1, r2) - max(l1, l2) + 1);
    f1 = r1 - l1 + 1;
    f2 = r2 - l2 + 1;
	g = __gcd(t1, t2);
	k = abs(l1 - l2) % g;
	if(t1 == t2){
		cout << ans;
		return 0;
	}
	if(k == 0){
		cout << max(ans, min(f1, f2));
		return 0;
	}
	if(l2 > l1){
		if(f2 + k <= f1){
			cout << max(ans, f2);
			return 0;
		}
		if(f1 + g - k <= f2){
			cout << max(ans, f1);
			return 0;
		}
		cout << max(ans, max(f1 - k, f2 - g + k));
		return 0;
	}
	if(l2 < l1){
		if(f1 + k <= f2){
			cout << max(ans, f1);
			return 0;
		}
		if(f2 + g - k <= f1){
			cout << max(ans, f2);
			return 0;
		}
		cout << max(ans, max(f2 - k, f1 - g + k));
		return 0;
	}
}