#include<bits/stdc++.h>
using namespace std;
const int inf = 1e9;
int f[3];
array <int,3> exgcd(int a,int b){
    if(b == 0)return {1,0,a};
    array <int,3> nr = exgcd(b,a%b);
    ///nr[0]*b + nr[1]*a - nr[1]*(a/b)*b = nr[2]
    ///nrn[0]*a + nrn[1]*b = nrn[2]
    return {nr[1],nr[0] - nr[1]*(a/b),nr[2]};
}
array <int,3> ans2;
array <int,3> p;
int ans = inf;
int step1,step2;
void pain(int r2,int i){
    for(int r = max(r2 - 5,0);r <= r2 + 5;r++){
        /*if(i == 246){
            cout<<r<<' '<<(p[0] + r*step1)*f[0] - i*f[1]<<' '<<i*f[1] - (p[1] - r*step2)*f[2]<<' '<<(i*f[1] - p[0]*f[0])/(step1*f[0])<<'\n';
        }*/
        if((p[0] + r*step1) <= i && i <= (p[1] - r*step2) && ans > abs((p[0] + r*step1)*f[0] - i*f[1]) + abs(i*f[1] - (p[1] - r*step2)*f[2])){
            ans = abs((p[0] + r*step1)*f[0] - i*f[1]) + abs(i*f[1] - (p[1] - r*step2)*f[2]);
            ans2 = {p[0] + r*step1,i,p[1] - r*step2};
        };
    }
}
int main(){
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    int n,s;
    cin>>n>>s;
    for(int i = 0;i < n;i++){
        int a;
        cin>>a;
        f[a - 3]++;
    }
    for(int i = 0;i*f[1] <= s;i++){
        ///f[0],f[2]
        p = exgcd(f[0],f[2]);
        s-=i*f[1];
        if(s%p[2] == 0){
            p[0]*=s/p[2];
            p[1]*=s/p[2];
            int lcm = f[0]*f[2]/p[2];
            int x = -p[0]/(lcm/f[0]);
            p[0]+=lcm/f[0]*x;
            p[1]-=lcm/f[2]*x;
            if(p[0] < 0){
                p[0]+=lcm/f[0];
                p[1]-=lcm/f[2];
            }
            //cout<<"cerc:\n";
            //cout<<p[0]<<' '<<i<<' '<<p[1]<<'\n';
            step1 = lcm/f[0],step2 = lcm/f[2];
            pain(0,i);
            pain((p[1] - i)/step2,i);
            pain((i - p[0])/step1,i);
            pain((i*f[1] - p[0]*f[0])/(step1*f[0]),i);
            pain((i*f[1] - p[1]*f[2])/(step2*f[2]),i);
            /*int r = 0;
            while((p[0] + r*step1) <= i && i <= (p[1] - r*step2)){
                if((p[0] + r*step1) <= i && i <= (p[1] - r*step2) && ans > abs((p[0] + r*step1)*f[0] - i*f[1]) + abs(i*f[1] - (p[1] - r*step2)*f[2])){
                    ans = abs((p[0] + r*step1)*f[0] - i*f[1]) + abs(i*f[1] - (p[1] - r*step2)*f[2]);
                    ans2 = {p[0] + r*step1,i,p[1] - r*step2};
                };
                r++;
            }*/
        }
        s+=i*f[1];
    }
    if(ans == inf)cout<<-1;
    else cout<<ans2[0]<<' '<<ans2[1]<<' '<<ans2[2];
    return 0;
}
/**
10 2950
4 3 4 5 4 4 3 3 3 3
0 246 1966
**/