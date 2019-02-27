package gospin

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const (
	defaultStartChar     = "{"
	defaultEndChar       = "}"
	defaultEscapeChar    = "\\"
	defaultDelimiterChar = "|"

	errBracketsNotMatching = "brackets aren't matching"
)

// Config is the struct that holds the configurable settings
// for the spinner
type Config struct {
	// StartChar is the starting character for the option strings,
	// the default is {
	StartChar string

	// EndChar is the ending character for the option strings,
	// the default is }
	EndChar string

	// EscapeChar is the escaping character for the startChars
	// and endChars the default is \
	EscapeChar string

	// DelimiterChar is the delimiter for the option strings the
	// default is |
	DelimiterChar string

	// UseGlobalRand is the setting to use your own rand seed,
	// if false, the `Spin` func will rand seed itself
	UseGlobalRand bool
}

// Spinner is the main struct for the package, it has the functions
// `Spin` and `SpinN` to spin a string
type Spinner struct {
	Config
}

// New returns a new `Spinner` - can also set the config here by passing
// as a variable or just passing nil to use the default values
func New(config *Config) *Spinner {
	if config == nil {
		config = &Config{
			StartChar:     defaultStartChar,
			EndChar:       defaultEndChar,
			DelimiterChar: defaultDelimiterChar,
			EscapeChar:    defaultEscapeChar,
		}
	}

	if config.StartChar == "" {
		config.StartChar = defaultStartChar
	}

	if config.EndChar == "" {
		config.EndChar = defaultEndChar
	}

	if config.DelimiterChar == "" {
		config.DelimiterChar = defaultDelimiterChar
	}

	if config.EscapeChar == "" {
		config.EscapeChar = defaultEscapeChar
	}

	return &Spinner{*config}
}

// Spin does the spinning of a string, if `UseGlobalRand` set to
// false it also sets the seed for the rand
func (s *Spinner) Spin(str string) (string, error) {
	if !s.Config.UseGlobalRand {
		rand.Seed(time.Now().UnixNano())
	}
	running, step := true, 0
	var seq string
	var err error
	for running && err == nil {
		running, err = s.walk(&seq, &step, &str, 0, 0)
		step++
	}

	return seq, err
}

// SpinN spins a string an N amount of times
func (s *Spinner) SpinN(str string, times int) ([]string, error) {
	var seqs []string

	for i := 0; i < times; i++ {
		str, err := s.Spin(str)
		if err != nil {
			return seqs, err
		}
		seqs = append(seqs, str)
	}

	return seqs, nil
}

func (s *Spinner) walk(seq *string, step *int, str *string, start int, level int) (bool, error) {
	if *step >= len(*str) {
		return false, nil
	}

	char := string((*str)[*step])
	prevCharNotEscape := *step == 0 || string((*str)[*step-1]) != s.Config.EscapeChar
	if char == s.Config.StartChar && prevCharNotEscape {
		var err error
		start = *step
		running := true
		level++
		for running && err == nil {
			start++
			running, err = s.walk(seq, &start, str, 0, level)
		}

		selected := s.selectOpt((*str)[*step : start+1])
		if level == 1 {
			if selected == "" {
				// trim due to optional params
				*seq = strings.TrimSpace(*seq)
			} else {
				*seq = *seq + selected
			}
		} else {
			// trim due to optional params
			selected = strings.TrimSpace(selected)

			// replace parameter string e.g. {hello|what} with selectedOpt
			stepDiff := len((*str)[*step:start+1]) - len(selected)
			*str = strings.Replace(*str, (*str)[*step:start+1], selected, 1)
			*step = *step - stepDiff
			start = start - stepDiff
		}

		*step = start
	} else if char == s.EndChar && prevCharNotEscape {
		if level == 0 {
			return false, errors.New(errBracketsNotMatching)
		}

		return false, nil
	} else if level == 0 {
		*seq = *seq + char
	}

	return true, nil
}

func (s *Spinner) selectOpt(strs string) string {
	strs = strs[1 : len(strs)-1]
	split := strings.Split(strs, s.Config.DelimiterChar)

	var str string
	var opts []string
	for _, i := range split {
		if !strings.HasSuffix(i, s.Config.EscapeChar) {
			opts = append(opts, str+i)
			str = ""
		} else {
			str = str + i
		}
	}

	return opts[rand.Int()%len(opts)]
}
