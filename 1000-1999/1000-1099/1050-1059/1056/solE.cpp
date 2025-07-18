#include <unordered_set>
#include <algorithm>
#include <iostream>
#include <fstream>
#include <iomanip>
#include <string>
#include <vector>
#include <bitset>
#include <queue>
#include <cmath>
#include <set>
#include <map>
using namespace std;

int main(){
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    //ifstream cin("digits.in");
    //ofstream cout("digits.out");
    string s, t;
    cin >> s >> t;
    int n = s.size(), m = t.size();
    vector<int> kek(2);
    for(int i = 0; i<n; i++)
        kek[s[i]-'0']++;
    int k = 228;
    vector<unsigned long long> h(m+1), pow_k(m+1);
    pow_k[0] = 1;
    for(int i = 0; i<m; i++)
        pow_k[i+1] = k*pow_k[i];
    h[0] = 0;
    for(int i = 0; i<m; i++)
        h[i+1] = h[i]*k+(t[i]-'a'+1);
    int ans = 0;
    for(int i = (m-1)/kek[0]; i>0; i--){
        if((m-i*kek[0])%kek[1])
            continue;
        int a = i, b = (m-i*kek[0])/kek[1];
        unsigned long long ha = 0, hb = 0;
        bool ia = true, ib = true;
        int it = 0;
        bool kek = true;
        for(int j = 0; j<n; j++){
            if(s[j]=='0'){
                if(ia){
                    ha = (h[it+a]-h[it]*pow_k[a]);
                    ia = false;
                }
                else{
                    if(ha!=(h[it+a]-h[it]*pow_k[a])){
                        kek = false;
                        break;
                    }
                }
                it += a;
            }
            else{
                if(ib){
                    hb = (h[it+b]-h[it]*pow_k[b]);
                    ib = false;
                }
                else{
                    if(hb!=(h[it+b]-h[it]*pow_k[b])){
                        kek = false;
                        break;
                    }
                }
                it += b;
            }
        }
        ans += (kek&&(ha!=hb||a!=b));
    }
    cout << ans;

    return 0;
}
///aaaaa.aaaab..ab
/*
-1 1 -1
0 0 2 4
*/