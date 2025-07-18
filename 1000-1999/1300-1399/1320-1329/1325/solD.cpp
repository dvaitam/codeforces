// LUOGU_RID: 102435206
#include<bits/stdc++.h>
using namespace std;
long long u,v,c;
int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    cin>>u>>v;
    c=v-u;
    if(c<0||(c&1)){
        cout<<-1;
        return 0;
    }
    if(!c){
        if(!u)
            cout<<0;
        else
            cout<<1<<endl<<u;
    }
    else{
        c>>=1;
        if(!(c&u))
            cout<<2<<endl<<c<<' '<<(c^u);
        else
            cout<<3<<endl<<c<<' '<<c<<' '<<u;
    }
    return 0;
}