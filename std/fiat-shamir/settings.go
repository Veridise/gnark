package fiatshamir

import (
	"github.com/Veridise/gnark/frontend"
	"github.com/Veridise/gnark/std/hash"
)

type Settings struct {
	Transcript     *Transcript
	Prefix         string
	BaseChallenges []frontend.Variable
	Hash           hash.FieldHasher
}

func WithTranscript(transcript *Transcript, prefix string, baseChallenges ...frontend.Variable) Settings {
	return Settings{
		Transcript:     transcript,
		Prefix:         prefix,
		BaseChallenges: baseChallenges,
	}
}

func WithHash(hash hash.FieldHasher, baseChallenges ...frontend.Variable) Settings {
	return Settings{
		BaseChallenges: baseChallenges,
		Hash:           hash,
	}
}
