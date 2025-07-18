#include <bits/stdc++.h>
using namespace std;

#define FasterIO ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0)
#define MOD 1000000007

int main(){
    FasterIO;

    string s;

    while(cin>>s){
        long long ans = 1;
        long long cnt = 0;

        s = 'b' + s;

        for(int i=s.size()-1; i>=0; i--){
            if(s[i]=='a')
                cnt++;
            else if(s[i]=='b'){
                cnt++;
                ans = (ans * cnt)%MOD;
                cnt=0;
            }
        }

        cout<<ans-1<<endl;
    }

    return 0;
}