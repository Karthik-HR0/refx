package main

import (
	
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const (
	reflectionMarker = "reflect_test_parameter"
	maxWorkers       = 20
)

type CrawlResult struct {
	URLs        map[string]bool
	Parameters  map[string][]string
	TargetFolder string
}

func printBanner() {
	banner := `
 
              ___________       
_______   ____\_   _____/__  ___
\_  __ \_/ __ \|    __) \  \/  /
 |  | \/\  ___/|     \   >    <     
 |__|    \___  >___  /  /__/\_ \
             \/    \/         \/    
                                    V1.0
                                    @Karthik-HR0
                                    
    Automated Reflected Parameter Finder Tool
                       
	`
	color.Cyan(banner)
}

func fetchURL(targetURL string) (string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func isInternalURL(urlStr string, targetDomain string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return strings.HasSuffix(parsedURL.Hostname(), targetDomain)
}

func crawlDomain(target string, crawlSubdomains bool) (*CrawlResult, error) {
	color.Yellow("[*] Crawling the domain for pages and parameters...")

	crawledURLs := make(map[string]bool)
	parameters := make(map[string][]string)
	toVisit := []string{target}

	parsedTarget, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	targetDomain := parsedTarget.Hostname()

	// Create target folder
	targetFolder := strings.ReplaceAll(targetDomain, ".", "_")
	if err := os.MkdirAll(targetFolder, 0755); err != nil {
		return nil, err
	}

	for len(toVisit) > 0 {
		currentURL := toVisit[0]
		toVisit = toVisit[1:]

		if crawledURLs[currentURL] {
			continue
		}
		crawledURLs[currentURL] = true

		// Fetch URL content
		response, err := fetchURL(currentURL)
		if err != nil {
			continue
		}

		// Save crawled page
		pageFilename := filepath.Join(targetFolder, strings.ReplaceAll(parsedTarget.Path, "/", "_"))
		if pageFilename == "" {
			pageFilename = "index.html"
		}
		if err := os.WriteFile(filepath.Join(targetFolder, pageFilename), []byte(response), 0644); err != nil {
			return nil, err
		}

		// Parse links and parameters
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
		if err != nil {
			continue
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}

			fullURL, err := url.Parse(href)
			if err != nil {
				return
			}
			absoluteURL := parsedTarget.ResolveReference(fullURL).String()

			// Crawl conditions
			if crawlSubdomains || isInternalURL(absoluteURL, targetDomain) {
				toVisit = append(toVisit, absoluteURL)
			}

			// Extract GET parameters
			parsedURL, err := url.Parse(absoluteURL)
			if err != nil {
				return
			}
			for param := range parsedURL.Query() {
				baseURL := strings.Split(absoluteURL, "?")[0]
				if _, exists := parameters[baseURL]; !exists {
					parameters[baseURL] = []string{}
				}
				if !sliceContains(parameters[baseURL], param) {
					parameters[baseURL] = append(parameters[baseURL], param)
				}
			}
		})
	}

	return &CrawlResult{
		URLs:        crawledURLs,
		Parameters:  parameters,
		TargetFolder: targetFolder,
	}, nil
}

func sliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func checkReflectedParameter(baseURL string, param string) string {
	testValue := reflectionMarker
	fullURL := fmt.Sprintf("%s?%s=%s", baseURL, param, testValue)

	resp, err := http.Get(fullURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	if strings.Contains(string(body), testValue) {
		return fullURL
	}
	return ""
}

func main() {
	printBanner()

	// Parse flags
	target := flag.String("t", "", "Target domain to crawl (e.g., http://example.com)")
	crawlSubdomains := flag.Bool("s", false, "Crawl subdomains as well")
	flag.Parse()

	if *target == "" {
		color.Red("[!] Target URL is required . use -h for help")
		os.Exit(1)
	}

	// Crawl domain
	crawlResult, err := crawlDomain(*target, *crawlSubdomains)
	if err != nil {
		color.Red("[!] Crawling error {Use http:// or https:// in domain & in  urls }: %v", err)
		os.Exit(1)
	}

	color.Yellow("[*] Crawled %d unique pages", len(crawlResult.URLs))
	color.Yellow("[*] Found %d unique parameters", len(crawlResult.Parameters))

	// Test reflected parameters
	color.Yellow("[*] Testing for reflected parameters...")
	var reflectedResults []string
	var wg sync.WaitGroup
	resultChan := make(chan string, 100)

	for baseURL, params := range crawlResult.Parameters {
		for _, param := range params {
			wg.Add(1)
			go func(baseURL, param string) {
				defer wg.Done()
				result := checkReflectedParameter(baseURL, param)
				if result != "" {
					resultChan <- result
				}
			}(baseURL, param)
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		reflectedResults = append(reflectedResults, result)
	}

	// Output results
	if len(reflectedResults) > 0 {
		color.Green("\n[+] Reflected Parameters Found:")
		for _, result := range reflectedResults {
			color.Green("[Reflected] %s", result)

			// Save to file
			resultFile := filepath.Join(crawlResult.TargetFolder, "reflected_parameters.txt")
			f, err := os.OpenFile(resultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			 if err == nil {
				defer f.Close()
				if _, err := f.WriteString(result + "\n"); err != nil {
					color.Red("[!] Error writing to file: %v", err)
				}
			}
		}
	} else {
		color.Red("\n[-] No reflected parameters found.")
	}
}
