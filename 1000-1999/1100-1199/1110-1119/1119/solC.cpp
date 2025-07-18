/*
 Author: LargeDumpling
 Email: LargeDumpling@qq.com
 Edit History:
	2019-04-06	File created.
*/

#include<iostream>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<cmath>
#include<algorithm>
using namespace std;
int n,m;
int c[505],r[505];
void read1n(int &x)
{
	char ch;
	for(ch=getchar();ch<'0'||'9'<ch;ch=getchar());
	for(x=0;'0'<=ch&&ch<='9';ch=getchar())
		x=(x<<1)+(x<<3)+ch-'0';
	return;
}
int main()
{
	int x;
	bool flag=true;
	read1n(n); read1n(m);
	for(int i=1;i<=n;i++)
		for(int j=1;j<=m;j++)
		{
			read1n(x);
			r[i]^=x;
			c[j]^=x;
		}
	for(int i=1;i<=n;i++)
		for(int j=1;j<=m;j++)
		{
			read1n(x);
			r[i]^=x;
			c[j]^=x;
		}
	for(int i=1;flag&&i<=n;i++)
		if(r[i])
			flag=false;
	for(int i=1;flag&&i<=m;i++)
		if(c[i])
			flag=false;
	puts(flag?"Yes":"No");
	fclose(stdin);
	fclose(stdout);
	return 0;
}