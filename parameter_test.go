package clif

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func mustParseTime(l, t string) time.Time {
	if r, err := time.Parse(l, t); err != nil {
		panic(err.Error())
	} else {
		return r
	}
}

type parameterResult struct {
	asString  []string
	asInt     []int
	asFloat   []float64
	asBool    []bool
	asTime    []time.Time
	asTimeErr error
	asJson    map[string]interface{}
	asJsonErr error
}

var testsParameterRead = []struct {
	from      string
	to        parameterResult
	assignErr error
}{
	{
		from: "",
		to: parameterResult{
			asString:  []string{""},
			asInt:     []int{0},
			asFloat:   []float64{0},
			asBool:    []bool{false},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"2006\""),
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: "foo",
		to: parameterResult{
			asString:  []string{"foo"},
			asInt:     []int{0},
			asFloat:   []float64{0},
			asBool:    []bool{false},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"foo\" as \"2006-01-02 15:04:05\": cannot parse \"foo\" as \"2006\""),
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: "10",
		to: parameterResult{
			asString:  []string{"10"},
			asInt:     []int{10},
			asFloat:   []float64{10},
			asBool:    []bool{true},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"10\" as \"2006-01-02 15:04:05\": cannot parse \"10\" as \"2006\""),
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: "123.234",
		to: parameterResult{
			asString:  []string{"123.234"},
			asInt:     []int{123},
			asFloat:   []float64{123.234},
			asBool:    []bool{true},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"123.234\" as \"2006-01-02 15:04:05\": cannot parse \"234\" as \"2006\""),
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: "true",
		to: parameterResult{
			asString:  []string{"true"},
			asInt:     []int{1},
			asFloat:   []float64{1},
			asBool:    []bool{true},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"true\" as \"2006-01-02 15:04:05\": cannot parse \"true\" as \"2006\""),
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: "2010-11-20 13:14:15",
		to: parameterResult{
			asString:  []string{"2010-11-20 13:14:15"},
			asInt:     []int{0},
			asFloat:   []float64{0},
			asBool:    []bool{false},
			asTime:    []time.Time{mustParseTime("2006-01-02 15:04:05", "2010-11-20 13:14:15")},
			asTimeErr: nil,
			asJson:    nil,
			asJsonErr: fmt.Errorf("unexpected end of JSON input"),
		},
		assignErr: nil,
	},
	{
		from: `{"foo":1.3,"bar":true,"baz":"bazoing"}`,
		to: parameterResult{
			asString:  []string{`{"foo":1.3,"bar":true,"baz":"bazoing"}`},
			asInt:     []int{0},
			asFloat:   []float64{0},
			asBool:    []bool{false},
			asTime:    nil,
			asTimeErr: fmt.Errorf("parsing time \"{\"foo\":1.3,\"bar\":true,\"baz\":\"bazoing\"}\" as \"2006-01-02 15:04:05\": cannot parse \"{\"foo\":1.3,\"bar\":true,\"baz\":\"bazoing\"}\" as \"2006\""),
			asJson: map[string]interface{}{
				"foo": 1.3,
				"bar": true,
				"baz": "bazoing",
			},
			asJsonErr: nil,
		},
		assignErr: nil,
	},
}

func TestParameterReadString(t *testing.T) {
	Convey("Read parameter as string", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As string: \"%s\"", i, test.from), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					So(p.Strings(), ShouldResemble, test.to.asString)
					if test.to.asString != nil {
						So(p.String(), ShouldResemble, test.to.asString[0])
					}
				}
			})
		}
	})
}

func TestParameterReadInt(t *testing.T) {
	Convey("Read parameter as int", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As int: \"%s\" -> %d", i, test.from, test.to.asInt[0]), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					So(p.Ints(), ShouldResemble, test.to.asInt)
					if test.to.asInt != nil {
						So(p.Int(), ShouldResemble, test.to.asInt[0])
					}
				}
			})
		}
	})
}

func TestParameterReadFloat(t *testing.T) {
	Convey("Read parameter as float", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As float: \"%s\" -> %g", i, test.from, test.to.asFloat[0]), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					So(p.Floats(), ShouldResemble, test.to.asFloat)
					if test.to.asFloat != nil {
						So(p.Float(), ShouldResemble, test.to.asFloat[0])
					}
				}
			})
		}
	})
}

func TestParameterReadBool(t *testing.T) {
	Convey("Read parameter as bool", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As bool: \"%s\" -> %v", i, test.from, test.to.asBool[0]), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					So(p.Bools(), ShouldResemble, test.to.asBool)
					if test.to.asBool != nil {
						So(p.Bool(), ShouldResemble, test.to.asBool[0])
					}
				}
			})
		}
	})
}

func TestParameterReadTime(t *testing.T) {
	Convey("Read parameter as time", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As time: \"%s\"", i, test.from), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					tt, err := p.Times()
					So(tt, ShouldResemble, test.to.asTime)
					if test.to.asTimeErr != nil {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldResemble, test.to.asTimeErr.Error())
					} else {
						So(err, ShouldBeNil)
					}
					if test.to.asTime != nil {
						t0, err := p.Time()
						So(t0, ShouldResemble, &test.to.asTime[0])
						So(err, ShouldBeNil)
					}
				}
			})
		}
	})
}

func TestParameterJson(t *testing.T) {
	Convey("Read parameter as json", t, func() {
		for i, test := range testsParameterRead {
			Convey(fmt.Sprintf("%d) As json: \"%s\"", i, test.from), func() {
				p := &parameter{
					Values: make([]string, 0),
				}
				err := p.Assign(test.from)
				So(err, ShouldResemble, test.assignErr)
				if test.assignErr == nil {
					jj, err := p.Json()
					So(jj, ShouldResemble, test.to.asJson)
					if test.to.asJsonErr != nil {
						So(err, ShouldNotBeNil)
					} else {
						So(err, ShouldBeNil)
					}
				}
			})
		}
	})
}
