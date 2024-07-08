package keydesc

type KeywordsDescriptor interface {
	// Describe keywords from text
	Describe(string) (string, error)
}

func NewKeywordsDescriptor(name string) KeywordsDescriptor {
	if name == "gpt" {
		return NewGPTClient()
	}
	return nil
}
