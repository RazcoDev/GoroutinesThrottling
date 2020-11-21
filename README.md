# GoroutinesThrottling

GoroutinesThrottling is a Golang command-line tool to demonstrate goroutines throttling.

This is done by choosing the max goroutines at once, and insert the amount of jobs as seconds in the cli arguments.

## Installation

```
go build -i -o GoroutinesThrottling
```


## Usage

```
./GoroutinesThrottling -c=2 5 3 5 6 7 4 3 2 6 7 4
```
