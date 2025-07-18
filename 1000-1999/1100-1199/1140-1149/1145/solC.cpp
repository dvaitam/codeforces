/********************************************************************
Author        : JiangYu
Created Time  : 2019年 04月 01日 23时 08分 51秒
********************************************************************/
/********************************************************************
Author        : JiangYu
Created Time  : 2019年 04月 01日 23时 08分 51秒
********************************************************************/
#include <bits/stdc++.h>
#define LL long long
#define DB double
#define fi first
#define se second
#define pb push_back
#define dep(i, a, b) for(int i = b; i > a; --i)
#define rep(i, a, b) for(int i = a; i < b; ++i)
//scanf printf
const int N = 2e5 + 1024;
const int mod = 1e9 + 7;
LL po(LL x, LL y) {LL ans=1;while(y){if(y&1)ans=ans*x%mod;y>>=1;x=x*x%mod;}return ans;}
using namespace std;

//-----


int a[20];
int ans = 1;

//#define LOCAL
int main(void) {
#ifdef LOCAL
    freopen("in.txt", "r", stdin);
    freopen("out.txt", "w", stdout);
#endif
	int n;
	cin >> n;
	if(n == 1) cout <<14;
	else if(n == 2) cout << 12;
	else if(n == 3) cout << 13;
	else if(n == 4) cout << 8;
	else if(n == 5) cout << 9;
	else if(n == 6) cout << 10;
	else if(n == 7) cout << 11;
	else if(n == 8) cout << 0;
	else if(n == 9) cout << 1;
	else if(n == 10) cout << 2;
	else if(n == 11) cout << 3;
	else if(n == 12) cout << 4;
	else if(n == 13) cout << 5;
	else if(n == 14) cout << 6;
	else if(n == 15) cout << 7;
	else cout << 15;

    return 0;
}