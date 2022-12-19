package mnemonic

import "strings"

type PhraseSequence []string

func (p PhraseSequence) String() string {
	var s string
	for _, phrase := range p {
		s += phrase + " "
	}
	return strings.TrimSpace(s)
}

func (p PhraseSequence) Len() int {
	count := 0
	for _, phrase := range p {
		if phrase == "*" {
			continue
		}
		count += 1
	}
	return count
}

func (p PhraseSequence) Fill(s string) PhraseSequence {
	a := make(PhraseSequence, len(p))
	copy(a, p)
	for i, phrase := range p {
		if phrase == "*" {
			a[i] = s
			return a
		}
	}
	return append(a, s)
}

func Pharse(phrases string) PhraseSequence {
	phrases = strings.TrimSpace(phrases)
	return PhraseSequence(strings.Split(phrases, " "))
}
