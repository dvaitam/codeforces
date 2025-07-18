#include <bits/stdc++.h>
#define in(x) (scanf("%d",&x))
#define in64(x) (scanf("%lld",&x))
#define out(x) (printf("%d",x))
#define out64(x) (printf("%lld",x))
#define en putchar('\n')
#define sp putchar(' ')
#define LOWBIT(x) (x&(-x))
using namespace std ;
///using namespace __gnu_cxx ;
const double eps (1e-7) ;
const int inf(0x3f3f3f3f) ;
const int MAXN = 200010;
const int DEG = 20;
char bg[MAXN] ;
char bgg[5] ;
	int n;
int getCnt()
{
	int cnt (0) ;
	for(int i(0);i<n;i+=3) {
		if(bg[i]!=bgg[0]) ++cnt ;
	}
	for(int i(1);i<n;i+=3) {
		if(bg[i]!=bgg[1]) ++cnt ;
	}
	for(int i(2);i<n;i+=3) {
		if(bg[i]!=bgg[2]) ++cnt ;
	}
	return cnt ;
}

int main()
{
    bgg[0] = 'B' ;bgg[1] = 'G';bgg[2] = 'R' ;bgg[3]=0 ;
	in(n) ;
	scanf("%s",bg) ;
	///RGB
	char bggg[5] ;
	int mins(inf) ;
	int cnt(0) ;
	do
    {
        int tmp(getCnt()) ;
        ///cout<<bgg<<endl;
        ///cout<<tmp<<"***"<<endl ;
        if(tmp < mins) {
			mins = tmp ;
		    memcpy(bggg,bgg,sizeof(bgg)) ;
        }
    }while(next_permutation(bgg,bgg+3));
	out(mins) ;en ;
	///cout<<bggg<<endl;
	for(int i(0);i<n;i+=3) {
		bg[i] = bggg[0] ;
	}
	for(int i(1);i<n;i+=3) {
		bg[i] = bggg[1] ;
	}
	for(int i(2);i<n;i+=3) {
		bg[i] = bggg[2] ;
	}
	puts(bg) ;
	en ;
}