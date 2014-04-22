package main

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/codegangsta/cli"
    "io"
    "os"
    "regexp"
    "strconv"
    "strings"
    "unicode/utf8"
)

func is_valid_isbn_10(s string) (bool, string) {
    s = strings.Replace(s, "-", "", -1)
    if len(s) != 10 {
        return false, s
    }
    sum := 0
    for i, c := range s {
        var value int
        if c == 'X' {
            value = 10
        } else {
            value = int(c)
        }
        sum += (10 - i) * value
    }
    return sum%11 == 0, s
}

// (10 - (sum(int(digit) * (3 if idx % 2 else 1)
//        for idx, digit in enumerate(isbn[:12])) % 10)) % 10
func is_valid_isbn_13(s string) (bool, string) {
    s = strings.Replace(s, "-", "", -1)
    if len(s) != 13 {
        return false, s
    }
    sum := 0
    var factor int
    buf := make([]byte, 1)
    for i, c := range s[:12] {
        if i%2 == 0 {
            factor = 1
        } else {
            factor = 3
        }
        _ = utf8.EncodeRune(buf, c)
        value, _ := strconv.Atoi(string(buf))
        sum += value * factor
    }
    check := strconv.Itoa((10 - (sum % 10)) % 10)
    return check == string(s[12]), s
}

// assume isbn is in 10 format already
func isbn_10_to_13(s string) string {
    var buffer bytes.Buffer
    buffer.WriteString("978")
    buffer.WriteString(s[:9])

    // s contains the ISBN13 w/o checksum
    s = buffer.String()
    // checksum calculation
    sum := 0
    var factor int
    buf := make([]byte, 1)
    for i, c := range s[:12] {
        if i%2 == 0 {
            factor = 1
        } else {
            factor = 3
        }
        _ = utf8.EncodeRune(buf, c)
        value, _ := strconv.Atoi(string(buf))
        sum += value * factor
    }
    check := strconv.Itoa((10 - (sum % 10)) % 10)
    buffer.WriteString(check)
    return buffer.String()
}

func main() {
    app := cli.NewApp()
    app.Flags = []cli.Flag{
        cli.BoolFlag{"verbose", "be verbose"},
        cli.BoolFlag{"uniq, u", "return a uniq list"},
        cli.BoolFlag{"only-10, 0", "only ISBN10"},
        cli.BoolFlag{"only-13, 3", "only ISBN13"},
        cli.BoolFlag{"normalize, n", "normalize to ISBN13"},
    }
    app.Name = "isbngrep"
    app.Usage = "find ISBNs in texts"
    app.Version = "1.0.2"
    app.Action = func(c *cli.Context) {
        bio := bufio.NewReader(os.Stdin)
        matches := 0
        seen := make(map[string]bool)
        re := regexp.MustCompile("[0-9X][X0-9-]{9,24}")
        for {
            line, _, err := bio.ReadLine()
            if err == io.EOF {
                break
            }
            occurences := re.FindAllString(string(line), -1)
            for _, occ := range occurences {
                _, exists := seen[occ]
                if c.Bool("uniq") && exists {
                    continue
                }
                if !c.Bool("only-13") {
                    ok, value := is_valid_isbn_10(occ)
                    if ok {
                        if c.Bool("normalize") {
                            value = isbn_10_to_13(value)
                        }
                        fmt.Println(value)
                        seen[occ] = true
                        matches += 1
                        continue
                    }
                }
                if c.Bool("only-10") {
                    continue
                }
                ok, value := is_valid_isbn_13(occ)
                if ok {
                    fmt.Println(value)
                    seen[occ] = true
                    matches += 1
                }
            }
        }
        if c.Bool("verbose") {
            fmt.Fprintf(os.Stderr, fmt.Sprintf("%d\n", matches))
        }
    }
    app.Run(os.Args)
}
