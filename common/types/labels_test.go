package types

import (
	"fmt"
	"github.com/noironetworks/cilium-net/common"

	. "github.com/noironetworks/cilium-net/Godeps/_workspace/src/gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
type LabelsSuite struct{}

var _ = Suite(&LabelsSuite{})

var (
	lblsArray = []string{`%=%ed`, `//=/=`, `foo=bar`, `foo2==bar2`, `foo=====`, `foo\\==\=`, `key=`}
	lbls      = Labels{
		"foo":    NewLabel("foo", "bar", common.CiliumLabelSource),
		"foo2":   NewLabel("foo2", "=bar2", common.CiliumLabelSource),
		"key":    NewLabel("key", "", common.CiliumLabelSource),
		"foo==":  NewLabel("foo==", "==", common.CiliumLabelSource),
		`foo\\=`: NewLabel(`foo\\=`, `\=`, common.CiliumLabelSource),
		`//=/`:   NewLabel(`//=/`, "", common.CiliumLabelSource),
		`%`:      NewLabel(`%`, `%ed`, common.CiliumLabelSource),
	}
)

func (s *LabelsSuite) TestSHA256Sum(c *C) {
	str, err := lbls.SHA256Sum()
	c.Assert(err, Equals, nil)
	c.Assert(str, Equals, "f4b2082334cdcf08e58c5bb5a04bb291ecc6a4de7555baaf4d5ac110edd7d222")
}

func (s *LabelsSuite) TestSortMap(c *C) {
	sortedMap := lbls.sortMap()
	c.Assert(sortedMap, DeepEquals, lblsArray)
}

type lblTest struct {
	label  string
	result *Label
}

func (s *LabelsSuite) TestLabelShortForm(c *C) {
	lbls := []lblTest{
		{"1foo", NewLabel("1foo", "", common.CiliumLabelSource)},
		{":2foo", NewLabel("2foo", "", common.CiliumLabelSource)},
		{":3foo=", NewLabel("3foo", "", common.CiliumLabelSource)},
		{"4blah=:foo=", NewLabel("foo", "", "4blah=")},
		{"5blah::foo=", NewLabel(":foo", "", "5blah")},
		{"6foo==", NewLabel("6foo", "=", common.CiliumLabelSource)},
		{"7foo=bar", NewLabel("7foo", "bar", common.CiliumLabelSource)},
		{"k8s:foo=bar:", NewLabel("foo", "bar:", "k8s")},
	}

	for _, v := range lbls {
		res := Label{}
		decodeLabelShortForm(v.label, &res)
		c.Assert(&res, DeepEquals, v.result)
	}
}

func (s *LabelsSuite) TestMap2Labels(c *C) {
	m := Map2Labels(map[string]string{
		"foo":    "bar",
		"foo2":   "=bar2",
		"key":    "",
		"foo==":  "==",
		`foo\\=`: `\=`,
		`//=/`:   "",
		`%`:      `%ed`,
	}, common.CiliumLabelSource)
	fmt.Printf("%+v\n", m)
	c.Assert(m, DeepEquals, lbls)
}
