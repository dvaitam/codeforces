#include<cstdio>
#include<cctype>
#include<algorithm>
using namespace std;
const int maxn=1000000,maxa=10000;
typedef long double DB;

int te,n,a[maxn+5],num[maxa+5],A,B;DB MIN;

#define Eoln(x) ((x)==10||(x)==13||(x)==EOF)
inline char readc(){
	static char buf[100000],*l=buf,*r=buf;
	if (l==r) r=(l=buf)+fread(buf,1,100000,stdin);
	if (l==r) return EOF;return *l++;
}
inline int readi(int &x){
	int tot=0;char ch=readc(),lst='+';
	while (!isdigit(ch)) {if (ch==EOF) return EOF;lst=ch;ch=readc();}
	while (isdigit(ch)) tot=(tot<<3)+(tot<<1)+(ch^48),ch=readc();
	return (lst=='-'?x=-tot:x=tot),Eoln(ch);
}
int main(){
	for (readi(te),MIN=1e100;te;te--,MIN=1e100){
		readi(n);for (int i=1;i<=n;i++) readi(a[i]),num[a[i]]++;sort(a+1,a+1+n);
		for (int i=1,j,lst=0;i<=n;i=j){
			for (j=i;j<=n&&a[i]==a[j];j++);
			if (num[a[i]]>3) MIN=1,A=a[i],B=a[i];
			if (lst&&num[a[i]]>1&&(DB)a[i]/lst<MIN) MIN=(DB)a[i]/lst,A=lst,B=a[i];
			if (num[a[i]]>1) lst=a[i];
		}
		for (int i=1;i<=n;i++) num[a[i]]=0;printf("%d %d %d %d\n",A,A,B,B);
	}return 0;
}