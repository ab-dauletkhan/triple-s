package core

const (
	DirPerm  = 0o755
	FilePerm = 0o644

	BucketsFile = "buckets.csv"
	ObjectsFile = "objects.csv"
)

var (
	BucketsCSVHeader = []string{"Name", "Status", "CreationDate", "LastUpdated"}
	ObjectsCSVHeader = []string{"ObjectKey", "ContentType", "ContentLength"}
)
