#include<bits/stdc++.h>
#define ll long long int
using namespace std;
main()
{
  ll n;
  cin>>n;
  ll a[n];
  ll b[n];
  ll m;
  for(ll i=0;i<n;i++)
  {
  	cin>>a[i];
  	cin>>b[i];
  }
  cin>>m;
  ll c=0;
  for(ll i=0;i<n;i++)
  {
  	if(a[i]>=m||b[i]>=m)
  	{
  		c++;
  		
	  }
  }
  cout<<c<<endl;
}