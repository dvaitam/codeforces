#include<bits/stdc++.h>

using namespace std;



#define f first

#define s second

#define mp make_pair

#define pb push_back

#define sqr(a) ((a)*(a))

#define v vector

#define p pair

#define ll long long

#define ld long double

#define ull unsigned long long

#define fixed(a) cout.precision((a));cout<<fixed;



int a,b,c;

ld x1,x2,d;



int main()

{

    cin>>a>>b>>c;



    fixed(6)



    d=sqrt(sqr(b)-4*a*c);



    x1=(-b+d)/(2*a);

    x2=(-b-d)/(2*a);



    cout<<max(x1,x2)<<"\n"<<min(x1,x2);

}