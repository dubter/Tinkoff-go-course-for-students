package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
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

func (r *StdinReader) Seek(offset int64, _ int) (n int64, err error) {
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

	fromReader, err := OpenFileFrom(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	toWriter, err := OpenFileTo(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = SeekFileFrom(fromReader, opts.Offset)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = TransferInfo(opts, fromReader, toWriter)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func TransferInfo(opts *Options, fromReader io.ReadSeekCloser, toWriter io.WriteCloser) (int64, error) {
	buf := make([]byte, opts.BlockSize)
	var totalRead []byte
	var n int
	var errRead error
	total := int64(0)
	if opts.BlockSize == 1 {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, fromReader)
		if err != nil {
			return 0, fmt.Errorf("error reading input file: %v", err)
		}
		totalRead = buf.Bytes()

		if opts.Conv != nil {
			for _, conv := range opts.Conv {
				switch conv {
				case "upper_case":
					totalRead = []byte(strings.ToUpper(string(totalRead)))
				case "lower_case":
					totalRead = []byte(strings.ToLower(string(totalRead)))
				case "trim_spaces":
					if len(string(totalRead)) > 1 {
						totalRead = []byte(strings.TrimSpace(string(totalRead)))
					}
				default:
					return 0, fmt.Errorf("invalid conversion option: %v", conv)
				}
			}
		}

		_, err = toWriter.Write(totalRead)
		if err != nil {
			return 0, fmt.Errorf("error writing output file: %v", err)
		}
	} else {
		for {
			totalRead = nil
			for {
				total += int64(opts.BlockSize)
				if opts.Limit > 0 && total-opts.Limit > 0 {
					buf = make([]byte, int64(opts.BlockSize)-(total-opts.Limit))
				}
				n, errRead = fromReader.Read(buf)
				totalRead = append(totalRead, buf[:n]...)

				if errRead != nil && errRead != io.EOF {
					return 0, fmt.Errorf("error reading input file: %v", errRead)
				}

				if opts.Limit > 0 && total-opts.Limit > 0 {
					break
				}
				if n == 0 {
					break
				}
				if utf8.Valid(totalRead) {
					break
				}
			}

			if opts.Conv != nil {
				for _, conv := range opts.Conv {
					switch conv {
					case "upper_case":
						totalRead = []byte(strings.ToUpper(string(totalRead)))
					case "lower_case":
						totalRead = []byte(strings.ToLower(string(totalRead)))
					case "trim_spaces":
						totalRead = []byte(strings.TrimSpace(string(totalRead)))
					default:
						return 0, fmt.Errorf("invalid conversion option: %v", conv)
					}
				}
			}

			_, err := toWriter.Write(totalRead)
			if n == 0 {
				break
			}
			if err != nil {
				return 0, fmt.Errorf("error writing output file: %v", err)
			}
			if opts.Limit > 0 && total-opts.Limit > 0 {
				break
			}
		}
	}
	return total, nil
}

func SeekFileFrom(fromReader io.ReadSeekCloser, position int64) (int64, error) {
	var n int64
	var err error

	if position > 0 {
		n, err = fromReader.Seek(position, io.SeekStart)

		if err != nil {
			return n, fmt.Errorf("can not seek to offset: %v", err)
		}
	} else if position < 0 {
		return n, fmt.Errorf("offset must be non-negative")
	}
	return n, nil
}

func OpenFileFrom(opts *Options) (io.ReadSeekCloser, error) {
	var fromReader io.ReadSeekCloser
	var err error

	if opts.From == "" {
		fromReader = &StdinReader{}
	} else {
		fromReader, err = os.Open(opts.From)

		if err != nil {
			return fromReader, fmt.Errorf("can not open input file: %v", err)
		}
	}
	return fromReader, nil
}

func OpenFileTo(opts *Options) (io.WriteCloser, error) {
	var toWriter io.WriteCloser
	var err error

	if opts.To == "" {
		toWriter = &StdoutWriter{}
	} else {
		toWriter, err = os.Open(opts.To)
		if err != nil {
			toWriter, err = os.Create(opts.To)
			if err != nil {
				return toWriter, fmt.Errorf("can not create output file: %v", err)
			}
		}
	}
	return toWriter, nil
}
