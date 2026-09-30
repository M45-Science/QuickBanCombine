package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type banDataType struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
}

func main() {
	exe, _ := os.Executable()

	//Help
	if len(os.Args) <= 1 {
		fmt.Println("Usage: " + filepath.Base(exe) + " <file1> <file2> ...")
		fmt.Print("Output file: composite.json\n")
		os.Exit(1)
	}

	if err := combineFiles(os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func combineFiles(files []string) error {
	composite := make([]banDataType, 0)
	indexes := make(map[string]int)
	var reasons [][]string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		var bans []banDataType
		var names []string
		if err := json.Unmarshal(data, &names); err == nil {
			for _, name := range names {
				bans = append(bans, banDataType{UserName: name})
			}
		} else if err := json.Unmarshal(data, &bans); err != nil {
			return fmt.Errorf("parse %s: %w", file, err)
		}
		for _, ban := range bans {
			if ban.UserName == "" {
				continue
			}
			name := strings.ToLower(ban.UserName)
			pos, found := indexes[name]
			if !found {
				pos = len(composite)
				indexes[name] = pos
				composite = append(composite, banDataType{UserName: name})
				reasons = append(reasons, nil)
			}
			if ban.Reason == "" {
				continue
			}
			duplicateReason := false
			for _, reason := range reasons[pos] {
				if strings.EqualFold(reason, ban.Reason) {
					duplicateReason = true
					break
				}
			}
			if !duplicateReason {
				reasons[pos] = append(reasons[pos], ban.Reason)
			}
		}
		log.Printf("Read %d bans from %s.\n", len(bans), file)
	}
	for pos, list := range reasons {
		composite[pos].Reason = strings.Join(list, ", ")
		if len(list) > 1 {
			composite[pos].Reason = "[dup] " + composite[pos].Reason
		}
	}
	data, err := json.MarshalIndent(composite, "", "\t")
	if err != nil {
		return fmt.Errorf("encode composite ban list: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile("composite.json", data, 0644); err != nil {
		return fmt.Errorf("write composite.json: %w", err)
	}
	log.Printf("Wrote banlist (%d) of %d bytes.\n", len(composite), len(data))
	return nil
}
