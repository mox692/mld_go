package main

import (
	"fmt"
	"os"

	"github.com/mox692/mld"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid Arg counts %d", len(os.Args))
		return
	}
	os.Exit(mld.Run())
}

/*
#include <stdint.h>
#include <stdio.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

typedef struct EFLInfo EFLInfo;
typedef struct TempRegion TempRegion;
struct EFLInfo {
    uint32_t sectionHeader; // セクションヘッダtableのaddr
};

// ELF headerを読んで、必要な情報をELFInfoに書き出す
EFLInfo readELFHeader() {

}

// それぞれのsectionを読んで,そのままtmpSectionsに書き写す。
// ただし、特定のセクション(.text, .bss, sym系)だったらそれぞれの処理に分岐.
int readSections(uint32_t sectionHeaderAddr, TempRegion *tmpSections) {}

// programヘッダを作る。readSectionsとは別のTempRegionを作成して、そこに書く.
// sectionの情報がいくらか必要になると思うので、それも引数で渡してる.
int createProgramHeader(TempRegion *tmpProgramHeader, int sectioninfo) {}

// 最後にelfに書いていくものを想定。tmpsectionsとtmpprogramheaderとelfheaderを書いていく想定。
// 書いたbyte数を返す.
int write(int fd, uint32_t offset, int size) {};

int main(int argc, char **argv) {
    if(argc != 2) {
        printf("invalid argc: %d\n", argc);
        return -1;
    }

    int fd = open(argv[1], O_RDWR);
    if (fd == -1) {
        printf("Err: %s\n", strerror(errno));
        return -1;
    }

    struct stat sb;
    if (fstat(fd, &sb) == -1) {
        printf("Err: %s\n", strerror(errno));
        return -1;
    }

    void *head = mmap(NULL, sb.st_size, PROT_WRITE | PROT_READ | PROT_EXEC, MAP_SHARED, fd, 0);
    printf("ad: %p\n", head);

    if(head == MAP_FAILED) {
        printf("mapfile err.\n");
        return -1;
    }

    return 0;
}


*/
