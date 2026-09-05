// This program finds every product page on the Aquasana water filter
// category, then visits each product page one at a time, finds every
// PDF linked on that page, and downloads any PDF it has not already
// saved into a local "PDFs" folder. Everything is hard-coded on
// purpose: there are no command-line flags and no concurrency, so the
// code runs top to bottom in a single, easy-to-follow order, and it
// can be safely re-run — anything already saved on disk is skipped.
package main // every Go program that runs on its own is package main

// bring in the standard library pieces this program needs
import (
	"fmt"           // fmt lets us build formatted error messages
	"io"            // io lets us read and copy response bodies
	"log"           // log lets us print timestamped status messages to the screen
	"net/http"      // net/http lets us make web requests
	"net/url"       // net/url lets us combine and clean up web addresses
	"os"            // os lets us create folders and files on disk
	"path/filepath" // path/filepath lets us build local file paths safely
	"regexp"        // regexp lets us search text for patterns
	"sort"          // sort lets us put lists into a predictable order
	"strconv"       // strconv lets us turn text numbers into real numbers
	"strings"       // strings gives us helper functions for text
	"time"          // time lets us pause between requests and measure how long they take
)

// the web page where product listings start
const startingPageWebAddress = "https://www.aquasana.com/water-filter-products/"

// the category identifier Aquasana's website uses internally for this listing
const productCategoryId = "water-filter-products"

// how many products the website shows per page of results
const productsShownPerPage = 12

// how long we will wait for any single web request before giving up
const secondsToWaitForEachRequest = 20 * time.Second

// how long we pause between one request and the next, to be polite to the server
const pauseBetweenRequests = 150 * time.Millisecond

// the identification text we send with every request so the server knows what is asking
const identifyOurselvesAs = "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"

// a safety limit on how many extra listing pages we will ask for, in case something goes wrong
const mostPaginationRequestsAllowed = 50

// the name of the local folder where downloaded PDFs are saved
const pdfsFolderName = "PDFs"

// a pattern that matches a link ending in a product's numeric id, e.g. "...-100236133.html"
var productLinkPattern = regexp.MustCompile(`href="([^"]+-\d{5,}\.html[^"]*)"`)

// a pattern that finds text like "of 67 Results" so we know how many products exist in total
var totalResultsPattern = regexp.MustCompile(`of\s+([\d,]+)\s+Results`)

// a pattern that matches any link ending in ".pdf", uppercase or lowercase
var pdfLinkPattern = regexp.MustCompile(`(?i)href="([^"]+\.pdf)"`)

// oneVisitOutcome holds what happened when we visited a single product page
// and tried to download whatever PDFs were linked from it
type oneVisitOutcome struct {
	webAddress        string        // the product page we visited
	statusCode        int           // the numeric HTTP status the server sent back, e.g. 200
	howLongItTook     time.Duration // how long the page request took
	problem           error         // set if something went wrong instead of getting a response
	pdfsFoundOnPage   int           // how many distinct PDF links were on this page
	pdfsDownloadedNow int           // how many of those PDFs we actually downloaded this run
	pdfsAlreadySaved  int           // how many of those PDFs were already on disk from a previous run
}

// main is where the program starts running
func main() {
	// create one shared web client, with a timeout so a stuck request cannot hang forever
	webClient := &http.Client{Timeout: secondsToWaitForEachRequest}

	// tell the person watching that discovery is starting
	log.Printf("Discovering product URLs starting from %s", startingPageWebAddress)

	// go find every product page URL on the site
	productWebAddresses, err := findAllProductWebAddresses(webClient)
	if err != nil {
		// if discovery failed outright, log why and stop
		log.Printf("Could not discover product URLs: %v", err)
		return // stop the program early
	}
	if len(productWebAddresses) == 0 {
		// if we found nothing, the site's layout may have changed; say so and stop
		log.Println("No product URLs were found — the page structure may have changed.")
		return // stop the program early
	}

	// put the URLs in alphabetical order so the printed list is easy to scan
	sort.Strings(productWebAddresses)

	// print out how many products we found and list them
	log.Printf("Found %d product URLs:", len(productWebAddresses))
	for _, oneWebAddress := range productWebAddresses { // loop over every URL we found
		log.Println(" ", oneWebAddress) // log it indented under the heading
	}

	// make sure the local folder for downloaded PDFs exists before we start
	if err := os.MkdirAll(pdfsFolderName, 0o755); err != nil {
		log.Printf("Could not create the %s folder: %v", pdfsFolderName, err) // explain why we cannot continue
		return                                                                // stop the program early
	}

	// tell the person watching that we are about to visit every page
	log.Printf("Visiting %d product pages one at a time to find and download PDFs...", len(productWebAddresses))

	// visit every product page in order, download any new PDFs, and collect what happened at each one
	visitOutcomes := visitEveryProductPageAndDownloadItsPdfs(webClient, productWebAddresses)

	// print a final pass/fail summary, including PDF totals
	printFinalSummary(visitOutcomes)
}

// findAllProductWebAddresses loads the category page, then keeps loading
// further pages of results from the site's "load more" endpoint until it
// has seen every product, returning the full list of product page URLs.
func findAllProductWebAddresses(webClient *http.Client) ([]string, error) {
	// parse the starting page address so we can safely build more addresses from it
	baseWebAddress, err := url.Parse(startingPageWebAddress)
	if err != nil {
		return nil, fmt.Errorf("the starting web address is invalid: %w", err) // bubble up the problem
	}

	// a set (map with no values) of every unique product URL we have found so far
	uniqueProductWebAddresses := make(map[string]struct{})

	// download the first page of the category listing
	firstPageText, err := downloadPageText(webClient, startingPageWebAddress)
	if err != nil {
		return nil, fmt.Errorf("could not download the starting page: %w", err) // bubble up the problem
	}

	// pull any product links out of the first page and remember them
	addProductLinksFoundOnPage(uniqueProductWebAddresses, baseWebAddress, firstPageText)

	// figure out how many products the site says exist in total, if it told us
	totalProductsAdvertised := readAdvertisedTotalCount(firstPageText)
	if totalProductsAdvertised > 0 {
		log.Printf("The category page reports %d total results.", totalProductsAdvertised) // let the user know
	}

	// build the web address of the endpoint the site uses to load more results
	loadMoreWebAddress := *baseWebAddress                                                          // copy the base address
	loadMoreWebAddress.Path = "/on/demandware.store/Sites-aquasana-Site/default/Search-UpdateGrid" // point it at the "load more" endpoint

	// how many extra pages we have requested so far, used for the safety limit
	extraPagesRequestedSoFar := 0

	// keep asking for more results, 12 at a time, starting right after the first page
	for startIndex := productsShownPerPage; ; startIndex += productsShownPerPage {
		// stop once we have reached the advertised total number of products
		if totalProductsAdvertised > 0 && startIndex >= totalProductsAdvertised {
			break // we are done paginating
		}
		// stop if we have hit the safety limit, so a bug cannot loop forever
		if extraPagesRequestedSoFar >= mostPaginationRequestsAllowed {
			log.Println("Reached the safety limit on extra pages, stopping pagination.")
			break // stop paginating for safety
		}
		extraPagesRequestedSoFar++ // count this page toward the safety limit

		// build the query parameters for this page of results
		queryParameters := url.Values{}                               // start with an empty set of parameters
		queryParameters.Set("cgid", productCategoryId)                // which category to load
		queryParameters.Set("start", strconv.Itoa(startIndex))        // which result to start at
		queryParameters.Set("sz", strconv.Itoa(productsShownPerPage)) // how many results to return
		thisPageWebAddress := loadMoreWebAddress                      // copy the "load more" address
		thisPageWebAddress.RawQuery = queryParameters.Encode()        // attach the query parameters to it

		// download this page of results
		pageText, err := downloadPageText(webClient, thisPageWebAddress.String())
		if err != nil {
			// if one page of pagination fails, stop paginating rather than crashing
			log.Printf("Warning: could not load results starting at %d: %v", startIndex, err)
			break // stop paginating
		}

		// remember how many unique URLs we had before this page, to detect new additions
		countBeforeThisPage := len(uniqueProductWebAddresses)

		// pull any product links out of this page and remember them
		addProductLinksFoundOnPage(uniqueProductWebAddresses, baseWebAddress, pageText)

		// if the site never told us a total, and this page added nothing new, assume we are done
		if totalProductsAdvertised == 0 && len(uniqueProductWebAddresses) == countBeforeThisPage {
			break // no new products were found, so stop paginating
		}

		// pause briefly before the next request so we are not hammering the server
		time.Sleep(pauseBetweenRequests)
	}

	// turn the set of unique URLs into a plain list to return
	allProductWebAddresses := make([]string, 0, len(uniqueProductWebAddresses)) // preallocate the right size
	for oneWebAddress := range uniqueProductWebAddresses {                      // loop over every unique URL
		allProductWebAddresses = append(allProductWebAddresses, oneWebAddress) // add it to the list
	}
	return allProductWebAddresses, nil // hand back the finished list
}

// readAdvertisedTotalCount looks for text like "of 67 Results" in the page
// and returns 67, or returns 0 if it could not find that text.
func readAdvertisedTotalCount(pageText string) int {
	match := totalResultsPattern.FindStringSubmatch(pageText) // try to find the pattern in the page
	if match == nil {
		return 0 // the pattern was not found, so we do not know the total
	}
	numberAsText := strings.ReplaceAll(match[1], ",", "") // remove any thousands separator like "1,234"
	numberAsInt, _ := strconv.Atoi(numberAsText)          // convert the text into a real number
	return numberAsInt                                    // hand back the number
}

// addProductLinksFoundOnPage searches the given page text for product links
// and adds each one, resolved into a full web address, into the given set.
func addProductLinksFoundOnPage(uniqueProductWebAddresses map[string]struct{}, baseWebAddress *url.URL, pageText string) {
	allMatches := productLinkPattern.FindAllStringSubmatch(pageText, -1) // find every matching link on the page
	for _, oneMatch := range allMatches {                                // loop over each match we found
		rawLinkText := unescapeHtmlText(oneMatch[1]) // clean up HTML escape codes like &amp;
		parsedLink, err := url.Parse(rawLinkText)    // parse the link into a structured web address
		if err != nil {
			continue // skip this one if it could not be parsed
		}
		fullWebAddress := baseWebAddress.ResolveReference(parsedLink)   // turn a relative link into a full one
		fullWebAddress.Fragment = ""                                    // drop any "#section" part, it is not needed
		uniqueProductWebAddresses[fullWebAddress.String()] = struct{}{} // add the finished address to the set
	}
}

// findPdfWebAddressesOnPage searches the given page text for PDF links and
// returns each one, resolved into a full, unique, sorted web address.
func findPdfWebAddressesOnPage(pageText string, baseWebAddress *url.URL) []string {
	allMatches := pdfLinkPattern.FindAllStringSubmatch(pageText, -1) // find every PDF link on the page
	uniquePdfWebAddresses := make(map[string]struct{})               // a set to avoid downloading the same PDF twice from one page

	for _, oneMatch := range allMatches { // loop over each match we found
		rawLinkText := unescapeHtmlText(oneMatch[1]) // clean up HTML escape codes like &amp;
		parsedLink, err := url.Parse(rawLinkText)    // parse the link into a structured web address
		if err != nil {
			continue // skip this one if it could not be parsed
		}
		fullWebAddress := baseWebAddress.ResolveReference(parsedLink) // turn a relative link into a full one
		uniquePdfWebAddresses[fullWebAddress.String()] = struct{}{}   // add the finished address to the set
	}

	pdfWebAddresses := make([]string, 0, len(uniquePdfWebAddresses)) // preallocate the right size
	for oneWebAddress := range uniquePdfWebAddresses {               // loop over every unique PDF address
		pdfWebAddresses = append(pdfWebAddresses, oneWebAddress) // add it to the list
	}
	sort.Strings(pdfWebAddresses) // sort them so the download order is predictable
	return pdfWebAddresses        // hand back the finished list
}

// unescapeHtmlText turns common HTML escape codes back into plain characters,
// for example turning "&amp;" back into "&".
func unescapeHtmlText(text string) string {
	replacer := strings.NewReplacer( // build a set of find-and-replace rules
		"&amp;", "&", // ampersand
		"&quot;", `"`, // double quote
		"&#39;", "'", // single quote
		"&lt;", "<", // less-than sign
		"&gt;", ">", // greater-than sign
	)
	return replacer.Replace(text) // apply all the rules and return the result
}

// visitEveryProductPageAndDownloadItsPdfs visits each product URL, one after
// another, finds and downloads any new PDFs linked from it, logs a line for
// each page as it finishes, and returns what happened at every page.
func visitEveryProductPageAndDownloadItsPdfs(webClient *http.Client, productWebAddresses []string) []oneVisitOutcome {
	// somewhere to collect the outcome of every page we visit
	allOutcomes := make([]oneVisitOutcome, 0, len(productWebAddresses)) // preallocate the right size

	// go through the product URLs one at a time, in order
	for _, oneWebAddress := range productWebAddresses {
		outcome := visitOneProductPageAndDownloadItsPdfs(webClient, oneWebAddress) // visit this single page and download its PDFs

		resultLabel := "OK" // assume success unless we learn otherwise
		if outcome.problem != nil || outcome.statusCode >= 400 {
			resultLabel = "FAIL" // mark it as a failure if there was an error or a bad status code
		}

		// log one line describing what happened at this page, including the PDF counts
		log.Printf("[%-4s] %-3d %8s  %s  (pdfs: %d found, %d downloaded, %d already saved)",
			resultLabel, outcome.statusCode, outcome.howLongItTook.Round(time.Millisecond), outcome.webAddress,
			outcome.pdfsFoundOnPage, outcome.pdfsDownloadedNow, outcome.pdfsAlreadySaved)

		allOutcomes = append(allOutcomes, outcome) // remember this outcome for the final summary

		time.Sleep(pauseBetweenRequests) // pause briefly before visiting the next page
	}

	return allOutcomes // hand back every outcome we recorded
}

// visitOneProductPageAndDownloadItsPdfs downloads a single product page,
// finds every PDF linked on it, downloads any that are not already saved,
// and reports back what happened.
func visitOneProductPageAndDownloadItsPdfs(webClient *http.Client, productWebAddress string) oneVisitOutcome {
	// start building the outcome we will return, with the web address already filled in
	outcome := oneVisitOutcome{webAddress: productWebAddress}

	request, err := http.NewRequest(http.MethodGet, productWebAddress, nil) // build the outgoing request
	if err != nil {
		outcome.problem = err // record that the request could not even be built
		return outcome        // nothing more we can do for this page
	}
	request.Header.Set("User-Agent", identifyOurselvesAs) // attach our identifying header

	startTime := time.Now()                       // remember when we started, to measure duration
	response, err := webClient.Do(request)        // actually send the request and wait for a response
	outcome.howLongItTook = time.Since(startTime) // work out how long that took

	if err != nil {
		outcome.problem = err // record the problem
		return outcome        // nothing more we can do for this page
	}
	defer response.Body.Close() // make sure we close the response body when this function ends

	outcome.statusCode = response.StatusCode // remember the status code either way

	bodyBytes, err := io.ReadAll(response.Body) // read the whole page into memory
	if err != nil {
		outcome.problem = err // record the problem
		return outcome        // nothing more we can do for this page
	}

	if response.StatusCode >= 400 {
		outcome.problem = fmt.Errorf("got status %s", response.Status) // treat error statuses as failures
		return outcome                                                 // do not look for PDFs on an error page
	}

	pageText := string(bodyBytes) // convert the page body to text so we can search it

	baseWebAddress, err := url.Parse(productWebAddress) // parse this page's own address, used to resolve relative PDF links
	if err != nil {
		outcome.problem = err // record the problem
		return outcome        // cannot safely resolve PDF links without a valid base address
	}

	pdfWebAddresses := findPdfWebAddressesOnPage(pageText, baseWebAddress) // find every PDF linked from this page
	outcome.pdfsFoundOnPage = len(pdfWebAddresses)                         // remember how many we found

	// go through every PDF link found on this page, one at a time
	for _, onePdfWebAddress := range pdfWebAddresses {
		wasAlreadySaved, downloadErr := downloadPdfUnlessAlreadySaved(webClient, onePdfWebAddress) // download it, unless we already have it
		if downloadErr != nil {
			log.Printf("  Warning: could not download %s: %v", onePdfWebAddress, downloadErr) // log the problem but keep going
			continue                                                                          // move on to the next PDF
		}
		if wasAlreadySaved {
			outcome.pdfsAlreadySaved++ // count it as already saved
		} else {
			outcome.pdfsDownloadedNow++      // count it as newly downloaded
			time.Sleep(pauseBetweenRequests) // pause briefly, since this was a real download request
		}
	}

	return outcome // hand back everything that happened on this page
}

// downloadPdfUnlessAlreadySaved works out the local file name for a PDF web
// address, skips it if that file already exists in the PDFs folder, and
// otherwise downloads it and saves it there. It reports whether the file
// was already saved (true) or was freshly downloaded (false).
func downloadPdfUnlessAlreadySaved(webClient *http.Client, pdfWebAddress string) (bool, error) {
	localFileName := localFileNameForPdfWebAddress(pdfWebAddress) // work out what to call the file on disk
	if localFileName == "" {
		return false, fmt.Errorf("could not work out a file name for %s", pdfWebAddress) // give up on this one
	}
	localFilePath := filepath.Join(pdfsFolderName, localFileName) // build the full local path, e.g. "PDFs/WH-1000_Install.pdf"

	if _, statErr := os.Stat(localFilePath); statErr == nil {
		return true, nil // the file is already there, so skip downloading it again
	}

	request, err := http.NewRequest(http.MethodGet, pdfWebAddress, nil) // build the outgoing request for the PDF
	if err != nil {
		return false, err // report the problem if it could not be built
	}
	request.Header.Set("User-Agent", identifyOurselvesAs) // attach our identifying header

	response, err := webClient.Do(request) // send the request and wait for a response
	if err != nil {
		return false, err // report the problem if the request failed
	}
	defer response.Body.Close() // make sure we close the response body when this function ends

	if response.StatusCode >= 400 {
		io.Copy(io.Discard, response.Body)                                               // drain the body so the connection can be reused
		return false, fmt.Errorf("got status %s for %s", response.Status, pdfWebAddress) // report the bad status
	}

	outputFile, err := os.Create(localFilePath) // create the local file to write the PDF into
	if err != nil {
		return false, err // report the problem if the file could not be created
	}
	defer outputFile.Close() // make sure we close the file when this function ends

	if _, err := io.Copy(outputFile, response.Body); err != nil {
		return false, err // report the problem if writing the file failed
	}

	log.Printf("  Downloaded %s", localFilePath) // let the person watching know a new file was saved
	return false, nil                            // this one was freshly downloaded, not already saved
}

// localFileNameForPdfWebAddress takes a full PDF web address and returns
// just the file name at the end of it, e.g. "WH-1000_Install.pdf".
func localFileNameForPdfWebAddress(pdfWebAddress string) string {
	parsedWebAddress, err := url.Parse(pdfWebAddress) // parse the address so we can look at just its path
	if err != nil {
		return filepath.Base(pdfWebAddress) // fall back to a simple guess if parsing fails
	}
	return filepath.Base(parsedWebAddress.Path) // return the last part of the path, ignoring any query string
}

// printFinalSummary prints how many pages succeeded, how many PDFs were
// downloaded and skipped, and lists any failures.
func printFinalSummary(allOutcomes []oneVisitOutcome) {
	successCount := 0          // how many pages loaded successfully
	failureCount := 0          // how many pages failed or errored
	totalPdfsDownloaded := 0   // how many PDFs were freshly downloaded across every page
	totalPdfsAlreadySaved := 0 // how many PDFs were already on disk across every page

	log.Println("Summary") // log a heading
	log.Println("-------") // log an underline for the heading

	for _, outcome := range allOutcomes { // loop over every outcome we recorded
		totalPdfsDownloaded += outcome.pdfsDownloadedNow  // add this page's fresh downloads to the running total
		totalPdfsAlreadySaved += outcome.pdfsAlreadySaved // add this page's already-saved PDFs to the running total

		if outcome.problem != nil {
			failureCount++                                                       // count this as a failure
			log.Printf("  ERROR  %s -> %v", outcome.webAddress, outcome.problem) // describe the error
			continue                                                             // move on to the next outcome
		}
		if outcome.statusCode >= 400 {
			failureCount++                                                 // count this as a failure too
			log.Printf("  %d  %s", outcome.statusCode, outcome.webAddress) // describe the bad status code
			continue                                                       // move on to the next outcome
		}
		successCount++ // this page loaded fine
	}

	log.Printf("%d/%d product pages returned a successful status.", successCount, len(allOutcomes)) // log the overall score
	if failureCount > 0 {
		log.Printf("%d page(s) failed or returned an error status — see above.", failureCount) // call out failures if any
	}
	log.Printf("%d PDF(s) downloaded this run, %d PDF(s) were already saved from before.", totalPdfsDownloaded, totalPdfsAlreadySaved) // report the PDF totals
}

// downloadPageText downloads the given web address and returns its body as text.
func downloadPageText(webClient *http.Client, webAddress string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, webAddress, nil) // build the outgoing request
	if err != nil {
		return "", err // report the problem if it could not be built
	}
	request.Header.Set("User-Agent", identifyOurselvesAs) // attach our identifying header

	response, err := webClient.Do(request) // send the request and wait for a response
	if err != nil {
		return "", err // report the problem if the request failed
	}
	defer response.Body.Close() // make sure we close the response body when this function ends

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("got status %s for %s", response.Status, webAddress) // treat error statuses as failures
	}

	bodyBytes, err := io.ReadAll(response.Body) // read the whole response body into memory
	if err != nil {
		return "", err // report the problem if reading failed
	}
	return string(bodyBytes), nil // convert the bytes to text and return it
}
