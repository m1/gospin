package main

import (
	"encoding/json"
	"fmt"
	"github.com/m1/gospin"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootCmd       *cobra.Command
	startChar     string
	endChar       string
	escapeChar    string
	delimiterChar string
	times         int
)

func main() {
	rootCmd = &cobra.Command{
		Run:   spinText,
		Use:   "gospin [text]",
		Short: "GoSpin is a article spinning and spintax engine.",
		Long:  "GoSpin is a fast and configurable article spinning and spintax engine written in Go.",
		Args:  cobra.MinimumNArgs(1),
	}

	rootCmd.PersistentFlags().StringVar(&startChar, "start", "{", "Start char for the spinning engine")
	rootCmd.PersistentFlags().StringVar(&endChar, "end", "}", "End char for the spinning engine")
	rootCmd.PersistentFlags().StringVar(&delimiterChar, "delimiter", "|", "Delimiter char")
	rootCmd.PersistentFlags().StringVar(&escapeChar, "escape", "\\", "Escape char")
	rootCmd.PersistentFlags().IntVar(&times, "times", 1, "How many articles to generate")
	rootCmd.Execute()
}

func spinText(_ *cobra.Command, args []string) {
	cfg := gospin.Config{
		StartChar:     startChar,
		EndChar:       endChar,
		DelimiterChar: delimiterChar,
		EscapeChar:    escapeChar,
	}
	spinner := gospin.New(&cfg)
	if times == 1 {
		spun, err := spinner.Spin(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(spun)
		return
	}

	spun, err := spinner.SpinN(args[0], times)
	js, err := json.Marshal(&spun)
	if err != nil {
		log.Panic(err)
		return
	}

	fmt.Println(string(js))
	return
}
