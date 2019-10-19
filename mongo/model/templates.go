package model

import (
	"fmt"
	"math/rand"
	"strings"
)

type Templates struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Data        string `json:"data"`
	CreatedTime int64  `json:"createdTime"`
}

func (t *Templates) FillTemplate(min, max int, variable map[string]string) string {
	str := t.Data
	str = t.fillString(str, min, max)
	str = t.fillNumber(str, min, max)
	str = t.fillChs(str, min, max)
	str = t.fillEmpty(str, min, max)
	str = t.fillEmpty(str, min, max)
	if variable != nil {
		str = t.fillVariable(str, variable)
	}
	return str
}

func (t *Templates) fillString(str string, min, max int) string {
	n := strings.Count(str, "[RANDSTRING]")
	if n == 0 {
		return str
	}
	for i := 0; i < n; i++ {
		str = strings.Replace(str, "[RANDSTRING]", t.randString(min, max), 1)
	}
	return str
}

func (t *Templates) fillNumber(str string, min, max int) string {
	n := strings.Count(str, "[RANDNUMBER]")
	if n == 0 {
		return str
	}
	for i := 0; i < n; i++ {
		str = strings.Replace(str, "[RANDNUMBER]", t.randNumber(min, max), 1)
	}
	return str
}

func (t *Templates) fillVariable(str string, variable map[string]string) string {
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("variable%d", i+1)
		varStr := variable[key]
		n := strings.Count(str, strings.ToUpper(key))
		if n == 0 {
			continue
		}
		if varStr == "" {
			str = strings.Replace(str, strings.ToUpper(key), "", -1)
			continue
		}
		varArr := make([]string, 0)
		data := strings.Replace(varStr, "\r", "", -1)
		for _, v := range strings.Split(data, "\n") {
			varArr = append(varArr, v)
		}
		for j := 0; j < n; j++ {
			str = strings.Replace(str, strings.ToUpper(key), varArr[t.rand(0, len(varArr))], 1)
		}
	}
	return str
}

func (t *Templates) fillChs(str string, min, max int) string {
	n := strings.Count(str, "[RANDCHS]")
	if n == 0 {
		return str
	}
	for i := 0; i < n; i++ {
		str = strings.Replace(str, "[RANDCHS]", t.randChs(min, max), 1)
	}
	return str
}

func (t *Templates) fillEmpty(str string, min, max int) string {
	n := strings.Count(str, "[RANDEMPTY]")
	if n == 0 {
		return str
	}
	for i := 0; i < n; i++ {
		str = strings.Replace(str, "[RANDEMPTY]", t.randEmpty(min, max), 1)
	}
	return str
}

func (t *Templates) fillLine(str string, min, max int) string {
	n := strings.Count(str, "[RANDLINE]")
	if n == 0 {
		return str
	}
	for i := 0; i < n; i++ {
		str = strings.Replace(str, "[RANDLINE]", t.randLine(min, max), 1)
	}
	return str
}

func (t *Templates) randEmpty(min, max int) string {
	randN := t.rand(min, max)
	re := ""
	for i := 0; i < randN; i++ {
		re += "&nbsp;"
	}
	return re
}
func (t *Templates) randLine(min, max int) string {
	randN := t.rand(min, max)
	re := ""
	for i := 0; i < randN; i++ {
		re += "<br>"
	}
	return re
}

func (t *Templates) randString(min, max int) string {
	randN := t.rand(min, max)
	re := ""
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < randN; i++ {
		re += string(str[rand.Intn(len(str))])
	}
	return re
}

func (t *Templates) randNumber(min, max int) string {
	randN := t.rand(min, max)
	re := ""
	str := "0123456789"
	for i := 0; i < randN; i++ {
		re += string(str[rand.Intn(len(str))])
	}
	return re
}

func (t *Templates) randChs(min, max int) string {
	a := make([]rune, 0)
	n := t.rand(min, max)

	for i := 0; i < n; i++ {
		a = append(a, rune(t.rand(19968, 40868)))
	}
	return string(a)
}

func (t *Templates) rand(min, max int) int {
	return min + rand.Intn(max-min)
}
