package config

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/kun98-liu/MyGodis/lib/logger"
)

// ServerProperties defines global config properties
type ServerProperties struct {
	Bind              string `cfg:"bind"`
	Port              int    `cfg:"port"`
	AppendOnly        bool   `cfg:"appendonly"`
	AppendFilename    string `cfg:"appendfilename"`
	MaxClients        int    `cfg:"maxclients"`
	RequirePass       string `cfg:"requirepass"`
	Databases         int    `cfg:"databases"`
	RDBFilename       string `cfg:"dbfilename"`
	MasterAuth        string `cfg:"masterauth"`
	SlaveAnnouncePort int    `cfg:"slave-announce-port"`
	SlaveAnnounceIP   string `cfg:"slave-announce-ip"`
	ReplTimeout       int    `cfg:"repl-timeout"`

	Peers []string `cfg:"peers"`
	Self  string   `cfg:"self"`
}

var Properties *ServerProperties

func init() {
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
	}
}

func SetupConfig(cfgFilename string) {
	file, err := os.Open(cfgFilename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	Properties = parse(file)
}

func parse(file io.Reader) *ServerProperties {

	config := &ServerProperties{}

	propMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, " ")

		if pivot > 0 && pivot < len(line)-1 {
			key := line[0:pivot]
			val := strings.Trim(line[pivot+1:], " ")
			propMap[strings.ToLower(key)] = val
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}

	// parse format
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("cfg")
		if !ok {
			key = field.Name
		}
		value, ok := propMap[strings.ToLower(key)]
		if ok {
			// fill config
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := strings.Compare("yes", value) == 0
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					slice := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(slice))
				}
			}
		}
	}
	return config

}
