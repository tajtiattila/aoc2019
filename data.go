package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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

func daydataInts(day int) ([]int, error) {
	rc, err := daydata(day)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	s := string(data)
	var sv []string
	if strings.IndexByte(s, ',') > 0 {
		sv = strings.Split(s, ",")
	} else {
		sv = strings.Fields(s)
	}

	var ints []int
	for _, s := range sv {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		ints = append(ints, n)
	}

	return ints, nil
}

func mustdaydatastr(day int) string {
	rc := mustdaydata(day)
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatal("error reading data:", err)
	}

	return string(b)
}

func mustdaydata(day int) io.ReadCloser {
	rc, err := daydata(day)
	if err != nil {
		log.Fatal("error getting data:", err)
	}
	return rc
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

func mustprocstr(s, sep string, nwant int, f func(string) error) {
	parts := strings.Split(strings.TrimSpace(s), sep)
	if len(parts) != nwant {
		log.Fatalf("error processing %.100q (sep %q), want %d parts, got %d",
			s, sep, nwant, len(parts))
	}
	for i, part := range parts {
		if err := f(part); err != nil {
			log.Fatalf("error processing %.100q (sep %q) at %d: %v",
				s, sep, i, err)
		}
	}
}
