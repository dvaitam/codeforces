#include <bits/stdc++.h>
using namespace std;
#define ll long long
int vis[5];
char s[100005];
int main()
{

    int t;
    scanf("%d",&t);
    while(t--){
        scanf("%s",s+1);
        int st=strlen(s+1);
        int flag=0;
        for(int i=1;i<=st;i++){
            if(s[i]==s[i-1]||s[i]==s[i+1]){
                if(s[i]!='?'){flag=1;break;}
            }
        }
        if(!flag){
            for(int i=1;i<=st;i++){
                memset(vis,0,sizeof(vis));
                if(s[i]=='?'){
                    if(i-1>=1)vis[s[i-1]-'a']=1;
                    if(i+1<=st&&s[i+1]!='?')vis[s[i+1]-'a']=1;
                    for(int j=0;j<=2;j++){
                        if(!vis[j])s[i]=char(j+'a');
                    }
                }
            }
            cout<<s+1<<endl;
        }else cout<<-1<<endl;


    }
    return 0;
}