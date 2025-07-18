#define MAIN signed main()
#define SPEEDUP ios::sync_with_stdio(0),cin.tie(0),cout.tie(0)
#define USING using namespace std
#define END return 0
#define LL long long
#define ULL unsigned long long
#define LD long double
#define STR string
#define EL '\n'
#define BK break
#define CTN continue
#define INF INT_MAX
#define UINF INT_MIN
#define IN(n) cin >> n
#define OUT(n) cout << n
#define OUTL(n) cout << n << EL
#define FP(i,a,b) for(i = a;i < b;i++)
#define FM(i,a,b) for(i = a;i > b;i--)
#define FL(i,a,b) for(i = a;i <= b;i++)
#define FG(i,a,b) for(i = a;i >= b;i--)

#include <bits/stdc++.h>
USING;
int n,i,j,l,r,p = -1,tar,value,d,a[1001];
bool first;
MAIN
{
	SPEEDUP;
	IN(n);
	a[0] = 0;
	l = n,r = 0;
	FL(i,1,n)
	{
		IN(a[i]);
		if(i < l && a[i] > 0)
			l = i;
		if(i > r && a[i] > 0)
			r = i;    
	}
	while(r > 0){
		for(i = p+2;i <= r && !a[i];i++)
			OUT("AR"),p++;
		value = 1;
		d = 0;
		first = true;
		FL(i,p+2,r){
			if(a[i] > 0){
				value += 4;
				if(first)
					d++;
			}
			else{
				value--;
				first = false;
			}
	    	if(value <= 0)
	    	{
		    	tar = d;
		    	FP(j,0,tar)
					OUT("AR");
		    	OUT('A');
		    	FP(j,0,tar)
					OUT('L');  
		    	OUT('A');
		    	FP(j,p+2,p+d+2)
					if(a[j])
						a[j]--;      
		    	BK;
	    	}
	    } 
	    if(value > 0){
		    tar = r - p - 1;
		    FP(i,0,tar)
		    	OUT("AR");
		    OUT("AL");
		    p = r - 2;
		    if(a[r] == 1)
		    	while(a[p+1] <= 1 && p+1 >= l){
					p--;
					OUT('L');
				}
		    while(a[p+1] > 1 && p+1 >= l){
				p--;
				OUT('L');
			}
		    OUT('A');
		    FL(i,p+2,r)
				if(a[i])
					a[i]--;  
	    }
	    l = n;r = 0;
	    FL(i,1,n){
	    	if((i < l) && (a[i] > 0))
				l = i;
	    	if((i > r) && (a[i] > 0))
				r = i;    
	    }
	}  
	END;
}