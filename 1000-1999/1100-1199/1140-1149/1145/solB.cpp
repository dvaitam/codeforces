#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
inline int read()
{
	   int x=0,f=1;char c=getchar();
	   while(c<'0'||c>'9'){if(c=='-')f=-1;c=getchar();}
	   while(c>='0'&&c<='9'){x=(x<<3)+(x<<1)+c-'0';c=getchar();}
	   return x*f;
}
int a;
int main()
{
	a=read();
	puts(a==2||a==3||a==5||a==46||a==4||a==12||a==31||a==30||a==35||a==43||a==52||a==64||a==86?"YES":"NO");
}