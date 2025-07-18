#pragma GCC optimize ("Ofast")
#include <bits/stdc++.h>
#define ll long long
#define pb push_back
#define mp make_pair
using namespace std;
int main()
{
	ios_base::sync_with_stdio(false);
	cin.tie(NULL);
	ll x,y,z,t1,t2,t3;
	cin >> x >> y >> z >> t1 >> t2 >> t3;
	if(abs(x-z)*t2 + abs(y-x)*t2 + 3*t3 <= abs(x-y)*t1) cout << "YES\n";
	else cout << "NO\n";
}