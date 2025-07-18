#include <bits/stdc++.h>

#define speed ios::sync_with_stdio(0);cin.tie(0);cout.tie(0);

#define ll long long 

#define pb push_back

#define mp make_pair

#define f first

#define s second



using namespace std;



const int N = 2001;



int n, p, cnt = 0, mx = 0;

bool f = 1;



int main(){

    speed;

    cin >> n;

    for(int i = 0;i < n;i++){

        int d;

        cin >> d;

        if(d == 1){

            cnt++;

            mx = max(mx, cnt);

        }else{

            if(f) p = cnt, f = 0;

            cnt = 0;

        }

    }

    mx = max(mx, cnt+p);

    cout << mx;

    return 0;

}