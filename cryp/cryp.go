package cryp
import (
 "crypto/aes"
 "os"
)
func AESEncrypt(src []byte, key []byte) (encrypted []byte) {
 cipher, _ := aes.NewCipher(generateKey(key))
 length := (len(src) + aes.BlockSize) / aes.BlockSize
 plain := make([]byte, length*aes.BlockSize)
 copy(plain, src)
 pad := byte(len(plain) - len(src))
 for i := len(src); i < len(plain); i++ {
 plain[i] = pad
 }
 encrypted = make([]byte, len(plain))
 // 分组分块加密 
 for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
 cipher.Encrypt(encrypted[bs:be], plain[bs:be])
 }
 return encrypted
}
func AESDecrypt(encrypted []byte, key []byte) (decrypted []byte) {
 cipher, _ := aes.NewCipher(generateKey(key))
 decrypted = make([]byte, len(encrypted))
 // 
 for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
 cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
 }
  trim  := 0
 if len(decrypted) > 0 {
 trim = len(decrypted) - int(decrypted[len(decrypted)-1])
 }
 return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
 genKey = make([]byte, 16)
 copy(genKey, key)
 for i := 16; i < len(key); {
 for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
 genKey[j] ^= key[i]
 }
 }
 return genKey
}
// 读取字节
func ReadFile(fileName string) []byte {
	file, err := os.Open(fileName)
	CheckErr(err)
	defer file.Close()

	fileInfo, err := file.Stat()
	CheckErr(err)

	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	CheckErr(err)

	return buf
}

// 字节存储成文件
func WriteFile(fileName string, data []byte) {
    file, err := os.Create(fileName)
    CheckErr(err)
    defer file.Close()

    _, err = file.Write(data)
    CheckErr(err)
}

// 文件名 key文件名
func Encrpyt(file []byte, keyPath string,fileName string)  {
    key:=ReadFile(keyPath)
    cryp:=AESEncrypt(file,key)
    WriteFile(fileName,cryp)
}
// 解密成字节
func Decrypt(file []byte,keyName string)[]byte {
    key:=ReadFile(keyName)
    rf:=AESDecrypt(file,key)
    return rf
}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}