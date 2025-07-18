#pragma comment(linker, "/STACK:36777216")
#include <fstream>
#include <iostream>
#include <vector>
#include <map>
#include <set>
#include <string>
#include <deque>
#include <algorithm>
#define _USE_MATH_DEFINES
#include <math.h>
using namespace std;
#define forn(i,n) for (int i=0;i<n;i++)
#define rforn(i,n) for (int i=n-1;i>=0;i--)
#define mp make_pair
#define __int64 long long
#define LL long long
int main()
{
	ios_base::sync_with_stdio(false);
	#ifndef ONLINE_JUDGE
	freopen("input.txt","r",stdin);
	freopen("output.txt","w",stdout);
	#endif
	int n,k;
	cin>>n>>k;
	vector <int> a(n+1,1);
	a[n]=0;
	int mx=1,l;
	forn(i,k){
		for (int i=1;i<=n;i++)
		{
			if (a[i]==n-i)
				cout<<n<<' ';
			else
			if (a[i]+mx>=n-i)
			{
				cout<<i+a[i]<<' ';
				a[i]=n-i;
			}
			else
			{
				cout<<"1 ";
				a[i]+=mx;
			}
		}
		mx=mx*2;
		cout<<"\n";
	}
}