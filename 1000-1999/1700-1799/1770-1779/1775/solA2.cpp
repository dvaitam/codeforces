#include <bits/stdc++.h>

using namespace std;

#define looklols    ios_base :: sync_with_stdio(0); cin.tie(0); cout.tie(0);

#define T           ll T;  cin>>T; while(T--)

#define ll          long long

#define lld         long double

#define ld          long double

#define F           first

#define S           second

#define pb          push_back

#define pf          push_front

#define all(x)      x.begin(),x.end()

#define allr(x)     x.rbegin(),x.rend()

#define ones(x) __builtin_popcountll(x)

#define sin(a) sin((a)*PI/180)

#define cos(a) cos((a)*PI/180)

#define endl        "\n"

const lld pi = 3.14159265358979323846;

const ll N=2e5+2;

const ll MOD = 998244353 , LOG = 25;

/*

  ℒ◎øкℓ☺łṧ

 */



int main () {



    looklols

    T{

        string s,str,str2;

        set<char>ss;

        cin>>s;

        int n=s.size();

        bool ok= false;

        for (int i = 0; i < s.size(); ++i) {

            ss.insert(s[i]);

        }

        if(ss.size()==1){

            cout<<s[0]<<" "<<s.substr(1,n-2)<<" "<<s[n-1]<<endl;



        }

        else{

            for (int i = 1; i < n-1; ++i) {

                if(s[i]=='a'){

                    str=s.substr(0,i);

                    str2=s.substr(i+1);

                    cout<<str<<" "<<'a'<<" "<<str2<<endl;

                    ok= true;

                    break;

                }

            }

            if(!ok){

                cout<<s[0]<<" "<<s.substr(1,n-2)<<" "<<s[n-1]<<endl;

            }

        }

    }



}