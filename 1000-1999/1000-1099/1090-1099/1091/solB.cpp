#include<bits/stdc++.h>
using namespace std;


int main(){
#ifdef LOCAL
    freopen("in.txt", "r", stdin);
#endif
    /*int t;
    scanf("%d",&t);
    while(t--){}*/
    int n;
    scanf("%d",&n);
    long long ansx=0,ansy=0;
    for(int i=0;i<2*n;i++){
        int x,y;
        scanf("%d%d",&x,&y);
        ansx+=x;
        ansy+=y;
    }
    cout<<ansx/(n)<<" "<<ansy/(n)<<endl;
    return 0;
}