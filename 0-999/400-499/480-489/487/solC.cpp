#include <algorithm>

#include <cmath>

#include <cstring>

#include <cassert>

#include <cstdio>

#include <iostream>

#define MOD 100000000LL



using namespace std;



int n;



bool isprime(int n) {

	bool b=true;

	int i = 2;

	while (i*i<=n && b) {

		b=(n%i!=0);

		i++;

	}

	return b;

}



int raizprimitiva(int p) {

	assert(isprime(p));

	bool es=false;

	int g = 2;

	while(!es) {

		int pot = g;

		int counter=0;

		while (pot!=1) {

			pot=(pot*g)%p;

			counter++;

		}

		if (counter==p-2) return g;

		g++;

	}

	return 0;

}



int main() {

	scanf("%d",&n);

	if(n==1) {

		printf("YES\n1\n");

	} else if(n==4) {

		printf("YES\n1\n3\n2\n4\n");

	} else if (!isprime(n)) {

		printf("NO\n");

	} else if (n==2) {

		printf("YES\n1\n2\n");

	} else {

		printf("YES\n");

		int g=raizprimitiva(n);

		int gg[n+2];

		gg[0]=1;

		for (int i=1;i<=n;i++) {

			gg[i]=(gg[i-1]*g)%n;

		}

		assert((g*gg[n-2])%n==1);

		printf("1\n");

		int a=2;

		int i=0;

		while(i<n/2) {

			int pos = gg[n-a];

			int qos = gg[a];

			if (pos!=1) printf("%d\n",pos);

			if (qos!=1) printf("%d\n",qos);

			a+=2;

			i++;

		}

		printf("%d\n",n);

	}

	return 0;

}