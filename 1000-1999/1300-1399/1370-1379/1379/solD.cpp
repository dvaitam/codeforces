#include<bits/stdc++.h>
using namespace std;
struct IO{
	static const int S=1<<21;
	char buf[S],*p1,*p2;int st[105],Top;
	~IO(){clear();}
	inline void clear(){fwrite(buf,1,Top,stdout);Top=0;}
	inline void pc(const char c){Top==S&&(clear(),0);buf[Top++]=c;}
	inline char gc(){return p1==p2&&(p2=(p1=buf)+fread(buf,1,1<<21,stdin),p1==p2)?EOF:*p1++;}
	IO&operator >> (char&x){while(x=gc(),x==' '||x=='\n');return *this;}
	template<typename T> IO&operator >> (T&x){
		x=0;bool f=0;char ch=gc();
		while(ch<'0'||ch>'9'){if(ch=='-') f^=1;ch=gc();}
		while(ch>='0'&&ch<='9') x=(x<<3)+(x<<1)+ch-'0',ch=gc();
		f?x=-x:0;return *this;
	}
	IO&operator << (const char c){pc(c);return *this;}
	template<typename T> IO&operator << (T x){
		if(x<0) pc('-'),x=-x;
		do{st[++st[0]]=x%10,x/=10;}while(x);
		while(st[0]) pc('0'+st[st[0]--]);return *this;
	}
}fin,fout;
const int N=100005;
using namespace std;
int n,h,m,k,non,H[N];
struct Node{
	int id,v;
	friend bool operator <(const Node a,const Node b){return a.v<b.v;}
}a[N<<1];
int main(){
	fin>>n>>h>>m>>k,m/=2;
	for(int i=1;i<=n;i++){
		fin>>H[i]>>a[i].v;a[i].id=a[n+i].id=i;
		if(a[i].v>=m) a[i].v-=m;a[n+i].v=a[i].v+m;
	}
	int j=1,l=0,r=n+1;sort(a+1,a+(n<<1|1));
	for(int i=n+1;i<=(n<<1);i++){
		while(a[j].v<a[i].v-k+1) ++j;
		if(i-r<j-l) l=j,r=i;
	}
	fout<<r-l<<' '<<a[r].v-m<<'\n';
	for(int i=l;i<r;i++) fout<<a[i].id<<' ';
	return rand()&0;
}