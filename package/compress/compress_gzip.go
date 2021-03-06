package compress

import (
	"bytes"
	"components/log"
	"compress/gzip"
	"fmt"
	"io"
)

type CompressOpGzip struct {
	direct bool
}

func (self *CompressOpGzip) Init(direct bool, params []interface{}) bool {
	self.direct = direct
	return true
}

func (self *CompressOpGzip) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
		tmpOutput, err := self.Compress(input.([]byte))
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		tmpOutput, err := self.Decompress(input.([]byte))
		if err != nil {
			fmt.Printf("unpack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil
	}

	return true, nil
}

func (self *CompressOpGzip) Compress(data []byte) ([]byte, error) {

	defer log.TraceLog("CompressOpGzip.Compress")()
	var buf bytes.Buffer
	c := gzip.NewWriter(&buf)

	_, err := c.Write(data)
	if err != nil {
		fmt.Printf("compress err: %s", err)
		return nil, err
	}

	//!!!注意：若是上面使用　defer c.Close() 会导致在下面解压时出现错误：　decompress err: unexpected EOF
	//这里使用c.Close则不会
	//但是对于decompress 中的　reader 却可以defer close
	c.Close()
	return buf.Bytes(), nil
}

func (self *CompressOpGzip) Decompress(data []byte) ([]byte, error) {

	defer log.TraceLog("CompressOpGzip.Decompress")()
	nr := bytes.NewReader(data)
	dc, err := gzip.NewReader(nr)
	if err != nil {
		fmt.Printf("decompress 1 err: %s", err)
		return nil, err
	}
	defer dc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, dc)
	if err != nil {
		fmt.Printf("decompress err: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
