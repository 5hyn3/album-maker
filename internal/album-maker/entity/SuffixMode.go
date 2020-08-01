package entity

type SuffixMode int

const (
	Noting SuffixMode = iota
	MD5
	DateTime
)

var SuffixModes = []SuffixMode{
	Noting,
	MD5,
	DateTime,
}

func (m SuffixMode) String() string {
	switch m {
	case Noting:
		return "noting"
	case MD5:
		return "md5"
	case DateTime:
		return "dateTime"
	default:
		return "nothing"
	}
}

func NewSuffixMode(str string) SuffixMode {
	for _, m := range SuffixModes {
		if m.String() == str {
			return m
		}
	}

	return Noting
}
