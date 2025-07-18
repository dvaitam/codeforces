#include <bits/stdc++.h>

#define _u(i,a,b) for(int i=(a);i<=(b);++i)
#define _d(i,a,b) for(int i=(a);i>=(b);--i)
#define FIO() ios::sync_with_stdio(0); cin.tie(0); cout.tie(0);
#define FI(task) freopen(task".inp","r",stdin);
#define FO(task) freopen(task".out","w",stdout);
#define fillchar(a, x) memset(a, x, sizeof(a))
#define ll long long
#define ii pair<int,int>
#define fi first
#define se second
#define endl "\n"

using namespace std;

const int N = 1e3+7;
const ll M = 1ll*1e9+7;

void write(){}
template<typename T, typename ...Ts>
void write(const T &first, const Ts &...rest){
    cout << first << ' ';
    write(rest...);
}
template<typename ...Ts>
void writeln(const Ts &... rest){
    write(rest...);
    cout << endl;
}

int readint(){
    int num(0); char c; bool neg(0);
    for(c=getchar();c<'0'||c>'9';c=getchar()) neg|=c=='-';
    for(;c>='0'&&c<='9';c=getchar()) num=num*10+c-48;
    return neg?-num:num;
}

void print(ll num){if(num>9)print(num/10);putchar(num%10+48);}

int n;

bool Case0(string s, bool trace = 0){
    vector<int>tr;
    _u(i,1,n-1){
        if(s[i] == 'B'){
            s[i] = 'W';
            if(s[i+1] == 'B') s[i+1] = 'W';
            else s[i+1] = 'B';
            if(trace) tr.push_back(i);
        }
    }
    _u(i,1,n) if(s[i] == 'B') return 0;
    if(trace){
        cout << tr.size() << endl;
        for(int v:tr) cout << v << " ";
    }
    return 1;
}

bool Case1(string s, bool trace = 0){
    vector<int>tr;
    _u(i,1,n-1){
        if(s[i] == 'W'){
            s[i] = 'B';
            if(s[i+1] == 'B') s[i+1] = 'W';
            else s[i+1] = 'B';
            if(trace) tr.push_back(i);
        }
    }
    _u(i,1,n) if(s[i] == 'W') return 0;
    if(trace){
        cout << tr.size() << endl;
        for(int v:tr) cout << v << " ";
    }
    return 1;
}

#define task "block"

int main(){
    FIO();

    cin >> n;
    string s; cin >> s; s = " " + s;
    if(Case0(s, 0)) Case0(s, 1);
    else{
            if(!Case1(s, 0)) cout << -1;
                else Case1(s, 1);
        }
    return 0;
}