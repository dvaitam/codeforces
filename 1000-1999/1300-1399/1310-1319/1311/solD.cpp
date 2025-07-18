#include<bits/stdc++.h>
#define ll long long
#define ull unsigned long long
#define mem(a,b) memset(a,b,sizeof(a))
#define inf 0x3f3f3f3f
#define IO std::ios::sync_with_stdio(false);cin.tie(0)
#define endl '\n'
#define pii pair<int,int>
using namespace std;

const int maxn = 3e4+5;
const double eps = 1e-18;
const ull base  = 13331;
const int mod = 2333;

int ans,ans1,ans2,ans3;
int a,b,c;

void dfs1(int x,int y,int num){
    if(num>ans)return ;
    for(int i = 1;;i++){
        int tmp = y*i;
        if(tmp>c){
            if(tmp-c+num>ans)return ;
            int d = tmp-c;
            int ff = num+d;
            if(ff<ans){
                ans = ff;
                ans1 = x;
                ans2 = y;
                ans3 = tmp;
            }
            break;
        }
        else{
            int d = c-tmp;
            int ff = num+d;
            if(ff<ans){
                ans = ff;
                ans1 = x;
                ans2 = y;
                ans3 = tmp;
            }
        }
    }
}


void dfs(int x,int num){
    for(int i = 1;;i++){
        int tmp = x*i;
        if(tmp>b){
            if(tmp-b+num>ans) return ;
            dfs1(x,tmp,num+tmp-b);
        }
        else{
            int res = b - tmp;
            dfs1(x,tmp,num+res);
        }
    }
}




int main(){
    IO;
    int t ;
    cin>>t;
    while(t--){
        ans = inf;
        cin>>a>>b>>c;
        for(int i = a;i>0;i--){
            if(a-i>ans)break;
            dfs(i,a-i);
        }
        for(int i = a+1;i<=1e4;i++){
            if(i-a>ans)break;
            dfs(i,i-a);
        }
        cout<<ans<<endl;
        cout<<ans1<<" "<<ans2<<" "<<ans3<<endl;
    }

}