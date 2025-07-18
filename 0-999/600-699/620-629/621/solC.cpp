#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cmath>
#include<string>
#include<vector>
#include<map>
#include<queue>
#include<list>
#include<cctype>
using namespace std;
#define LL long long
#define PB push_back
#define oo 0x3f3f3f3f
template<class T> T read(T &x) {
	int f = 1;
	char ch = getchar();
	while(!isdigit(ch)) { if(ch=='-') f = -f; ch = getchar(); }
	for(x = 0; isdigit(ch); ch = getchar()) x = x*10 + ch -48;
	return x = f *x;
}
//EOT
LL n,mo,FL,FR,Ll,LR,NL,NR;
long double Exp=0.0;

void work(LL a,LL b,LL c,LL d)
{
	LL A=b/mo-(a-1)/mo, B=d/mo-(c-1)/mo, C=A*B;
	long double Ans;
	Ans=(long double)(A*(d-c+1)+B*(b-a+1)-C)*2000;
	Exp+=(long double)Ans/(b-a+1)/(d-c+1);
}

int main()
{
	read(n), read(mo), read(FL), read(FR);
	Ll=FL; LR=FR;
	for (int i=2;i<=n;i++)
	{
		read(NL), read(NR);
		work(Ll,LR,NL,NR);
		Ll=NL; LR=NR;
	}
	work(FL,FR,NL,NR);
	printf("%.10lf\n",(double)Exp);
}