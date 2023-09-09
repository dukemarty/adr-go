/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	adrexport "github.com/dukemarty/adr-go/export"
	"github.com/dukemarty/adr-go/logic"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve ADRs as website",
	Long: `Serve the managed ADRs as a website.
	
Basically, the page which would be exported as html is served via http.
	
Additionally, endpoints are provided to get a list of the available ADRs
(at 'URL/adrs') and to fetch a single ADR in its original format (at
'URL/adr/INDEX').`,
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetUint16("port")

		logger.Println("Command 'serve' called.")

		http.HandleFunc("/", showMainSiteHandler)
		http.HandleFunc("/adr/", adrHandler)
		http.HandleFunc("/adrs/", adrsHandler)

		serveUrl := fmt.Sprintf("%s:%d", address, port)
		fmt.Printf("Serving at http://%s\n", serveUrl)
		log.Fatal(http.ListenAndServe(serveUrl, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("address", "a", "localhost", "blabla")
	serveCmd.Flags().Uint16P("port", "p", 8080, "blabla")
}

func showMainSiteHandler(w http.ResponseWriter, r *http.Request) {
	dataPath, data := loadAdrData(logger)

	exporter, err := adrexport.CreateExporter(logger, "html")
	if err != nil {
		logger.Printf("Error when creating exporter: %v\n", err)
		return
	}

	exportData := exporter.Export(logger, data, dataPath)

	fmt.Fprint(w, exportData)
}

func adrHandler(w http.ResponseWriter, r *http.Request) {
	// parts := strings.Split(r.URL., ("/"))
	fmt.Fprintf(w, "ADR-URL = %q\n", r.URL)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(404)
		return
	}
	reqIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	adrFile, err := logic.GetAdrFilePathByIndex(reqIndex, logger)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	content, err := os.ReadFile(adrFile)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(content)
}

type adrInfoForRest struct {
	Index int
	Title string
}

func adrsHandler(w http.ResponseWriter, r *http.Request) {
	_, data := loadAdrData(logger)

	parts := strings.Split(r.URL.Path, "/")

	var cmd string
	if len(parts) == 2 || len(parts[2]) == 0 {
		cmd = "LIST"
	} else {
		cmd = strings.ToUpper(parts[2])
	}

	switch cmd {
	case "COUNT":
		fmt.Fprintf(w, "%d", len(data))
	case "LIST":
		v := toRestResult(data)
		res, _ := json.MarshalIndent(v, "", "  ")
		fmt.Fprintf(w, "%s\n", string(res))
	default:
		w.WriteHeader(404)
		// fmt.Fprintf(w, "ADRs-URL = %s\n", parts)
	}
}

func toRestResult(asl []logic.AdrStatus) []adrInfoForRest {
	res := make([]adrInfoForRest, 0)

	for _, as := range asl {
		next := adrInfoForRest{
			Index: as.Index,
			Title: as.Title,
		}
		res = append(res, next)
	}

	return res
}
