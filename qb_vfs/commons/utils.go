package commons

import (
	"strings"

	qbc "github.com/rskvp/qb-core"
)

func ErrorContains(err error, text string) bool {
	lowErr := strings.ToLower(err.Error())
	lowTxt := strings.ToLower(text)
	return strings.Contains(lowErr, lowTxt)
}

func IsFile(path string) bool{
	return len(qbc.Paths.Extension(path))>0
}

func Absolutize(root, path string) string {
	if strings.HasPrefix(path, ".") {
		return qbc.Paths.Concat(root, path)
	}
	return path
}

func Relativize(root, path string) string {
	if !strings.HasPrefix(path, ".") {
		response :=  strings.Replace(path, root, "", 1)
		if strings.HasPrefix(response, "/"){
			response = "." + response
		} else {
			response = "./" + response
		}
		return response
	}
	return path
}

func ReadKey(pathOrKey string) ([]byte, error) {
	if len(pathOrKey) == 0 {
		return []byte{}, nil
	}
	if b, err := qbc.Paths.IsFile(pathOrKey); b && nil == err {
		data, err := qbc.IO.ReadBytesFromFile(pathOrKey)
		if nil != err {
			return []byte{}, err
		}
		return data, nil
	}
	return []byte(pathOrKey), nil
}

func SplitHost(settings *VfsSettings) (host string, port int) {
	port = 22
	_, full := settings.SplitLocation()
	tokens := strings.Split(full, ":")
	switch len(tokens) {
	case 1:
		host = tokens[0]
	case 2:
		host = tokens[0]
		port = qbc.Convert.ToInt(tokens[1])
	}
	return
}
