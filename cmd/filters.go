package cmd

import (
	"fmt"
	"gfuzz/utils"

	"github.com/spf13/cobra"
)

var filtersCmd = &cobra.Command{
	Use:   "filters",
	Short: "Show filters",
	Long:  `Show all filters for --filter`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.Pcyan("Filters Language (Supported by github.com/Knetic/govaluate)"))
		fmt.Println(utils.Pcyan("\nHelp Manual: https://github.com/Knetic/govaluate/blob/master/MANUAL.md\n"))
		fmt.Println(utils.Pcyan(" - Modifiers:"))
		fmt.Print("   - + - / * & | ^ ** % >> <<\n\n")
		fmt.Println(utils.Pcyan(" - Comparators:"))
		fmt.Print("   > >= < <= == != =~ !~\n\n")
		fmt.Println(utils.Pcyan(" - Logical ops:"))
		fmt.Print("   || &&\n\n")
		fmt.Println(utils.Pcyan(" - Prefixes:"))
		fmt.Print("   ! - ~\n\n")
		fmt.Println(utils.Pcyan(" - Boolean:"))
		fmt.Print("   true false\n\n")
		fmt.Println(utils.Pcyan(" - PlaceHolders:"))
		fmt.Print("   FUZZ FUZnZ\n\n")
		fmt.Println(utils.Pcyan(" - Filters:"))
		fmt.Println("   url                 HTTP request's value")
		fmt.Println("   method              HTTP request's verb")
		fmt.Println("   scheme              HTTP request's scheme")
		fmt.Println("   host                HTTP request's host")
		fmt.Println("   content             HTTP response's content")
		fmt.Println("   req_cookies_<name>  HTTP request's cookie")
		fmt.Println("   res_cookies_<name>  HTTP response's cookie")
		fmt.Println("   req_headers_<name>  HTTP request's headers")
		fmt.Println("   res_headers_<name>  HTTP response's headers")
		fmt.Println("   c|code              HTTP response status code")
		fmt.Println("   h|chars             HTTP response chars")
		fmt.Println("   w|words             HTTP response words")
		fmt.Println("   l|lines             HTTP response lines")

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(filtersCmd)
}
