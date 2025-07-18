#include<iostream>
#include<vector>
#include<algorithm>
#include<string>
#include<map>
#include<iterator>
#include<stack>
#include<string>
#include<climits>
#include<bitset>
#include<queue>
#include<cmath>
#include<cstdlib>
#include<cstdio>
#include<sstream>

#define MOD 1000000007
#define ll long long int
#define vi vector<int>
#define vll vector<long long int>
#define vvi vector<vector<int> >
#define vvl vector<vector<long long int> >
#define vp vector<pair<int, int> >
#define sc(n) scanf("%d", &n)
#define ssync ios_base::sync_with_stdio(false), cin.tie(0), cout.tie(0)

using namespace std;

string toBin(long long int a)
{
	return bitset<64>(a).to_string();
}

void printArr(int* a, int n)
{
	for(int i = 0; i < n; i++)
		printf("%d ", a[i]);
	printf("\n");
}

string intToString (ll a)
{
	ostringstream temp;
	temp<<a;
	return temp.str();
}

ll pow(ll base, ll exponent)
{
	ll result = 1;
	while(exponent > 0)
	{
		if (exponent % 2 == 1)
			result *= base;
		exponent = exponent >> 1;
		base *= base;
	}
	return result;
}

ll powerWithMod(ll base, ll exponent, ll modulus)
{
	ll result = 1;
	base %= modulus;
	while(exponent > 0)
	{
		if (exponent % 2 == 1)
			result = (result * base) % modulus;
		exponent = exponent >> 1;
		base = (base * base) % modulus;
	}
	return result;
}

int main()
{
	ssync;
	bool *isPrime = new bool[2000001];
	fill(isPrime, isPrime+2000001, true);
	for(int i=4; i<2000001; i+=2)
		isPrime[i]=false;
	for(int d=3; d <= ceil(sqrt(2000001)); d+=2)
	{
		if(isPrime[d])
		{
			for(int i=2; d*i < 2000001; i++)
				isPrime[d*i]=false;
		}
	}
	int n, count=0;
	cin>>n;
	vi a(n), ans;
	for(int i=0; i<n; i++)
	{
		cin>>a[i];
		if(a[i] == 1)
			ans.push_back(1);
	}
	for(int i=0; i<n; i++)
	{
		if(isPrime[a[i] + 1] && a[i]>1)
		{
			ans.push_back(a[i]);
			break;
		}
	}
	if(ans.size() > 1)
	{
		cout<<ans.size()<<"\n";
		for(int i=0; i<ans.size(); i++)
			cout<<ans[i]<<" ";
		cout<<"\n";
		return 0;
	}
	ans.clear();
	for(int i=0; i<n; i++)
	{
		for(int j=i+1; j<n; j++)
		{
			if(isPrime[a[i]+a[j]])
			{
				cout << "2\n" << a[i] << " " << a[j] << "\n";
				return 0;
			}
		}
	}
	cout << "1\n" << a[0] << "\n";
	return 0;
}