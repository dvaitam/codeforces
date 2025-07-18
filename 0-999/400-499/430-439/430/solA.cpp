#include <bits/stdc++.h>
using namespace std;
pair<int,int> a[1000];
int ans[1000];
int main()
{
	int n,m;
	cin >> n >> m;
	for (int i=0;i<n;i++)
	{
		cin>>a[i].first;
		a[i].second=i;
	}
	sort(a,a+n);
	int q=0;
	for (int i=0;i<n;i++)
	{
		ans[a[i].second]=q;
		q^=1;
	}
	for (int i = 0; i<n; i++)
		cout << ans[i]<<" ";
	return 0;
}