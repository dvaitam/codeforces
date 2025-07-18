#include<bits/stdc++.h>
using namespace std;
const int maxn = 2e5 + 5;
const long long mod = 998244353;
struct Node{
	int id, val;
	bool operator < (const Node& a) {
		return val < a.val || (val == a.val && id < a.id);
	}
} a[maxn];
struct Node1{
	int l, r, val;
	bool operator < (const Node1& a) {
		return l < a.l;
	}
} b[maxn];
int ans[maxn];
int main()
{
	int n; scanf("%d", &n);
	for(int i=0;i<n;i++){
		scanf("%d",&a[i].val);
		a[i].id = i;
	}
	sort(a, a+n);
	int cnt = 0;
	for(int i = 1; i < n; i++)
		if(a[i].val == a[i-1].val)
			b[cnt++] = { a[i-1].id, a[i].id, a[i].val };
	sort(b, b+cnt);
	int lp = b[0].l, ip = 0;
	int col = 1;
	while(lp < n && ip < cnt) {
		while(lp > b[ip].r) {
			if(++ip >= cnt) break;
			if(b[ip].val != b[ip-1].val && lp <= b[ip].l) col++;
			lp = max(lp, b[ip].l);
		}
		if(ip >= cnt) break;
		ans[lp++] = col;
	}
	long long anss = 1;
	for(int i = 1; i < n; i++)
		if(ans[i]==0 || ans[i]!=ans[i-1])
			anss = anss*2%mod;
	printf("%I64d\n", anss);
	return 0;
}