/*

                ____________    ______________       __

               / _________  /\ /_____   _____/\     / /\

              / /\       / /  \\    /  /\    \ \   / /  \

             / /  \_____/ /   / \__/  /  \____\/  / /   /

            / /   /    / /   /    /  /   /       / /   /

           / /   /    / /   /    /  /   /       / /   /

          / /   /    / /   /    /  /   /       / /   /

         / /___/____/ /   /    /  /   /       / /___/________

        /____________/   /    /__/   /       /______________/\

        \            \  /     \  \  /        \              \ \

         \____________\/       \__\/          \______________\/

           ___       ___               ___    __________

          /  /\     /  /\             /  /\  /_______  /\

         /  /__\___/  /  \           /  /  \ \      /  /  \

        /____    ____/   /          /  /   /  \____/  /   /

        \   /   /\   \  /          /  /   /       /  /   /

         \_/   /  \___\/ ___      /  /   /       /  /   /

          /   /   /     /  /\    /  /   /       /  /   /

         /   /   /     /  /__\__/  /   /       /  /___/____

        /___/   /     /___________/   /       /___________/\

        \   \  /      \           \  /        \           \ \

         \___\/        \___________\/          \___________\/

       

          A FAN OF FIZZYDAVID

           

*/



#include<bits/stdc++.h>



#define HEAP priority_queue

#define rep(i,n) for(int i=0;i<(n);i++)

#define per(i,n) for(int i=(n)-1;i>=0;i--)

#define forn(i,l,r) for(int i=(l);i<=(r);i++)

#define nrof(i,r,l) for(int i=(r);i>=(l);i--)

#define mp make_pair

#define mt make_tuple

#define pb push_back

#define X first

#define Y second

#define eps 1e-20

#define pi 3.1415926535897932384626433832795

#define orz int

#define yjz main

#define fizzydavid return

#define ak 0

#define la ;

#define debug puts("OK");

#define SZ(x) (int)x.size()

#define ALL(x) x.begin(),x.end()



using namespace std;



typedef long long LL;

typedef double flt;

typedef vector<int> vi;

typedef pair<int,int> pii;



const int iinf=1e9+7;

const LL linf=1ll<<60;

const flt dinf=1e10;



template <typename T>

inline void scf(T &x)

{

	bool f=0; x=0; char c=getchar();

	while((c<'0' || c>'9') && c!='-') c=getchar();

	if(c=='-') { f=1; c=getchar(); }

	while(c>='0' && c<='9') { x=x*10+c-'0'; c=getchar(); }

	if(f) x=-x; return;

}



template <typename T1,typename T2>

void scf(T1 &x,T2 &y) { scf(x); return scf(y); }



template <typename T1,typename T2,typename T3>

void scf(T1 &x,T2 &y,T3 &z) { scf(x); scf(y); return scf(z); }



//---------------------------head----------------------------



int n,a,d;

flt ans,tmp;



inline flt calc(flt v)

{

	flt t1=(v*1.0)/(a*1.0);

	flt d1=t1*(v*1.0)/2.0;

	if(d1<=d) return (1.0*(d-d1))/(v*1.0)+t1;

	return tmp;

}



orz yjz()

{

	scf(n,a,d);

	tmp=sqrt((d*2.0)/(a*1.0));

	forn(i,1,n)

	{

		int t,v;

		scf(t,v);

		flt tmp=t*1.0+calc(v);

		if(tmp<ans) tmp=ans;

		printf("%.5f\n",tmp);

		ans=tmp;

	}

	fizzydavid ak la

}