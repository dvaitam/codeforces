#include<iostream>
#include<cstdio>
#include<algorithm>
#include<cstring>
#include<cmath>
#include<vector>
#include<queue>
#include<map>
#include<set>
using namespace std;
#define A first
#define B second
#define mp make_pair
#define pii pair<int,int>
#define ins insert
#define pb push_back
int n,k,ls,lt,lp,tt;
char st[1000],ss[1000],sp[1000],s[1000];
bool cheak(char *st){
    int tot=0;
    for (int i=0 ;i<lt; i++){
        int f=1;
        for (int j=0; j<ls; j++)
        if (ss[j]!=st[i+j]){f=0; break;}
        if (f==1) tot++;
    }
    if (tot!=tt) return 0;
    return 1;
}
int main(){
    //freopen("input.txt","r",stdin);
    cin>>n>>k;
    scanf("%s",ss); scanf("%s",sp);
    st[n]='\0';
    for (int i=0; i<n; i++) st[i]='!';
    lt=n; ls=strlen(ss); lp=strlen(sp);
    for (int i=0; i<lp; i++){
         if (sp[i]=='1'){
             tt++;
            for (int j=0; j<ls; j++)
                if (st[j+i]!='!'){
                    if (st[j+i]!=ss[j]){
                        puts("No solution");
                        return 0;
                    }
                } else st[j+i]=ss[j];
        }
    }
    for (int i='a'; i<='a'+k; i++){
        for (int j=0; j<lt; j++)
        if (st[j]=='!') s[j]=i;else s[j]=st[j];
        if (cheak(s)) {printf("%s",s);return 0;};
    }
    puts("No solution");
}