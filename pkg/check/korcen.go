// pkg/check/korcen.go

package check

import (
	"encoding/xml"
	"errors"
	"strings"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/fluffy-melli/korcen-go"
)

type Header struct {
	XMLName xml.Name `json:"-" xml:"header"`
	Input   string   `json:"input" xml:"input"`
	Start   string   `json:"replace-front" xml:"replace-front"`
	End     string   `json:"replace-end" xml:"replace-end"`
}

type Respond struct {
	XMLName   xml.Name `json:"-" xml:"respond"`
	Detect    bool     `json:"detect" xml:"detect"`
	Swear     string   `json:"swear" xml:"swear"`
	String    string   `json:"input" xml:"input"`
	NewString string   `json:"output" xml:"output"`
}

type KorcenResult = korcen.CheckInfo

var korcenPool = sync.Pool{
	New: func() interface{} {
		return &KorcenResult{}
	},
}

func freeKorcenResult(result *KorcenResult) {
	*result = KorcenResult{}
	korcenPool.Put(result)
}

var workerPool = NewWorkerPool(10)

func ShutdownWorkerPool() {
	workerPool.Shutdown()
}

var (
	globalLRU = NewShardedLRUCache(16, 64, workerPool)
)

func Korcen(header *Header) (*Respond, error) {
	if header == nil {
		return nil, errors.New("Korcen: header is nil")
	}

	info, ok := globalLRU.Get(header.Input)
	if ok {
		return buildRespond(header, info), nil
	}

	info = korcen.Check(header.Input)
	if info == nil {
		return nil, errors.New("Korcen: korcen.Check returned nil")
	}

	err := globalLRU.Set(header.Input, info)
	if err != nil {
		return nil, err
	}

	return buildRespond(header, info), nil
}

func buildRespond(header *Header, info *KorcenResult) *Respond {
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
	if start < 0 || end > len(text) || start > end {
		return text
	}

	var sb strings.Builder
	sb.WriteString(text[:start])
	if prefix != "" {
		sb.WriteString(prefix)
	}
	sb.WriteString(text[start:end])
	if suffix != "" {
		sb.WriteString(suffix)
	}
	sb.WriteString(text[end:])
	return sb.String()
}

// ---------------------------------------------------------------------
// Actor
// ---------------------------------------------------------------------

type KorcenRequest struct {
	Header *Header
}

type KorcenResponse struct {
	Respond *Respond
	Err     error
}

type KorcenResponseMessage struct {
	Respond *Respond
	Err     error
}

type KorcenActor struct{}

func (k *KorcenActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *KorcenRequest:
		if msg.Header == nil {
			context.Respond(&KorcenResponse{
				Respond: nil,
				Err:     errors.New("KorcenActor: Header is nil"),
			})
			return
		}

		workerPool.Submit(func() {
			resp, err := Korcen(msg.Header)
			context.Send(context.Self(), &KorcenResponseMessage{
				Respond: resp,
				Err:     err,
			})
		})

	case *KorcenResponseMessage:
		context.Respond(&KorcenResponse{
			Respond: msg.Respond,
			Err:     msg.Err,
		})

	default:
	}
}
