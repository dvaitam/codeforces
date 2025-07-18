#include <cmath>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>
#include <set>
#include <map>
#include <stack>
#include <queue>
#include <algorithm>
using namespace std;

#define SET(p) memset(p, -1, sizeof(p))
#define CLR(p) memset(p, 0, sizeof(p))
#define MEM(p, v) memset(p, v, sizeof(p))
#define CPY(d, s) memcpy(d, s, sizeof(s))
#define ll long long
#define ld long double
#define mod 1000000007
#define inf 1LL<<60
#define pii pair< int, int >
#define pll pair< ll, ll >
#define psi pair< string, int >
#define N 100010
int main(int argc, char const *argv[])
{
	ld n,m;
	cin>>n>>m;
	if(n*m==1){
		printf("1\n");
		return 0;
	}
	ld ans = 1/n + ((n-1)/n)*((m-1)/(n*m-1));
	printf("%.9Lf\n", ans);
	return 0;
}