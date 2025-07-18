#include<bits/stdc++.h> 
using namespace std;
typedef long long ll;
typedef vector<int> vi;
typedef vector<long long> vll;
typedef pair<int, int> pii;
typedef pair<long long, long long> pll;
void setIO() {
    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
}
void solve() {
    ll a, b; cin >> a >> b;
    int digitnuma = 0, digitnumb = 0;
    vi digita, digitb;
    ll result = 0;
    ll tmpa = a, tmpb = b;
    while(tmpa > 0) { digitnuma++; digita.push_back(tmpa%10); tmpa /= 10; }
    while(tmpb > 0) { digitnumb++; digitb.push_back(tmpb%10); tmpb /= 10; }
    if(digitnuma < digitnumb) {
        for(int i=0; i<digitnuma; i++) result = result * 10 + 9;
        cout << result << "\n"; return;
    }
    else {
        reverse(digita.begin(), digita.end());
        reverse(digitb.begin(), digitb.end());
        for(int diff=0; diff<=9; diff++) {
            for(int l=0; l+diff<=9; l++) {
                int r = l+diff;
                ll front = 0;
                for(int i=0; i<digitnuma; i++) {
                    if(digita[i] == digitb[i]) {
                        if(digita[i] < l || digita[i] > r) break;
                        front = front * 10 + digita[i];
                        if(i == digitnuma-1) { cout << front << "\n"; return; }
                    }
                    else {
                        if(l <= digita[i] && digita[i] <= r) {
                            result = front * 10 + digita[i];
                            for(int j=i+1; j<digitnuma; j++) result = result * 10 + r;
                            if(result >= a) { cout << result << "\n"; return; }
                        }
                        if(l <= digitb[i] && digitb[i] <= r) {
                            result = front * 10 + digitb[i];
                            for(int j=i+1; j<digitnuma; j++) result = result * 10 + l;
                            if(result <= b) { cout << result << "\n"; return; }
                        }
                        for(int k=digita[i]+1; k<=digitb[i]-1; k++) {
                            if(l <= k && k <= r) {
                                result = front * 10 + k;
                                for(int j=i+1; j<digitnuma; j++) result = result * 10 + l;
                                cout << result << "\n"; return;
                            }
                        }
                        break;
                    }
                }
            }
        }
    }
}
int main() {
    setIO();
    int t; cin >> t;
    while(t--) solve();
}