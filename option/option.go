
package option

import (
	"os"
	"container/list"
	"encoding/json"
	"flag"
)

type OptionSet struct {
	name string
	options *list.List
}

func NewOptionSet(name string) *OptionSet {
	set := &OptionSet {
		name: name,
		options: list.New(),
	}
	optionSets.PushBack(set)
	return set
}

func (self *OptionSet) String(place *string, arg string, name string, def string, desc string) {
	*place = def
	o := &stringOption {
		place: place,
		v: nil,
		arg: arg,
		name: name,
		def: def,
		desc: desc,
	}
	self.options.PushBack(o)
}

func (self *OptionSet) Bool(place *bool, arg string, name string, def bool, desc string) {
	*place = def
	o := &boolOption {
		place: place,
		v: nil,
		arg: arg,
		name: name,
		def: def,
		desc: desc,
	}
	self.options.PushBack(o)
}

func (self *OptionSet) Int(place *int, arg string, name string, def int, desc string) {
	*place = def
	o := &intOption {
		place: place,
		v: nil,
		arg: arg,
		name: name,
		def: def,
		desc: desc,
	}
	self.options.PushBack(o)
}

func Parse(configFile string) {
	initDefaults()
	parseConfig(configFile)
	registerFlags()
	flag.Parse()
	setFlags()
}

func initDefaults() {
	for set := optionSets.Front(); set != nil; set = set.Next() {
		set.Value.(*OptionSet).initDefaults()
	}
}

func registerFlags() {
	for set := optionSets.Front(); set != nil; set = set.Next() {
		set.Value.(*OptionSet).registerFlags()
	}
}

func setFlags() {
	all := list.New()
	for set := optionSets.Front(); set != nil; set = set.Next() {
		all.PushBackList(set.Value.(*OptionSet).options)
	}

	findOpt := func (arg string) optionI {
		return findOptionByArg(all, arg)
	}

	flag.Visit(func (f *flag.Flag) {
		o := findOpt(f.Name)
		if o != nil {
			o.setFlag()
		}
	})
}

type optionI interface {
	initDefault()
	registerFlag()
	getArg() string
	getName() string
	setFlag()
}

var optionSets *list.List = list.New()

type stringOption struct {
	place *string
	v *string
	arg string
	name string
	def string
	desc string
}

func (self *stringOption) initDefault() {
	*self.place = self.def
}

func (self *stringOption) registerFlag() {
	self.v = flag.String(self.arg, self.def, self.desc)
}

func (self *stringOption) getArg() string {
	return self.arg
}

func (self *stringOption) getName() string {
	return self.name
}

func (self *stringOption) setFlag() {
	*self.place = *self.v
}

type boolOption struct {
	place *bool
	v *bool
	arg string
	name string
	def bool
	desc string
}

func (self *boolOption) initDefault() {
	*self.place = self.def
}

func (self *boolOption) registerFlag() {
	self.v = flag.Bool(self.arg, self.def, self.desc)
}

func (self *boolOption) getArg() string {
	return self.arg
}

func (self *boolOption) getName() string {
	return self.name
}

func (self *boolOption) setFlag() {
	*self.place = *self.v
}

type intOption struct {
	place *int
	v *int
	arg string
	name string
	def int
	desc string
}

func (self *intOption) initDefault() {
	*self.place = self.def
}

func (self *intOption) registerFlag() {
	self.v = flag.Int(self.arg, self.def, self.desc)
}

func (self *intOption) getArg() string {
	return self.arg
}

func (self *intOption) getName() string {
	return self.name
}

func (self *intOption) setFlag() {
	*self.place = *self.v
}

func (self *OptionSet) initDefaults() {
	for o := self.options.Front(); o != nil; o = o.Next() {
		o.Value.(optionI).initDefault()
	}
}

func (self *OptionSet) registerFlags() {
	for o := self.options.Front(); o != nil; o = o.Next() {
		o.Value.(optionI).registerFlag()
	}
}

func parseConfig(file string) {
	var conf map[string]interface{}
	f, err := os.Open(file)
	if err != nil {
		return
	}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return
	}

	for k, v := range conf {
		opts, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		set := findOptionSet(optionSets, k)
		if set == nil {
			continue
		}
		set.applyConfig(opts)
	}
}

func (self *OptionSet) applyConfig(conf map[string]interface{}) {
	for k, v := range conf {
		o := findOptionByName(self.options, k)
		if o == nil {
			continue
		}

		switch opt := o.(type) {
		case *stringOption:
			str, ok := v.(string)
			if ok {
				*opt.place = str
			}
		case *boolOption:
			flg, ok := v.(bool)
			if ok {
				*opt.place = flg
			}
		case *intOption:
			flg, ok := v.(float64)
			if ok {
				*opt.place = int(flg)
			}
		}
	}
}

func findOptionSet(sets *list.List, name string) *OptionSet {
	for set := sets.Front(); set != nil; set = set.Next() {
		if set.Value.(*OptionSet).name == name {
			return set.Value.(*OptionSet)
		}
	}
	return nil;
}

func findOptionByArg(opts *list.List, arg string) optionI {
	for o := opts.Front(); o != nil; o = o.Next() {
		if o.Value.(optionI).getArg() == arg {
			return o.Value.(optionI)
		}
	}
	return nil
}

func findOptionByName(opts *list.List, name string) optionI {
	for o := opts.Front(); o != nil; o = o.Next() {
		if o.Value.(optionI).getName() == name {
			return o.Value.(optionI)
		}
	}
	return nil
}
