package simpleyaml_test

import (
	"testing"

	"github.com/smallfish/simpleyaml"
	"github.com/smallfish/simpleyaml/helper/util"
	"github.com/stretchr/testify/assert"
)

var (
	data = []byte(`
name: smallfish
age: 99
float: 3.14159
bool: true
0: IntKey
emails:
   - xxx@xx.com
   - yyy@yy.com
bb:
    cc:
        dd:
            - 111
            - 222
            - 333
        ee: aaa
`)
)

func TestBool(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("bool").Bool()
	assert.NoError(t, err)

	t.Logf("result: %v", v)
	assert.True(t, v)
}

func TestString(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("name").String()
	assert.NoError(t, err)

	t.Logf("result: %v", v)
	assert.Equal(t, v, "smallfish")
}

func TestStringFromIntKey(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get(0).String()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, v, "IntKey")
}

func TestFloat(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("float").Float()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, v, 3.14159)
}

func TestInt(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("age").Int()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, v, 99)
}

func TestGetIndex(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("bb").Get("cc").Get("dd").GetIndex(1).Int()
	assert.NoError(t, err)

	t.Logf("result: %v", v)
}

func TestString2(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("bb").Get("cc").Get("ee").String()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, v, "aaa")
}

func TestGetPath(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.GetPath("bb", "cc", "ee").String()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, v, "aaa")
}

func TestGetAllPaths(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := util.GetAllPaths(y)
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, len(v), 10)
}

func TestArray(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	v, err := y.Get("emails").Array()
	assert.NoError(t, err)

	t.Logf("result: %v", v)

	assert.Equal(t, len(v), 2)
}

func TestMap(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)
	assert.True(t, y.IsMap())

	keys, err := y.GetMapKeys()
	assert.NoError(t, err)

	t.Logf("result: %v", keys)

	assert.Equal(t, len(keys), 6)
}

func TestIsFound(t *testing.T) {
	y, err := simpleyaml.NewYaml(data)
	assert.NoError(t, err)

	assert.True(t, y.Get("name").IsFound())
	assert.False(t, y.Get("xx").IsFound())
}
