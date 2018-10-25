# Golang go-routine pool (Executor)

Difference with Java, the perpose of golang executor is not save the cost of thread create. We use it for concurrency control.

## Introduction

There are a lot of use cases to control the parallel executing of some work. We design a "Executor" just like it in "Java" to control the concurrency.

This design is not for same the creation of go-routing (which is the most objective in Java). The creation of go-routing is so cheap in golang.  

## Use cases

### Case 1. Limit the CPU usage

In same CPU sensitive case, we want to limit the CPU number we used. If a lot of tasks works in parallel, each of the work slow. But if we running tasks in sequential, each task would be running fast. For this case, we need create a executor with the pool size same with CPU core size. 

### Case 2. Limit the concurrency of IO related access

For some IO related cases (like DB access, CDN access), we want to limit the concurrent number, in this case we need start a fixed number of go-routing to running the work and let other work waiting in the queue.

### Case 3. In case of IO error

In case of some error occurred on IO operation, we expect user doing retry on Executor, then it would slow the rate of following IO opertion to reduce the access rate of IO, it would protect the target service (DB/Cache/CDN etc.)  
