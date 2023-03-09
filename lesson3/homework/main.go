package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    int64
	Limit     int64
	BlockSize int
	Conv      []string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "number of bytes to skip before copying")
	flag.Int64Var(&opts.Limit, "limit", -1, "maximum number of bytes to copy. 0 - copy all data")
	flag.IntVar(&opts.BlockSize, "block-size", 1024, "size of one block in bytes")
	convStr := flag.String("conv", "", "comma-separated list of conversions to apply to the data. possible values: upper_case, lower_case, trim_spaces")
	flag.Parse()

	if *convStr != "" {
		convList := strings.Split(*convStr, ",")
		for _, conv := range convList {
			conv = strings.TrimSpace(conv)
			switch conv {
			case "upper_case", "lower_case", "trim_spaces":
				opts.Conv = append(opts.Conv, conv)
			default:
				return nil, fmt.Errorf("unsupported conversion: %s", conv)
			}
		}
	}

	if contains(opts.Conv, "upper_case") && contains(opts.Conv, "lower_case") {
		_, _ = fmt.Fprintln(os.Stderr, "can not edit to upper_case and lower_case at the same time")
		os.Exit(1)
	}

	return &opts, nil
}

type StdinReader struct{}

func (r *StdinReader) Seek(offset int64, whence int) (n int64, err error) {
	return io.CopyN(io.Discard, os.Stdin, offset)
}

func (r *StdinReader) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (r *StdinReader) Close() error {
	return nil
}

type StdoutWriter struct{}

func (w *StdoutWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (w *StdoutWriter) Close() error {
	return nil
}

func contains(container []string, certain string) bool {
	for _, conv := range container {
		if conv == certain {
			return true
		}
	}
	return false
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	var fromReader io.ReadSeekCloser
	var toWriter io.WriteCloser

	if opts.From == "" {
		fromReader = &StdinReader{}
	} else {
		fromReader, err = os.Open(opts.From)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not open input file:", err)
			os.Exit(1)
		}

		defer func(fromReader io.ReadSeekCloser) {
			err := fromReader.Close()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "can not close file:", opts.From, err)
				os.Exit(1)
			}
		}(fromReader)
	}

	if opts.To == "" {
		toWriter = &StdoutWriter{}
	} else {
		toWriter, err = os.Open(opts.To)
		if err != nil {
			toWriter, err = os.Create(opts.To)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "can not create output file: ", err)
				os.Exit(1)
			}
		}
		defer func(toWriter io.WriteCloser) {
			err := toWriter.Close()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "can not close file: ", opts.To, err)
				os.Exit(1)
			}
		}(toWriter)
	}

	SeekFileFrom(&fromReader, opts.Offset)
	TransferInfo(opts, &fromReader, &toWriter)
}

func TransferInfo(opts *Options, fromReader *io.ReadSeekCloser, toWriter *io.WriteCloser) {
	buf := make([]byte, opts.BlockSize)
	var tmp []byte
	total := int64(0)
	for {
		total += int64(opts.BlockSize)
		if opts.Limit > 0 && total-opts.Limit > 0 {
			buf = make([]byte, int64(opts.BlockSize)-(total-opts.Limit))
		}
		n, errRead := (*fromReader).Read(buf)
		tmp = buf[:n]
		if errRead != nil && errRead != io.EOF {
			_, _ = fmt.Fprintln(os.Stderr, "error reading input file:", errRead)
			os.Exit(1)
		}
		if opts.Conv != nil {
			for _, conv := range opts.Conv {
				switch conv {
				case "upper_case":
					tmp = []byte(strings.ToUpper(string(tmp)))
				case "lower_case":
					tmp = []byte(strings.ToLower(string(tmp)))
				case "trim_spaces":
					tmp = []byte(strings.TrimSpace(string(tmp)))
				default:
					_, _ = fmt.Fprintln(os.Stderr, "invalid conversion option:", conv)
					os.Exit(1)
				}
			}
		}

		if n == 0 {
			break
		}
		_, err := (*toWriter).Write(tmp)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error writing output file:", err)
			os.Exit(1)
		}
		if opts.Limit > 0 && total-opts.Limit > 0 {
			break
		}
	}
}

func SeekFileFrom(fromReader *io.ReadSeekCloser, position int64) {
	if position > 0 {
		_, err := (*fromReader).Seek(position, io.SeekStart)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not seek to offset:", err)
			os.Exit(1)
		}
	} else if position < 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Offset must be non-negative")
		os.Exit(1)
	}
}
