package cmd

type CMD string

const (
	GetAll         CMD = "getall"
	Add            CMD = "add"
	Remove         CMD = "remove"
	ViewOrginMsg   CMD = "viewOriginMsg"
	ViewBannedMsg  CMD = "viewBannedMsg"
	Unban          CMD = "unban"
	Clear          CMD = "clear"
	RestoreDefault CMD = "restoredefault"
)

func (c CMD) String() string {
	return string(c)
}
