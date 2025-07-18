#include <bits/stdc++.h>
#define FR first
#define SE second
#define next next2

using namespace std;

typedef long long ll;
typedef pair<int,int> pr;

char s1[500005],s2[500005];

int next[500005];

void kmp(int n) {
  next[0]=next[1]=0;
  for(int i=1,j=0;i<n;i++) {
  	while (j&&s2[j]!=s2[i]) j=next[j];
  	if (s2[j]==s2[i]) j++;
  	next[i+1]=j;
  }
}

int main() {
  scanf("%s%s",s1,s2);
  int n=strlen(s1),m=strlen(s2);
  kmp(m);
  int a=0,b=0;
  for(int i=0;i<n;i++)
    if (s1[i]=='0') a++; else b++;
  int c=0,d=0;
  for(int i=0;i<m;i++)
    if (s2[i]=='0') c++; else d++;
  if (a<c||b<d) {
  	puts(s1);
  	return 0;
  }
  for(int i=0;i<m;i++) s1[i]=s2[i];
  a-=c;b-=d;
  for(int i=0;i<next[m];i++)
    if (s2[i]=='0') c--; else d--;
  int cur=m;
  while (a>=c&&b>=d) {
  	for(int i=next[m];i<m;i++) s1[cur++]=s2[i];
  	a-=c;b-=d;
  }
  while (a) {
  	s1[cur++]='0';
  	a--;
  }
  while (b) {
  	s1[cur++]='1';
  	b--;
  }
  puts(s1);
  return 0;
}