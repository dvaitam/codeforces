#include<bits/stdc++.h>
using namespace std;

int main()//prime(int g)
{
	//ios_base::sync_with_stdio(false);
//cin.tie(NULL);
	int n,i,d,a[101],t=0;
	cin>>n>>d;
	for(i=0;i<n;i++)
		cin>>a[i];

	for(i=0;i<n;i++)
	{
		if(i!=0)
		{
			if(abs(a[i]-d-a[i-1])>=d&&a[i]-a[i-1]>2*d)
				t++;
		}
		if(i!=n-1)
		{
			if(abs(a[i]+d-a[i+1])>=d)
				t++;
		}
	}
	cout<<t+2;

}