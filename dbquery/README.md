# Performance comparison between prepared statements and straight statements

Experiment to compare the performance of running straight SQL statements vs prepared statements.

Benchmarking the following:
1. Running straight SQL statements in a loop
1. Opening a prepared statement and running queries in a loop in the statement
1. Opening a prepared statement, running a query, closing the statement, then rinse and repeat

Table should be kept very small to minimise latency from disk reads. If at all possible, the network path should also be kept minimal to minimise network latency.
