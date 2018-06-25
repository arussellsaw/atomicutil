# atomicutil

helpers built on top of sync/atomic primitives

MapStringUint64 is a handy tool that allows you to write to a map[string]uint64 only taking a read lock. this is useful for situations where you're tracking a large number of counters and want to avoid contention on increment operations. This works because when incrementing an existing value inside the map, as the atomic operations on the keys do not mutate the map itself. New values will require a write lock, as the map must be mutated to add the new value, but in the majority of use cases i've found this will only happen a few times on service startup.
