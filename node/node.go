package node

type Node struct {
	Name            string
	Ip              string
	Cores           int    // Total CPU cores
	Memory          int    // Total memory in MB
	Disk            int    // Total disk in MB
	MemoryAllocated int    // Allocated memory in MB
	DiskAllocated   int    // Allocated disk in MB
	Role            string // "manager" or "worker"
	TaskCount       int
}
