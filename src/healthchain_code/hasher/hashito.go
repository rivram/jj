package hasher

import (
  "crypto/md5"
)

func Hash(data []byte) string {
  hashed :=  md5.Sum(data)
  return string(hashed[:])
}
