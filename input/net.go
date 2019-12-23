package input

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	urlprefix   = "https://adventofcode.com/2019"
	cachesubdir = "tajtiattila/AdventOfCode/2019"
)

func downloaddaydata(day int) error {
	dataclient.once.Do(func() {
		dataclient.c, dataclient.err = getdataclient()
	})

	if dataclient.err != nil {
		return dataclient.err
	}

	u := fmt.Sprintf("%s/day/%d/input", urlprefix, day)
	if ReportDownload {
		log.Printf("Downloading %s", u)
	}
	resp, err := dataclient.c.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := createFile(dayinputfn(day))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func getdataclient() (*http.Client, error) {
	const fn = ".cookie"
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("Cookie file %q missing: %w", fn, err)
	}
	defer f.Close()

	var cookies []*http.Cookie

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		i := strings.IndexByte(t, '=')
		if i < 0 {
			return nil, fmt.Errorf("invalid line in %q", fn)
		}

		cookies = append(cookies, &http.Cookie{
			Name:  t[:i],
			Value: t[i+1:],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	j, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(urlprefix)
	if err != nil {
		return nil, err
	}
	j.SetCookies(u, cookies)

	return &http.Client{Jar: j}, nil
}
