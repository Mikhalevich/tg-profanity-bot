package cmd

type CMD string

const (
	GetAll         CMD = "getall"
	Add            CMD = "add"
	Remove         CMD = "remove"
	ClearAll       CMD = "clearall"
	RestoreDefault CMD = "restoredefault"
	Rankings       CMD = "rankings"
)

func (c CMD) String() string {
	return string(c)
}
