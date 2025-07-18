#include<bits/stdc++.h>
using namespace std;
#define rep(i,j,k) for(int i = j;i <= k;++i)
#define repp(i,j,k) for(int i = j;i >= k;--i)
#define rept(i,x) for(int i = linkk[x];i;i = e[i].n)
#define P pair<int,int>
#define Pil pair<int,ll>
#define Pli pair<ll,int>
#define Pll pair<ll,ll>
#define pb push_back 
#define pc putchar
#define mp make_pair
#define file(k) memset(k,0,sizeof(k))
#define ll long long
#define ls root * 2
#define rs root * 2 + 1
const int INF = 1e9+7;
int ans;
int n;
int read()
{
    int sum = 0;char c = getchar();bool flag = true;
    while(c < '0' || c > '9') {if(c == '-') flag = false;c = getchar();}
    while(c >= '0' && c <= '9') sum = sum * 10 + c - 48,c = getchar();
    if(flag) return sum;
    else return -sum;
}
struct cmp1{
	bool operator ()(int &a,int &b)
	{
		return a<b;
	}
};
struct cmp2{
	bool operator ()(int &a,int &b)
	{
		return a>b;
	}
};
priority_queue<int,vector<int>,cmp2>k2;//k1大根k2小根
priority_queue<int,vector<int>,cmp1>k1;
int tmp[401000],o;
void work(int x)
{
    rep(i,1,o)
    {
    	if(tmp[i] == x) continue;
    	if(tmp[i] < x) k1.push(tmp[i]);
    	else k2.push(tmp[i]);
    }
	return;
}
int main()
{
    n = read();
    ans = 1;
    int l = 0 , r = INF;
    rep(i,1,n)
    {
    	char c = getchar();
    	while(c != 'D' && c != 'C') c = getchar();
        if(c == 'D')
        {
        	int x = read();
        	tmp[++o] = x;
        }
        else
        {
        	int x = read();
        	if(x < l || x > r)
        	{
        		printf("0");
        		return 0;
        	}
        	bool flag = true;
        	if(x == l) k1.pop(),flag = false;
        	if(x == r) k2.pop(),flag = false;
        	work(x);
        	if(flag) ans = ans * 2 % INF;
        	if(!k1.empty())l = k1.top();
        	else l = 0;
        	if(!k2.empty())r = k2.top();
        	else r = INF;
        	o = 0;
        }
    }
    if(o == 0) printf("%d\n",ans);
    else
    {
    	int sum = 0;
    	rep(i,1,o) 
    	    if(tmp[i] < l) k1.push(tmp[i]);
            else if(tmp[i] > r) k2.push(tmp[i]);
            else sum++;
        ans = 1ll * ans * (sum+1) % INF;
        printf("%d\n",ans);
    }
    return 0;
}