#include<bits/stdc++.h>
using namespace std;
int main()
{
	int i,j,n,x,y,t;
	cin>>t;
	while(t--)
	{
		cin>>n>>x>>y;
		for(i=1;i<y;i++)cout<<((y-i)%2?-1:1)<<' ';
		for(;i<=x;i++)cout<<"1 ";
		for(;i<=n;i++)cout<<((i-x)%2?-1:1)<<' ';
		cout<<'\n';
	}
	
	
	return 0;
}