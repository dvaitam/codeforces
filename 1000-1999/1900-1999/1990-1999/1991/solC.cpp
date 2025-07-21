#include<bits/stdc++.h>
using namespace std;
int main()
{
	int T,n,x;
	cin>>T;
	while(T--)
	{
		cin>>n;
		int t=0,flag=0;
		for(int i=1;i<=n;i++)
		{
			cin>>x;
			if(i==1)t=x&1;
			else if(t!=(x&1))flag=1;
		}
		if(flag)cout<<"-1\n";
		else
		{
			t=!t;
			cout<<30+t<<'\n';
			for(int i=29;i>=0;i--)
				cout<<(1<<i)<<' ';
			if(t)cout<<t;
			cout<<'\n';
		}
	}
	return 0;
}