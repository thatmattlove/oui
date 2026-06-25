package oui

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
)

func DownloadCSV(registry *Registry) (string, error) {
	client := http.DefaultClient
	client.Timeout = 30 * time.Second

	req, err := http.NewRequest(http.MethodGet, registry.URL().String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("user-agent", "oui")
	res, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return "", fmt.Errorf("request timed out: %w", err)
		}
		return "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("[%s] failed to download data from %s error: [%s] %v", registry.Name, registry.URL(), res.Status, string(b))
	}
	file, err := os.CreateTemp("", registry.TempFilePattern())
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(b)
	return file.Name(), nil
}

func ReadCSV(registry *Registry, fileName string, logger LoggerType) ([]*VendorDef, error) {
	results := make([]*VendorDef, 0)
	file, err := os.Open(fileName)
	if err != nil {
		if logger != nil {
			logger.Err(err)
		}
		return nil, err
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
			if logger != nil {
				logger.Success("finished parsing vendors from %s registry", registry.Name)
			}
			break
		} else if err != nil {
			if logger != nil {
				logger.Err(err, "failed to read file '%s'", registry.FileName())
			}
		}
		if place == 0 {
			// Ignore header row.
			place++
			continue
		}
		place++
		if len(row) < 3 {
			// Ignore rows that don't conform to expected structure.
			if logger != nil {
				logger.Warn("skipping row %s", row)
			}
			continue
		}
		assignment := strings.TrimSpace(row[1])
		if !strings.Contains(assignment, "/") {
			assignment += fmt.Sprintf("/%d", registry.DefaultPrefixLen)
		}
		organization := row[2]
		org := strings.TrimSpace(organization)
		base, mp, err := macaddr.ParseMACPrefix(assignment)
		if err != nil {
			if logger != nil {
				logger.Err(err, "failed to parse OUI assignment")
			}
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
	return results, nil
}

func CollectAll(p *progress.Progress, logger LoggerType) ([]*VendorDef, error) {
	registries := Registries()
	defs := make([]*VendorDef, 0)
	errs := make([]error, 0)
	for _, reg := range registries {
		if p != nil {
			p.Advance(uint(88 / len(registries)))
		}
		fileName, err := DownloadCSV(reg)
		if err != nil {
			errs = append(errs, err)
			if logger != nil {
				logger.Err(err, "failed to download file '%s'", reg.FileName())
			}
		}
		results, err := ReadCSV(reg, fileName, logger)
		if err != nil {
			return nil, err
		}
		defs = append(defs, results...)
	}
	err := errors.Join(errs...)
	return defs, err
}
