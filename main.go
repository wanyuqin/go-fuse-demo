package main

import (
	"context"
	"flag"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func main() {
	var mountpoint string
	flag.StringVar(&mountpoint, "mountpoint", "", "mount point(dir)?")
	flag.Parse()

	if mountpoint == "" {
		log.Fatal("please input invalid mount point\n")
	}
	// 建立一个负责解析和封装 FUSE 请求监听通道对象；
	c, err := fuse.Mount(mountpoint, fuse.FSName("hellworld"), fuse.Subtype("hellofs"))
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	// 把 FS 结构体注册到 server，以便可以回调处理请求
	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}
}

// FS 文件系统主体
type FS struct {
}

func (F FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

// Dir hellofs 文件系统中，Dir是目录操作的主体
type Dir struct {
}

func (d Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 20210601
	attr.Mode = os.ModeDir | 0555

	return nil
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hello" {
		return File{}, nil
	}

	return nil, syscall.ENOENT
}

// File hellofs 文件系统中，File 结构体实现了文件系统中关于文件的调用实现
type File struct{}

const fileContent = "hello world\n"

func (f File) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 20210606
	attr.Mode = 0444
	attr.Size = uint64(len(fileContent))
	return nil
}

// ReadAll 当 cat 这个文件的时候，文件内容返回 hello，world
func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(fileContent), nil
}
