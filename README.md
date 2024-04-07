# Pseudorandom number generators

## Description
This is simple application to test implementations of random number generators.

For now it supports LGC, and two ways of generating random numbers in the Cauchy distribution.

By default the app generates 10 000 000 numbers for each generator, prints relevant statistics for the distribution tested and creates a histogram of these numbers.

## LGC
The implementation uses `glibc` parameters and of course tries to provide numbers from the uniform distribution. Generated numbers are positive integers up to $2^{31}$. Statistics are calculated for normalized values the are in range of [0,1].

Theoretical value of the mean should be 0.5, and standard deviation should be $\frac{1}{2\sqrt{3}}$ which is around 0,288675134594813.

Below are statistics calculated by the program.
```
Generating numbers using LGC took 40.918709ms
LGC stats (normalized):
	Mean: 0.499957
	Standard deviation: 0.288711
```

Statistics are quite good and close to theoretical values.

Below is the histogram of generated numbers.

![](./plots/Linear%20Congruential%20Generator.png)

Theoretically the histogram should be flatter but of course the generator is not perfect and some numbers repeat more often than others.

## Cauchy

To generate numbers in the Cauchy distribution two algorithms have been used - classic approach using tangent and an alternative that leverages the probability density function - tangent is not needed here.

`X0` and `Gamma` parameters of the distribution are set to 0 and 1 respectively - they can be adjusted in the code.

