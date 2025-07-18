#include<bits/stdc++.h>
typedef long long ll;

using namespace std;

int main()
{
	ll t;
	cin >> t;
	while(t--) {
		ll a, b, c, d, k;
		cin >> a >> b >> c >> d >> k;
		ll pen = ceil((double) a/c);
		ll pencil = ceil((double) b/d);
		
		if(pen + pencil <= k) {
			cout << pen <<" " << pencil << endl;
		}
		else cout << -1 <<endl;
	}
	
	return 0;
}