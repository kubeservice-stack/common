/*
Copyright 2022 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
将io.Reader转换成io.ReadSeeker。
功能
 1. 超过一定大小、将临时缓冲转成临时文件
 2. 异步读取源信息
 3. 异步读取源信息后关闭

Attention

	使用异步读取数据时候, 在调用其他方法之前一定要先使用Await。
	该操作等待全部数据的加载, 同时可以处理读取数据时的错误。
*/
package temporary

import (
	"io"
	"sync"
)

type Temporary interface {
	Await() error
	Size() int64
	Type() string
	Name() string
	Bytes() []byte
	io.Reader
	io.Seeker
	io.Closer
}

const (
	TemporaryBuffer string = "Buffer"
	TemporaryFile   string = "File"
)

type temporaryItemer interface {
	Size() int64
	Type() string
	Name() string
	Bytes() []byte
	Sync() error
	io.Writer
	io.Reader
	io.Seeker
	io.Closer
}

type temporary struct {
	item         temporaryItemer
	itemConvert  bool
	itemWg       *sync.WaitGroup
	itemAsyncErr error

	maxBufferSize int64
	fileDir       string
	filePattern   string
}

func NewTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary, err error) {
	var temp = &temporary{
		item:        newTemporaryBuffer(),
		itemConvert: false,
		itemWg:      new(sync.WaitGroup),

		maxBufferSize: maxBufferSize,
		fileDir:       fileDir,
		filePattern:   filePattern,
	}
	if _, err = io.Copy(temp, reader); err != nil {
		return
	}
	if err = temp.item.Sync(); err != nil {
		return
	}
	return temp, nil
}

func NewAsyncTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary) {
	var temp = &temporary{
		item:        newTemporaryBuffer(),
		itemConvert: false,
		itemWg:      new(sync.WaitGroup),

		maxBufferSize: maxBufferSize,
		fileDir:       fileDir,
		filePattern:   filePattern,
	}
	temp.itemWg.Add(1)
	go func(temp *temporary, reader io.Reader) {
		if _, temp.itemAsyncErr = io.Copy(temp, reader); temp.itemAsyncErr != nil {
			goto end
		}
		if temp.itemAsyncErr = temp.item.Sync(); temp.itemAsyncErr != nil {
			goto end
		}
	end:
		temp.itemWg.Done()
	}(temp, reader)
	return temp
}

func NewMustCloseReaderAsyncTemporary(readcloser io.ReadCloser, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary) {
	var temp = &temporary{
		item:        newTemporaryBuffer(),
		itemConvert: false,
		itemWg:      new(sync.WaitGroup),

		maxBufferSize: maxBufferSize,
		fileDir:       fileDir,
		filePattern:   filePattern,
	}
	temp.itemWg.Add(1)
	go func(temp *temporary, readcloser io.ReadCloser) {
		if _, temp.itemAsyncErr = io.Copy(temp, readcloser); temp.itemAsyncErr != nil {
			goto end
		}
		if temp.itemAsyncErr = temp.item.Sync(); temp.itemAsyncErr != nil {
			goto end
		}
	end:
		io.Copy(io.Discard, readcloser)
		readcloser.Close()
		temp.itemWg.Done()
	}(temp, readcloser)
	return temp
}

func (temp *temporary) toTemporaryFile() (err error) {
	if temp.itemConvert {
		return
	}
	var (
		seekPosition int64
		tf           *temporaryFile
	)
	if seekPosition, err = temp.item.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	if tf, err = newTemporaryFile(temp.fileDir, temp.filePattern); err != nil {
		return
	}
	if _, err = temp.item.Seek(0, io.SeekStart); err != nil {
		return
	}
	if _, err = io.Copy(tf, temp.item); err != nil {
		return
	}
	temp.item.Close()
	if _, err = tf.Seek(seekPosition, io.SeekStart); err != nil {
		return
	}
	temp.item = tf
	temp.itemConvert = true
	return
}

func (temp *temporary) Write(p []byte) (n int, err error) {
	if temp.item.Size()+int64(len(p)) > temp.maxBufferSize {
		if err = temp.toTemporaryFile(); err != nil {
			return
		}
	}
	if n, err = temp.item.Write(p); err != nil {
		if err == ErrBufferTooLarge { // 内存分配出错
			if err = temp.toTemporaryFile(); err != nil {
				return
			}
			var m int
			m, err = temp.item.Write(p[n:])
			n += m
			if err != nil {
				return
			}
			return
		}
		return
	}
	return
}

func (temp *temporary) Await() error {
	temp.itemWg.Wait()
	return temp.itemAsyncErr
}

func (temp *temporary) Size() int64 {
	return temp.item.Size()
}

func (temp *temporary) Type() string {
	return temp.item.Type()
}

func (temp *temporary) Name() string {
	return temp.item.Name()
}

func (temp *temporary) Bytes() []byte {
	return temp.item.Bytes()
}

func (temp *temporary) Read(p []byte) (n int, err error) {
	return temp.item.Read(p)
}

func (temp *temporary) Seek(offset int64, whence int) (abs int64, err error) {
	return temp.item.Seek(offset, whence)
}

func (temp *temporary) Close() (err error) {
	return temp.item.Close()
}
