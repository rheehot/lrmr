package playground

import (
	"bufio"
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/therne/lrmr/lrdd"
	"github.com/therne/lrmr/output"
	"github.com/therne/lrmr/transformation"
	"io"
	"os"
)

type ndjsonDecoder struct {
	transformation.Simple
}

func NewLDJSONDecoder() transformation.Transformation {
	return &ndjsonDecoder{}
}

func (l *ndjsonDecoder) DescribeOutput() *transformation.OutputDesc {
	return transformation.DescribingOutput().WithRoundRobin()
}

func (l *ndjsonDecoder) Run(row lrdd.Row, out output.Output) error {
	path := row["path"].(string)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file : %w", err)
	}

	r := bufio.NewReader(file)
	for {
		line, err := readline(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		msg := make(lrdd.Row)
		if err := jsoniter.Unmarshal(line, &msg); err != nil {
			return err
		}
		if err := out.Send(msg); err != nil {
			return err
		}
	}
	return file.Close()
}

func readline(r *bufio.Reader) (line []byte, err error) {
	var isPrefix = true
	var ln []byte
	var buf bytes.Buffer
	for isPrefix && err == nil {
		ln, isPrefix, err = r.ReadLine()
		buf.Write(ln)
	}
	line = buf.Bytes()
	return
}
