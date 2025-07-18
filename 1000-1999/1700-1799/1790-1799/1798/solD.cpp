#pragma GCC optimize("Ofast,unroll-loops")
#include <bits/stdc++.h>
using namespace std;
vector <int> v[2];
int main(){
    int T = 1, kase = 0;
    cin >> T;
    while (T--) {
        int n;
        scanf("%d",&n);
        v[0].clear(),v[1].clear();
        for(int i=1,x;i<=n;i++){
            scanf("%d",&x);
            if(x>=0) v[0].push_back(x);
            else v[1].push_back(x);
        }
        sort(v[0].begin(),v[0].end(),greater<int>());
        sort(v[1].begin(),v[1].end());
        if(v[1].size()==0&&v[0][0]==0){
            puts("NO");
            continue;
        }
        long long now=0;
        int p[2]={};
        puts("YES");
        for(int i=1;i<=n;i++){
            if(now<=0) printf("%d",v[0][p[0]]),now+=v[0][p[0]++];
            else printf("%d",v[1][p[1]]),now+=v[1][p[1]++];
            printf("%c",i==n?'\n':' ');
        }
    }
    return 0;
}