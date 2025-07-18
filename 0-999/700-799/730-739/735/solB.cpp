#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <set>
#include <map>
#include <iostream>
#include <algorithm>
#include <cstring>
#include <queue>
#include <vector>
#include <stack>
#include <string>
using namespace std;
#define lson l,m,rt<<1
#define rson m+1,r,rt<<1|1
#define LL long long
#define rep1(i,a,b) for (int i = a;i <= b;i++)
#define rep2(i,a,b) for (int i = a;i >= b;i--)
#define mp make_pair
#define pb push_back
#define fi first
#define se second

typedef pair<int,int> pii;
typedef pair<LL,LL> pll;

void rel(LL &r)
{
	r = 0;
	char t = getchar();
	while (!isdigit(t) && t!='-') t = getchar();
	LL sign = 1;
	if (t == '-')sign = -1;
	while (!isdigit(t)) t = getchar();
	while (isdigit(t)) r = r * 10 + t - '0', t = getchar();
	r = r*sign;
}

void rei(int &r)
{
	r = 0;
	char t = getchar();
	while (!isdigit(t)&&t!='-') t = getchar();
	int sign = 1;
	if (t == '-')sign = -1;
	while (!isdigit(t)) t = getchar();
	while (isdigit(t)) r = r * 10 + t - '0', t = getchar();
	r = r*sign;
}

const int MAXN = 1e5+100;
const int dx[5] = {0,1,-1,0,0};
const int dy[5] = {0,0,0,-1,1};
const double pi = acos(-1.0);

int n,n1,n2;
int a[MAXN];

int main()
{
    //freopen("F:\\rush.txt","r",stdin);
    rei(n);rei(n1);rei(n2);
    rep1(i,1,n)
        rei(a[i]);
    sort(a+1,a+1+n);
    double ans = 0;
    int t = min(n1,n2);
    rep1(i,n-t+1,n)
        ans+=a[i];
    ans=ans/(t*1.0);
    double ans2= 0;
    int t1 = max(n1,n2);
    rep1(i,n-t-t1+1,n-t)
        ans2+=a[i];
    ans2 = ans2/(t1*1.0);
    ans+=ans2;
    printf("%.6lf\n",ans);
    return 0;
}