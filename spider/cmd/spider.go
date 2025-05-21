package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

)

func main() {
	urlPtr := flag.String("url", "", "URL du site à crawler")
	destPtr := flag.String("dest", "images", "Dossier de destination des images")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Usage: spider -url <url> [-dest <dossier>]")
		os.Exit(1)
	}

	absDest, err := filepath.Abs(*destPtr)
	if err != nil {
		fmt.Println("Erreur chemin destination:", err)
		os.Exit(1)
	}

	visited := make(map[string]bool)
	err = downloader.DownloadImages(*urlPtr, absDest, visited)
	if err != nil {
		fmt.Println("Erreur lors du téléchargement:", err)
		os.Exit(1)
	}

	fmt.Println("Téléchargement terminé dans:", absDest)
}