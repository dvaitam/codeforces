#include<bits/stdc++.h>
#define rgi register int
#define ll long long
class fastin{
	private: int _ch,_f;
	public: inline fastin& operator>>(char&c) {c=getchar(); return *this;}
		template<typename _Tp> inline fastin& operator>>(_Tp&_x){
		_x=0; while(!isdigit(_ch)) _f|=(_ch==45),_ch=getchar();
		while(isdigit(_ch)) _x=(_x<<1)+(_x<<3)+(_ch^48),_ch=getchar();
		_f&&(_x=-_x,_f=0); return *this;}
	fastin() {_ch=_f=0;}
}fin;
class fastout{
#define endl '\n'
	private: int _num[32],_head;
	public: inline fastout& operator<<(char c) {putchar(c); return *this;}
		template<typename _Tp> inline fastout& operator<<(_Tp _x){
		_Tp _k; if(_x==0) {putchar('0'); return *this;} if(_x<0) putchar('-'),_x=-_x;
		while(_x>0) _k=_x/10,++_head,_num[_head]=(_x-(_k<<1)-(_k<<3))^48,_x=_k;
		while(_head>0) putchar(_num[_head]),--_head; return *this;}
	fastout() {_head=0;}
}fout;
// ----------------------------
// #define int ll
// using namespace std;
const int maxn=304;
const int mod=998244353,inf=1000000007;

int a[maxn][maxn],b[maxn][maxn];
inline int calc(int i,int j) {return int(bool(b[i+1][j]))+bool(b[i-1][j])+bool(b[i][j+1])+bool(b[i][j-1]);}
signed main() // B
{
	int T;
	fin>>T;
	while(T--)
	{
		int n,m;
		fin>>n>>m;
		for(rgi i=0;i<=n+1;++i) for(rgi j=0;j<=m+1;++j) b[i][j]=0;
		for(rgi i=1;i<=n;++i) for(rgi j=1;j<=m;++j) fin>>a[i][j],b[i][j]=1;
		bool tag=1;
		for(rgi i=1;i<=n&&tag;++i) for(rgi j=1;j<=m&&tag;++j)
		{
			b[i][j]=calc(i,j);
			if(a[i][j]>b[i][j]) tag=0;
		}
		if(!tag) {puts("NO"); continue;}
		else
		{
			puts("YES");
			for(rgi i=1;i<=n;++i)
			{
				for(rgi j=1;j<=m;++j) fout<<b[i][j]<<' ';
				fout<<endl;
			}
		}
	}
	return 0;
}
// ----------------------------
// by imzzy