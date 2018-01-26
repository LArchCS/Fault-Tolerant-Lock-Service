# Simple-Fault-Tolerant-RPC-Lock-Service
This is a simple fault tolerant system, contains one primary and one backup servers. It tolerates 1 fault, and can be easily scaled up to N servers to tolerate (N - 1) faults. However, the mechanism here is simple, only assumes server crashes without consideration of other possible failures (network failures, etc).
