//*Anup Ghosh*/

#include<bits/stdc++.h>

using namespace std;

#define ll long long int

#define pi acos(-1)

#define ull unsigned long long

#define nl printf("\n")

#define MAXN 1000005

#define gcd(a,b) __gcd(a,b)

#define pb push_back

#define     all(x)      x.begin() , x.end()

#define fio() ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);

#define sp(n)  fixed<<setprecision(n)

#define    mod    998244353

#define no cout<<"No\n";

#define yes cout<<"Yes\n";

#define INF 100000000

vector<ll> st;

// prime factorization

set<ll> primefactorize(ll n)

{

	ll i,j;

	ll c=0;

	set<ll> s;

	for(i=2;i*i<=n;i++)

	{

		if(n%i==0)

		{

			while(n%i==0)

			{

				s.insert(i);

				n=n/i;

			}

		}

	}

	if(n>1) s.insert(n);

	return s;

}

// seive of eratosthenes

void seive()

{

	int n = 90000000;

	bool prime[90000001];

	int i,j;

	for(i=2;i*i<=n;i++)

	{

		if(prime[i]==false)

		{

			for(j=i*i;j<=n;j+=i)

			{

				prime[j]=true;

			}

		}

	}

	for(i=2;i<=n;i++)

	{

		if(prime[i]==false)

		st.pb(i);

	}

}

//binary exponentiation

ll power(ll n,ll p)

{

	ll res=1;

	while(p)

	{

			if(p%2)

	{

		res=(res*n)%10;

		p--;

	}

	else

	{

		n=(n*n)%10;

		p/=2;

	}



	}



	return res;

}

void print(ll a[],ll n)

{

	for(ll i=0;i<n;i++)

	cout<<a[i]<<" ";

	cout<<endl;

}

bool find(ll n)

{

	return n && (!(n&(n-1)));

}

void solve()

{

	string s;

	cin>>s;

	s[s.size()-1]=s[0];

	cout<<s<<endl;

	

}

int main()

{

	fio();

//	poW();

    int t=1;

  	cin>>t;

	while(t--)

	{

		solve();

	}



	return 0;

}