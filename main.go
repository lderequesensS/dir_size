package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
)

func main() {
	working_directory, err := os.Getwd()
	if err != nil {
		panic(1)
	}

	directory := flag.String("d", working_directory, "Use another path")
	recursive := flag.Bool("r", true, "Calculate size recursively") // To change this -r=false
	use_decimal := flag.Bool("u", false, "Use decimal units instead of binary")
	ignore_unit := flag.Bool("b", false, "Print Bytes, ignores u flag")
	flag.Parse()

	file_system := os.DirFS(*directory)
	var total_bytes int64 = 0

	fs.WalkDir(file_system, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		stat, err := fs.Stat(file_system, path)
		// Here if breaking with some links, for example in tor-browser
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("Found weid link: %s\n", path)
				return nil
			}
			panic(err)
		}

		total_bytes += stat.Size()

		if !*recursive && stat.IsDir() && stat.Name() != "." {
			return fs.SkipDir
		}
		return nil
	})

	extra := ""
	if !*recursive {
		extra = "(Not recursive result)"
	}

	if *ignore_unit {
		fmt.Printf("Size of directory %s: %v %s\n", *directory, total_bytes, extra)
	} else {
		fmt.Printf("Size of directory %s: %v %s\n", *directory, human_numbers(total_bytes, *use_decimal), extra)
	}
}

func human_numbers(total_bytes int64, use_decimal bool) string {
	var divisor int64 = 1024
	binary_unit := map[int]string{
		0: "B",
		1: "KiB",
		2: "MiB",
		3: "GiB",
		4: "TiB",
	}
	decimal_unit := map[int]string{
		0: "B",
		1: "KB",
		2: "MB",
		3: "GB",
		4: "TB",
	}

	if use_decimal {
		divisor = 1000
	}

	if total_bytes < divisor {
		return fmt.Sprintf("%vB", total_bytes)
	}

	order := len(fmt.Sprint(total_bytes)) / 3
	number := float64(total_bytes) / (math.Pow(float64(divisor), float64(order)))

	unit := binary_unit[order]
	if use_decimal {
		unit = decimal_unit[order]
	}

	return fmt.Sprintf("%.2v%s", number, unit)
}
