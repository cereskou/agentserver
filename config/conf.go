package config

// bytes size
const (
	_       = iota
	KB uint = 1 << (10 * iota)
	MB
	GB
	TB
)

//ServerSetting -
type ServerSetting struct {
	Host     string `toml:"host"`
	Port     uint   `toml:"port"`
	LogLevel string `toml:"log"`
}

//AwsSetting -
type AwsSetting struct {
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	Region          string `toml:"region"`
	Token           string
}

//CommonSetting -
type CommonSetting struct {
	BroadcastPort uint   `toml:"broadcast_port"`
	Proxy         string `toml:"proxy"`
	Retry         int    `toml:"retry"`
	Timeout       int    `toml:"timeout"`
}

//FileSetting -
type FileSetting struct {
	ListSize uint `toml:"list_buffer_size"`
	Threads  uint `toml:"thread"`
	PartSize uint `toml:"partsize"`
}

//DistributedSetting -
type DistributedSetting struct {
	ListSize uint `toml:"list_buffer_size"`
	Threads  uint `toml:"thread"`
	PartSize uint `toml:"partsize"`
	JobSize  int  `toml:"job_size"`
}

//AppSetting -
type AppSetting struct {
	Common CommonSetting      `toml:"common"`
	Server ServerSetting      `toml:"server"`
	Aws    AwsSetting         `toml:"aws"`
	File   FileSetting        `toml:"file"`
	Dist   DistributedSetting `toml:"distributed"`
}

//Config -
type Config struct {
	Cmd           string
	Dir           string
	S3            string
	Host          string //Master/Agent hostname
	Port          uint   //Master/Agent port
	LocalHost     string
	AccessKey     string
	SecretKey     string
	Region        string
	Token         string
	Proxy         string
	BroadcastPort uint
	Retry         int
	Timeout       int
	JobID         string
	Log           string
	Env           string
	DryRun        bool
	Distributed   bool
	MasterServer  string
	Md5           bool
	ListSize      uint
	Threads       uint
	PartSize      uint64
	DistListSize  uint
	DistThreads   uint
	DistPartSize  uint64
	DistJobSize   int
	Storage       string
	Wait          bool
}
