package shuai_utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func EncodeGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	// 确保所有数据都被写入
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CalcFuncCostTime(logKey string, f func()) {
	timeBegin := time.Now()
	defer func() {
		log.Printf("calcFuncCostTime processed, logKey: %s, cost time: %s", logKey, time.Since(timeBegin).String())
	}()
	f()
}
