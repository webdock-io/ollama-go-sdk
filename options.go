package ollama

// Options contains Ollama runtime options keyed by OptionKey constants.
type Options map[OptionKey]any

// OptionKey is an enum-style key for Ollama runtime options.
type OptionKey string

const (
	NumCtx        OptionKey = "num_ctx"
	RepeatLastN   OptionKey = "repeat_last_n"
	RepeatPenalty OptionKey = "repeat_penalty"
	Temperature   OptionKey = "temperature"
	Seed          OptionKey = "seed"
	Stop          OptionKey = "stop"
	NumPredict    OptionKey = "num_predict"
	TopK          OptionKey = "top_k"
	TopP          OptionKey = "top_p"
	MinP          OptionKey = "min_p"
)

func (k OptionKey) String() string {
	return string(k)
}
