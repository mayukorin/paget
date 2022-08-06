package model

type Paper struct {
	ID             string
	SubmittedMonth string
}

type Papers []Paper

func (p Papers) Len() int {
	return len(p)
}

func (p Papers) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Papers) Less(i, j int) bool {
	return p[i].SubmittedMonth > p[j].SubmittedMonth // 逆
}
