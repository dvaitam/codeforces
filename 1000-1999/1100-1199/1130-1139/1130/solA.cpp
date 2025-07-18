#include<bits/stdc++.h>
using namespace std;
 
#define MP make_pair
#define pb push_back
#define rep(i,n) for(int i=0;i<n;i++)
#define REP(i,a,b) for(int i=a;i<=b;i++)
#define PER(i,a,b) for(int i=b;i>=a;i--)
#define X first
#define Y second
#define gcd __gcd
 
//i/o
#define inp(n) scanf("%d",&n)
#define inpl(n) scanf("%lld",&n)
#define inp2(n,m) inp(n), inp(m)
#define inp2l(n,m) inpl(n), inpl(m)
 
 
//cost
#define lli long long int
#define MOD 1000000007
#define MOD_INV 1000000006
#define MAX 100009
#define INF 999999999
#define mp make_pair

//debug
#define debug() printf("here\n")
#define chk(a) cerr << endl << #a << " : " << a << endl
#define chk2(a,b) cerr << endl << #a << " : " << a << "\t" << #b << " : " << b << endl


//iterators
#define vitr std::vector<int>::iterator

int main() 
{
	int n, A[105];
	inp(n);
	for(int i=0; i<n; ++i)
	{
		inp(A[i]);
	}
	for(int i=-1000; i<1000; ++i)
	{
		if(i == 0)
			continue;
		int cnt = 0;
		for(int j=0; j<n; ++j)
		{
			double x = (double)A[j]/i;
			if(x > 0)
				cnt++;
		}
		if(cnt >= ceil((double)n/2))
		{
			printf("%d\n",i);
			return 0;
		}
	}
	printf("0\n");
	return 0;
}