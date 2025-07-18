#include<stdio.h>
#include<string.h>
#include<algorithm>
using namespace std;
const int maxn=1e5+10;
int n,father[30],top,ansr[110],ansl[110];
char ch1[maxn],ch2[maxn];
int find(int x){
 if(father[x]!=x) father[x]=find(father[x]);
 return father[x];
}
void mis(int a,int b){
 a=find(a),b=find(b);
 if(a!=b) father[a]=b;
}
int main(){
 int i,j;
 top=0;
 scanf("%d",&n);
 scanf("%s%s",ch1,ch2);
 for(i=1;i<=26;i++) father[i]=i;
 for(i=0;i<n;i++){
  int a=ch1[i]-'a'+1;
  int b=ch2[i]-'a'+1;
  if(a!=b&&find(a)!=find(b)){
   ansr[top]=a;ansl[top++]=b;
   mis(a,b);
  }
 }
 printf("%d\n",top);
 for(i=0;i<top;i++){
  printf("%c %c\n",ansr[i]-1+'a',ansl[i]-1+'a');
 }
 return 0;
}