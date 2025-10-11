# FF - Fuzzy Finder

- Example usage:

```golang
baseDirectory := "D:\\Downloads"
targetFilename := "The.Awesome.Show.S01E09.1080p.WEB.H264-WESTSIDEGUNN.mkv"

log.Println("Starting fuzzy search...")

// --- EXECUTION ---
relatedFiles, err := FindRelatedFiles(baseDirectory, targetFilename)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// --- DISPLAY RESULTS ---
fmt.Println("\n--- Search Results ---")
if len(relatedFiles) == 0 {
    fmt.Println("No related directories found.")
} else {
    for dir, files := range relatedFiles {
        fmt.Printf("\n[+] Files in: %s\n", dir)
        if len(files) == 0 {
            fmt.Println("  (No files found in this directory)")
        } else {
            for _, file := range files {
                // Print just the filename for cleaner output
                fmt.Printf("  - %s\n", filepath.Base(file))
            }
        }
    }
}
fmt.Println("\n--- Search Complete ---")

// --- FIND NEXT EPISODE ---
log.Println("\nSearching for the next episode...")
targetInfo, err := ParseEpisodeInfo(targetFilename)
if err != nil {
    log.Fatalf("Could not parse target filename to find next episode: %v", err)
}

nextEp, found := FindNextEpisode(targetInfo, relatedFiles)

// --- DISPLAY NEXT EPISODE RESULT ---
fmt.Println("\n--- Next Episode ---")
if found {
    fmt.Printf("Next episode found!\n")
    fmt.Printf("  Show:    %s\n", nextEp.ShowName)
    fmt.Printf("  Season:  %d\n", nextEp.Season)
    fmt.Printf("  Episode: %d\n", nextEp.Episode)
    fmt.Printf("  File:    %s\n", nextEp.FilePath)
} else {
    fmt.Println("Could not find a subsequent episode in the related files.")
    fmt.Println("You might be on the last available episode.")
}
fmt.Println("\n--- Search Complete ---")
```
