// pkg/check/korcen.go

package check

import (
	"github.com/fluffy-melli/korcen-go"
)

type Header struct {
	Input string `json:"input"`
	Start string `json:"replace-front"`
	End   string `json:"replace-end"`
}

type Respond struct {
	Detect    bool   `json:"detect"`
	Swear     string `json:"swear"`
	String    string `json:"input"`
	NewString string `json:"output"`
}

func Korcen(header *Header) *Respond {
	info := korcen.Check(header.Input)
	if !info.Detect {
		return &Respond{
			Detect:    false,
			Swear:     "",
			String:    header.Input,
			NewString: header.Input,
		}
	}

	firstSwear := info.Swear[0]
	formattedMessage := formatMessage(info.NewText, firstSwear.Start, firstSwear.End, header.Start, header.End)

	return &Respond{
		Detect:    true,
		Swear:     firstSwear.Swear,
		String:    header.Input,
		NewString: formattedMessage,
	}
}

func formatMessage(text string, start, end int, prefix, suffix string) string {
	switch {
	case prefix != "" && suffix != "":
		return text[:start] + prefix + text[start:end] + suffix + text[end:]
	case prefix != "":
		return text[:start] + prefix + text[start:end] + text[end:]
	case suffix != "":
		return text[:start] + text[start:end] + suffix + text[end:]
	default:
		return text
	}
}
