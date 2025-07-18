#include<iostream>
#include<vector>
#include<string>
#include<cstring>
#include<algorithm>
#include<map>
#include<set>
#include<cmath>
#include<cassert>
#include<queue>
#include<string.h>
typedef long long ll;
const double PI = 3.14159265358979323846;
const double EPS = 1e-12;
const ll INF = 1LL << 42;
const ll MOD = (ll)998244353;

using namespace std;

int main() {
	ios::sync_with_stdio(false); cin.tie(0);
    
    int n;
    ll ans;
	cin >> n;
	if(n >= 3){
	    ans = n-1 + n-2;
	    for (int k = 3; k < n; k++){
	        ans = (k * ans - 1) % MOD;
	    }
	    ans = (n * ans) % MOD;
	}
    if(n == 1){
        ans = 1;
    }
    if(n == 2){
        ans = 2;
    }

	cout << ans << endl;

	return 0;
}