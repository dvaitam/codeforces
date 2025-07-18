#include <cmath>
#include <cstdio>
#include <cstring>
#include <iostream>
#include <algorithm>
#define SG string
#define DB double
#define LL long long
using namespace std;
const LL Max=2e5+5;
LL N,M,K,Ans[Max];
inline LL Read(){
	LL X=0;char CH=getchar();bool F=0;
	while(CH>'9'||CH<'0'){if(CH=='-')F=1;CH=getchar();}
	while(CH>='0'&&CH<='9'){X=(X<<1)+(X<<3)+CH-'0';CH=getchar();}
	return F?-X:X;
}
inline void Write(LL X){
	if(X<0)X=-X,putchar('-');
	if(X>9)Write(X/10);
	putchar(X%10+48);
}
int main(){
	LL I,J;
	N=Read(),M=Read(),K=Read();
	if(K<M||(N-1)*M<K){
		puts("NO");return 0;
	}puts("YES");
	LL P=K/M;
	for(I=1;I<=M;I++){
		if(I&1){
			Ans[I]=1+P;
		} else {
			Ans[I]=1;
		}
	}
	if((K%M)&1){
		for(I=1;I<=K%M-1;I++){
			if(I&1){
				Ans[I]++;
			}	
		}
		if(M&1){
			Ans[M]++;
		} else {
			Ans[M-1]=1+P+1;Ans[M]=2; 
		}
	} else {
		for(I=1;I<=K%M;I++){
			if(I&1){
				Ans[I]++;
			}
		}
	}
	for(I=1;I<=M;I++){
		Write(Ans[I]),putchar(' ');
	}
	return 0;
}