//In the name of Allah
#include <bits/stdc++.h>
#define fs first
#define sc second
#define pb push_back

using namespace std;

typedef long long ll;
const int N=1e5+7, INF=1e7;
const double PI = 3.141592653589793238462643383279502884197, eps=1e-6;

int arr[N];

int getXor(int n)
{
	int ans=0;

	while(n--)
		ans ^= n;

	return ans;
}

int main()
{
	ios::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);

	int n, x, y, m;
	bool is=0;

	cin >> n >> x;

	if(n==1)
		return cout << "YES\n" << x, 0;

	if(n==2 && !x)
		return cout << "NO\n", 0;

	y=getXor(n-2);

	y |= (1<<17);
	x |= (1<<17);

	m = n-2;

	while(x==y && m<=1e6)
		is=1, y^=m, m++;

	cout << "YES\n";

	for(int i=is ; i<n-2 ; i++)
		cout << i << ' ';

	if(is)
		cout << m-1 << ' ';

	cout << y << ' ' << x;

	return 0;
}