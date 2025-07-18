#include<bits/stdc++.h>

using namespace std;



int main(){

	int a,b,c,d;

	cin>>a>>b>>c>>d;

	if(a==c){

		printf("%d %d %d %d\n",a-abs(b-d),b,c-abs(b-d),d);

	}

	else if(b==d){

		printf("%d %d %d %d\n",a,b-abs(a-c),c,d-abs(a-c));

	}

	else if(abs(a-c)==abs(b-d)){

		printf("%d %d %d %d\n",a,d,c,b);

	}

	else{

		printf("-1\n");

	}

	return 0;

}