#include<bits/stdc++.h>
using namespace std;
#define FOR(a,b,c) for(int a=(b),a##_end__=(c);a<a##_end__;++a)
template<class T>inline bool chkmin(T&a,T const&b){return a>b?a=b,true:false;}
template<class T>inline bool chkmax(T&a,T const&b){return a<b?a=b,true:false;}
const int M=100005;
int A[M],B[M],C[M],D[M];
int n,m,ans;
long long d1,d2;
bool check(int X1,int Y1,int X2,int Y2,long long a1,long long b1,long long a2,long long b2){
	long long A=(a1-a2)*(a1-a2)+(b1-b2)*(b1-b2);
	long long B=((a1-a2)*(X1-X2)+(b1-b2)*(Y1-Y2))*2;
	long long C=(X1-X2)*(X1-X2)+(Y1-Y2)*(Y1-Y2);
	if(C<=d1*d1) return true;
	if(A+B+C<=d1*d1) return true;
	if(A and B<=0 and -B<=2*A){
		if(4*A*C-B*B<=4*d1*d1*A) return true;
	}
	return false;
}
long long DIS(int X1,int Y1,int X2,int Y2){
	return (X1-X2)*(X1-X2)+(Y1-Y2)*(Y1-Y2);
}

int main(){
	scanf("%d %lld %lld",&n,&d1,&d2);
	FOR(i,0,n) scanf("%d%d%d%d",A+i,B+i,C+i,D+i);
	bool flag=false;
	if(DIS(A[0],B[0],C[0],D[0])<=d1*d1){
		flag=true;
		++ans;
	}
	FOR(i,1,n){
		bool f=check(A[i-1],B[i-1],C[i-1],D[i-1],A[i]-A[i-1],B[i]-B[i-1],C[i]-C[i-1],D[i]-D[i-1]);
		if(!flag and f) ++ans;
		if(f) flag=true;
		if(DIS(A[i],B[i],C[i],D[i])>d2*d2) flag=false;
	}
	printf("%d\n",ans);
	return 0;
}