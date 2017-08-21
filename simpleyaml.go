// a Go package to interact with arbitrary YAML.
//
// Example:
//      var data = []byte(`
//      name: smallfish
//      age: 99
//      bool: true
//      bb:
//          cc:
//              dd:
//                  - 111
//                  - 222
//                  - 333
//      `
//
//      y, err := simpleyaml.NewYaml(data)
//      if err != nil {
//      	// ERROR
//      }
//
//      if v, err := y.Get("name").String(); err == nil {
//      	fmt.Println("value:", v)
//      }
//
//      // y.Get("age").Int()
//      // y.Get("bool").Bool()
//      // y.Get("bb").Get("cc").Get("dd").Array()
//      // y.Get("bb").Get("cc").Get("dd").GetIndex(1).Int()
//      // y.GetPath("bb", "cc", "ee").String()
package simpleyaml

import (
	"errors"
	"strconv"
	"gopkg.in/yaml.v2"
	"bytes"
)

type Yaml struct {
	data interface{}
}
var array_of_paths = make([]string, 0)
// NewYaml returns a pointer to a new `Yaml` object after unmarshaling `body` bytes
func NewYaml(body []byte) (*Yaml, error) {
	var val interface{}
	err := yaml.Unmarshal(body, &val)
	if err != nil {
		return nil, errors.New("unmarshal []byte to yaml failed: " + err.Error())
	}
	return &Yaml{val}, nil
}

// Check if the given branch was found
func (y *Yaml) IsFound() bool {
	if y.data == nil {
		return false
	}
	return true
}

// Get returns a pointer to a new `Yaml` object for `key` in its `map` representation
//
// Example:
//      y.Get("xx").Get("yy").Int()
func (y *Yaml) Get(key interface{}) *Yaml {
	m, err := y.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Yaml{val}
		}
	}
	return &Yaml{nil}
}

// GetPath searches for the item as specified by the branch
//
// Example:
//      y.GetPath("bb", "cc").Int()
func (y *Yaml) GetPath(branch ...interface{}) *Yaml {
	yin := y
	for _, p := range branch {
		yin = yin.Get(p)
	}
	return yin
}

// Array type asserts to an `array`
func (y *Yaml) Array() ([]interface{}, error) {
	if a, ok := (y.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

func (y *Yaml) IsArray() bool {
	_, err := y.Array()

	return err == nil
}

// return the size of array
func (y *Yaml) GetArraySize() (int, error) {
	a, err := y.Array()
	if err != nil {
		return 0, err
	}
	return len(a), nil
}

// GetIndex returns a pointer to a new `Yaml` object.
// for `index` in its `array` representation
//
// Example:
//      y.Get("xx").GetIndex(1).String()
func (y *Yaml) GetIndex(index int) *Yaml {
	a, err := y.Array()
	if err == nil {
		if len(a) > index {
			return &Yaml{a[index]}
		}
	}
	return &Yaml{nil}
}

// Int type asserts to `int`
func (y *Yaml) Int() (int, error) {
	if v, ok := (y.data).(int); ok {
		return v, nil
	}
	return 0, errors.New("type assertion to int failed")
}

// Bool type asserts to `bool`
func (y *Yaml) Bool() (bool, error) {
	if v, ok := (y.data).(bool); ok {
		return v, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (y *Yaml) String() (string, error) {
	if v, ok := (y.data).(string); ok {
		return v, nil
	}
	return "", errors.New("type assertion to string failed")
}

func (y *Yaml) Float() (float64, error) {
	if v, ok := (y.data).(float64); ok {
		return v, nil
	}
	return 0, errors.New("type assertion to float64 failed")
}

// Map type asserts to `map`
func (y *Yaml) Map() (map[interface{}]interface{}, error) {
	if m, ok := (y.data).(map[interface{}]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[interface]interface{} failed")
}

// Check if it is a map
func (y *Yaml) IsMap() bool {
	_, err := y.Map()
	return err == nil
}

// Get all the keys of the map
func (y *Yaml) GetMapKeys() ([]string, error) {
	m, err := y.Map()

	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for k, _ := range m {
		if s, ok := k.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}

// GetAllPaths retrieves all possible paths in the YAML file
//
// Example:
//      y.GetAllPaths(yaml_file, array_of_strings)
func GetAllExistingPaths(y *Yaml, path_slice []string) []string {
	if y.IsMap() {
		keys, err := y.GetMapKeys()
		if err != nil {
			panic(err)
		}
		for k, _ := range keys {
			if k != 0 {
				path_slice = path_slice[:len(path_slice)-1]
			}
			path_slice = append(path_slice, keys[k])
			GetAllExistingPaths(y.Get(keys[k]), path_slice)
		}
	} else if y.IsArray() {
		arr, err := y.Array()
		if err != nil {
			panic(err)
		}
		for k, _ := range arr {
			if k != 0 {
				path_slice = path_slice[:len(path_slice)-1]
			}
			path_slice = append(path_slice, strconv.Itoa(k))
			GetAllExistingPaths(y.GetIndex(k), path_slice)
		}
	} else {
		var buffer bytes.Buffer
		for k, _ := range path_slice {
				if k == len(path_slice)-1 {
				  buffer.WriteString(path_slice[k])
				}else{
				  buffer.WriteString(path_slice[k]+".")
				}
		}
		array_of_paths = append(array_of_paths, buffer.String())
	}
	return array_of_paths
}

// Helper function to invoke GetAllExistingPaths
func (y *Yaml) GetAllPaths() []string {
	yin := y
	all_paths := make([]string, 0)
	initial_path := make([]string, 0)
	all_paths = GetAllExistingPaths(yin, initial_path)
	return all_paths
}
