#include<iostream>

#include<math.h>

using namespace std;

int main(){

	

	int n , p , flag1 , flag2, cnt=1;	

	cin >> n >> p;

	string s;

	cin >> s;

	p--;

	if(p == 0){

		flag1= 1;

		flag2 = 0;

	}

	else if(p== n-1){

		flag1=0;

		flag2 = 0;

	}

	else if(p < n-1-p){

		flag1=0;

		flag2=1;

	}

	else{

		flag1=1;

		flag2= 1;

	}

	

	if(flag1== 1){

		cout<<"PRINT "<<s[p];

		cnt= p+1;

		while(cnt<n-1){

			cout<<endl<<"RIGHT"<<endl<<"PRINT "<<s[cnt];

			cnt++;

		}

		if(cnt<n)

			cout<<endl<<"RIGHT"<<endl<<"PRINT "<<s[cnt];

	}

	else{

		cout<<"PRINT "<<s[p];

		cnt= p - 1;

		while(cnt>0){

			cout<<endl<<"LEFT"<<endl<<"PRINT "<<s[cnt];

			cnt--;

		}

		if(cnt>=0)

			cout<<endl<<"LEFT"<<endl<<"PRINT "<<s[cnt];	

	}

	

	if(flag2==1){

		if(flag1==1){

			while(cnt>=p){

				cout<<endl<<"LEFT";

				cnt--;

			}

			while(cnt>0){

				cout<<endl<<"PRINT "<<s[cnt]<<endl<<"LEFT";

				cnt--;

			}

			if(cnt>=0)

			cout<<endl<<"PRINT "<<s[cnt];

		}

		else{

			while(cnt<=p){

				cout<<endl<<"RIGHT";

				cnt++;

			}

			while(cnt<n-1){

				cout<<endl<<"PRINT "<<s[cnt]<<endl<<"RIGHT";

				cnt++;

			}

			if(cnt<n)

			cout<<endl<<"PRINT "<<s[cnt];

		}

		

		

	}

	

	

	

	

	

	

}