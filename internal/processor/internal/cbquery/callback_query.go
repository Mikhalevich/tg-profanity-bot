package cbquery

type CBQUERY string

const (
	Add           CBQUERY = "add"
	Remove        CBQUERY = "remove"
	ViewOrginMsg  CBQUERY = "viewOriginMsg"
	ViewBannedMsg CBQUERY = "viewBannedMsg"
	Unban         CBQUERY = "unban"
)

func (c CBQUERY) String() string {
	return string(c)
}
