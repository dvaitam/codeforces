#include<cstdio>
using namespace std;
#define ff for(i=0;i<n;++i)
#define ast a[s[i]-'a'][t[i]-'a']
#define ats a[t[i]-'a'][s[i]-'a']
const int N=2e5+1;
char s[N],t[N];
int n,cnt,i,j;
int a[26][26];
int main(){
    scanf("%d",&n);
    scanf("%s %s",s,t);
    ff  if(s[i]!=t[i]) {ast=i+1;cnt++;}
    ff  if(s[i]!=t[i]&&ats) {printf("%d\n%d %d\n",cnt-2,ast,ats);return 0;}
    ff  {if(s[i]!=t[i]){
        for(j=0;j<26;++j)
            if(a[t[i]-'a'][j]){printf("%d\n%d %d\n",cnt-1,ast,a[t[i]-'a'][j]);return 0;}
        }
    }
    printf("%d\n-1 -1\n",cnt);
    return 0;
}