package replacer

type staticText struct {
	tmpl string
}

func NewStaticText(tmpl string) *staticText {
	return &staticText{
		tmpl: tmpl,
	}
}

func (st *staticText) Replace(text string) string {
	return st.tmpl
}
