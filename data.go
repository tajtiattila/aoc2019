package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/publicsuffix"
)

var (
	IgnoreCache bool
)

func daydataInts(day int) (<-chan int, <-chan error) {
	ch, cherr := make(chan int), make(chan error, 1)
	go func() {
		defer close(ch)
		defer close(cherr)
		rc, err := daydata(day)
		if err != nil {
			cherr <- err
			return
		}

		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			t := strings.TrimSpace(scanner.Text())
			if t == "" {
				continue
			}
			n, err := strconv.Atoi(t)
			if err != nil {
				cherr <- err
				return
			}
			ch <- n
		}
		cherr <- scanner.Err()
	}()
	return ch, cherr
}

func daydata(day int) (io.ReadCloser, error) {
	if !IgnoreCache {
		r, err := os.Open(daydatafn(day))
		if err == nil {
			return r, nil
		}
	}

	if err := downloaddaydata(day); err != nil {
		return nil, err
	}

	return os.Open(daydatafn(day))
}

func daydatafn(day int) string {
	return fmt.Sprintf("data/%d", day)
}

var dataclient struct {
	once sync.Once

	c   *http.Client
	err error
}

const urlprefix = "https://adventofcode.com"

func downloaddaydata(day int) error {
	dataclient.once.Do(func() {
		dataclient.c, dataclient.err = getdataclient()
	})

	if dataclient.err != nil {
		return dataclient.err
	}

	u := fmt.Sprintf("%s/2019/day/%d/input", urlprefix, day)
	log.Printf("Downloading %s", u)
	resp, err := dataclient.c.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll("data", 0777); err != nil {
		return err
	}

	f, err := os.Create(daydatafn(day))
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
