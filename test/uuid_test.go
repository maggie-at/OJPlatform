package test

import (
	"fmt"
	"testing"
)
import "github.com/satori/go.uuid"

func TestGenerateUUID(t *testing.T) {
	// 生成一个长度为36的UUID
	UUID := uuid.NewV4().String()
	fmt.Printf("%T(%d) %v", UUID, len(UUID), UUID)
}
