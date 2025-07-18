#include <iostream>
#include<string>
#include <bits/stdc++.h> 
#include<cmath>
#include<vector>
#include<algorithm>
typedef unsigned long long int ull;
typedef long double ldb;
typedef long long int ll;

#define ForA1(n)  for (ll i=0; i<n; i++)
#define ForA2(n)  for (ll j=0; j<n; j++)
#define ForA3(n)  for (ll k=0; k<n; k++)
#define ForN1(n)  for (ll i=1; i<=n; i++)
#define ForN2(n)  for (ll j=1; j<=n; j++)
#define ForN3(n)  for (ll k=1; k<=n; k++)

#define mod 1000000007
#define pb push_back
#define vi vector<int>
#define F first
#define S second
#define mem(x) memset(x,0,sizeof(x))
#define PI 3.1415926535897932384626433832795l
#define deci(n) cout<<fixed<<setprecision(n);

#define number 100000

using namespace std;
  

void doit(ull n , ll k)
{

	ll n2 = n/2;

if(k<1)
{
	return;
}
	if(k==1)
	{
		cout<<n<<" ";
		return;
	}

	if(k==2)
	{
		cout<<n2<<" "<<n2<<" ";
		return; 


	}

	if(k <= n2)
	{
		cout<<n2<<" ";
		doit(n2,k-1); 
		return;
	}

	else
	{
		for(ll i=0;i<n2;i++)
		{
			cout<<"1 ";
		}

		k-=n2;
		doit(n2,k);
		return;
	}





}




int main()
{
    
    ios_base::sync_with_stdio(false);
	cin.tie(NULL);

	ull n;
	cin>>n;

	ll k;
	cin>>k;

	ull n2=n;

	vector<ll>arr;
	ll ctr=0;

	while(n2>0)
	{
		
		int t= n2%2;

		if(t==1)
		{
			arr.pb(pow(2,ctr));
		}

		ctr++;
		n2/=2;
	}

	// for(int i=0;i<arr.size();i++)
	// cout<<arr[i]<<" ";

	sort(arr.rbegin(),arr.rend());

	ll len = arr.size();


	if( k > n || k < arr.size()  )
	{
		cout<<"NO";
		return 0;
	}




	cout<<"YES\n";
	for(ll i=0;i<arr.size();i++)
	{
		if( k <= (arr[i] + len - i - 1 ) )
		{
			// cout<<"!!!!";
			doit(arr[i],k - (len - i - 1  ));

			for( ll j=i+1 ; j<arr.size() ;j++ )
			{
				cout<<arr[j]<<" ";
	
			}
		return 0;
		
		}

		else{
			// cout<<arr[i];
			for( ll l=0;l<arr[i];l++)
			{
				cout<<"1 ";
			}

			k-=arr[i];

			// for( ll j=i+1 ; j<arr.size() ;j++ )
			// {
			// 	cout<<arr[j]<<" ";
	
			// }

			// return 0;
	
		}

	}


















return 0;
}