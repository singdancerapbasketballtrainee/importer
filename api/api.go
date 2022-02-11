package api

// NegtServer 对外接口
type NegtServer interface {
	Ping() error
	ImportRps(date string) error
	ImportFxj() error
}
