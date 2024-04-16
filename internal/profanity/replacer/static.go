package replacer

type static struct {
	wildcard string
}

func NewStatic(wildcard string) *static {
	return &static{
		wildcard: wildcard,
	}
}

func (st *static) Replace(text string) string {
	return st.wildcard
}
