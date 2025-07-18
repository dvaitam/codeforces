#include <bits/stdc++.h>

using namespace std;

int get(char a, char b, char c, char d){
	/*int list[256];
	list[(int)a]=0;
	list[(int)b]=0;
	list[(int)c]=0;
	list[(int)d]=0;
	list[0]=0;
	list[(int)a]++;
	list[(int)b]++;
	list[(int)c]++;
	list[(int)d]++;
	if(a==b or a==c or a==d)list[(int)a]=0;	
	if(b==c or b==d)list[(int)b]=0;
	if(c==d)list[(int)c]=0;
	int ret = list[(int)a]%2 + list[(int)b]%2 + list[(int)c]%2 + list[(int)d]%2;
	return (ret==1?1:(ret/2));*/
	char pa = a;
	
	if(a==b and c==d)return 0;
	if(a==c and b==d)return 0;
	if(a==d and b==c)return 0;
	
	a = b;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	a = c;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	a = d;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	a = pa;
	c = a;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	c = b;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	c = d;
	
	if(a==b and c==d)return 1;
	if(a==c and b==d)return 1;
	if(a==d and b==c)return 1;
	
	return 2;
	
	
	
}

int main(){
	int n;
	char s[2][200005];
	scanf("%d %s %s", &n, s[0], s[1]);
	int ans = 0;
	for(int i=0;i<n/2;i++){
		//printf("%d\t", ans);
		ans+=get(s[0][i],s[1][i],s[0][n-i-1],s[1][n-i-1]);
		//printf("%d %d %c %c %c %c\n", ans, get(s[0][i],s[1][i],s[0][n-i-1],s[1][n-i-1]), s[0][i],s[1][i],s[0][n-i-1],s[1][n-i-1]);
	}
	//printf("%d\n", ans);
	if(n%2==1 and s[0][n/2]!=s[1][n/2])
		ans++;
	
	printf("%d\n", ans);
	return 0;
}