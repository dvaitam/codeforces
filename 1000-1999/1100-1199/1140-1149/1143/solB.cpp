#include<bits/stdc++.h>
#include<math.h>
using namespace std;
#define ll long long
int main()
{
	ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);
	ll n;
	cin>>n;
	int x=0,n1=n;
	while(n1>0)
	{
		x++;
		n1=n1/10;
	}
	ll ans=1;
	ll mans=1;
	n1=n;
	for(int i=0;i<x;i++)
	{
	    ans*=n1%10;
	    n1/=10;
	}
	if(ans>mans)mans=ans;
	ans=1;
	n1=n-1;
	for(int i=0;i<x;i++)
	{
	    ans*=n1%10;
	    n1/=10;
	}
	if(ans>mans)mans=ans;
	ans=1;
	int x1=x-1;
	while(x1>=0)
	{
	    n1=n-(n%(ll)pow(10,x1))-1;
	    for(int i=0;i<x;i++)
    	{
    	    if(n1%10==0)break;
    	    ans*=n1%10;
    	    n1/=10;
    	}
    	if(ans>mans)mans=ans;
	    ans=1;
	    x1--;
	}
	cout<<mans;
	return 0;
}