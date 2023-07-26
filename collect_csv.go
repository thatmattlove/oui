package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
)

func DownloadCSV(registry *Registry) (fileName string, err error) {
	file, err := os.CreateTemp("", registry.TempFilePattern())
	if err != nil {
		return
	}
	defer file.Close()
	fileName = file.Name()
	res, err := http.Get(registry.URL().String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	_, err = file.Write(b)
	return
}

func ReadCSV(registry *Registry, fileName string, logger *Logger) (results []*VendorDef) {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Err(err)
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	var place int64
	for {
		var row []string
		row, err = reader.Read()
		if err == io.EOF {
			// Exit loop when file is done being read.
			logger.Success("finished parsing vendors from %s registry", registry.Name)
			break
		} else if err != nil {
			logger.Err(err, "failed to read file '%s'", registry.FileName())
		}
		if place == 0 {
			// Ignore header row.
			place++
			continue
		}
		place++
		if len(row) < 3 {
			// Ignore rows that don't conform to expected structure.
			logger.Warn("skipping row %s", row)
			continue
		}
		assignment := strings.TrimSpace(row[1])
		if !strings.Contains(assignment, "/") {
			assignment += "/24"
		}
		organization := row[2]
		org := strings.TrimSpace(organization)
		base, mp, err := macaddr.ParseMACPrefix(assignment)
		if err != nil {
			logger.Err(err, "failed to parse OUI assignment")
			continue
		}
		prefixLen := mp.PrefixLen()
		prefix := fmt.Sprintf("%s/%d", base.String(), prefixLen)
		v := &VendorDef{
			Org:      org,
			Length:   prefixLen,
			Prefix:   prefix,
			Registry: registry.Name,
		}
		results = append(results, v)
	}
	return
}

func CollectAll(p *progress.Progress, logger *Logger) (results []*VendorDef) {
	registries := Registries()
	for _, reg := range registries {
		p.Advance(uint(88 / len(registries)))
		fileName, err := DownloadCSV(reg)
		if err != nil {
			logger.Err(err, "failed to download file '%s'", reg.FileName())
		}
		results = append(results, ReadCSV(reg, fileName, logger)...)
	}
	return
}
