#include <iostream>
#include <cstring>
using namespace std;
const int M=2048;
char s[M],b[M],e[M],fb[M],fe[M];
int r[2][M],*p=r[0],*q=r[1],h[M];
int main(void) {
  cin>>s>>b>>e;
  int ns=strlen(s),nb=strlen(b),ne=strlen(e),nn=nb>ne?nb:ne,w=0;
  for(int i=ns;i-->0;){
    if(!strncmp(s+i,b,nb))fb[i]=1;
    if(!strncmp(s+i,e,ne))fe[i]=1;
  }
  for(int i=0,*t;i<ns;++i,t=p,p=q,q=t)
    for(int j=i;j-->0;){
      q[j+1]=(s[i]==s[j])?(p[j]+1):0;
      if(h[i]<q[j+1])h[i]=q[j+1];
    }
  for(int i=ns;i-->0;)if(fb[i])
    for(int j=i+nn;j<=ns;++j)if(fe[j-ne]&&h[j-1]<j-i)++w;
  cout<<w<<endl;
}