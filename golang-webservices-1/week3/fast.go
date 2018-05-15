package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		user := &User{}
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false
		browsers := user.Browsers

		for _, browser := range browsers {
			notSeenBefore := false

			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore = true
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore = true
			}

			if isAndroid || isMSIE {
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}

				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
