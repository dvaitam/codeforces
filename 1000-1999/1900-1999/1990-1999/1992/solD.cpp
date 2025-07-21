#include <bits/stdc++.h>
using namespace std;

int n,m,k;
string s;

main()
{
	int t; 
	cin>>t; 
	while (t--)
	{
		cin>>n>>m>>k>>s;

		int j=m;

		for (int i=0; i<n && j; ++i)
		{
			if (s[i]=='L') 
				j=m; 
			else if (s[i]=='W' && j==1 && k)
				k--;
			else
				j--;
		}

		cout<<(j?"YES\n":"NO\n");
	}
}