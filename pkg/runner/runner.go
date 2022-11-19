package runner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"pocGoby2Xray/pkg/gobypoc"
	"pocGoby2Xray/pkg/xraypocv2"
	"strings"
)

type PocRunner struct {
}

func NewPocRunner() *PocRunner {
	r := PocRunner{}
	return &r
}

func (p *PocRunner) LoadGobyJsonPocFile(fileName string) (poc gobypoc.PocJson, err error) {
	var fileContent []byte
	fileContent, err = os.ReadFile(fileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileContent, &poc)
	return
}

func (p *PocRunner) SaveXrayPocFile(filename string, xrayPoc xraypocv2.Poc) (err error) {
	var fileContent []byte
	fileContent, err = yaml.Marshal(xrayPoc)
	if err != nil {
		return
	}
	err = os.WriteFile(filename, p.transformVars(fileContent), 0666)
	return
}

func (p *PocRunner) transformVars(fileContent []byte) (results []byte) {
	results = bytes.ReplaceAll(fileContent, []byte("{{{"), []byte("{{"))
	results = bytes.ReplaceAll(results, []byte("}}}"), []byte("}}"))

	return
}

func (p *PocRunner) Run(gobyPoc gobypoc.PocJson) (xrayPoc xraypocv2.Poc, err error) {
	if len(gobyPoc.ScanSteps) < 2 {
		err = errors.New("no ScanSteps")
		return
	}
	// 基础信息部份
	xrayPoc.Name = gobyPoc.Name
	xrayPoc.Manual = true
	xrayPoc.Query = gobyPoc.Query
	xrayPoc.Transport = "http"
	xrayPoc.Detail.Author = gobyPoc.Author
	xrayPoc.Detail.Links = gobyPoc.References
	xrayPoc.Detail.Tags = strings.Join(gobyPoc.Tags, ",")
	xrayPoc.Detail.Description = gobyPoc.Description
	// 规则 Rule
	var ruleKey []string
	xrayPoc.Rules = make([]yaml.MapItem, 0)
	ScanStepOpts := gobyPoc.ScanSteps[1:]
	for index, ScanStepOpt := range ScanStepOpts {
		// 解析每一个rule
		var ruleItem xraypocv2.Rule
		var sets []xraypocv2.VariableMapItem
		ruleItem, sets, err = p.analyseScanStep(ScanStepOpt.(map[string]interface{}))
		if err != nil {
			return
		}
		xrayPoc.Rules = append(xrayPoc.Rules, yaml.MapItem{
			Key:   fmt.Sprintf("r%d", index),
			Value: ruleItem,
		})
		// request中的set_variable，要作为xray的全局变量
		for _, set := range sets {
			xrayPoc.Set = append(xrayPoc.Set, yaml.MapItem{
				Key:   set.Key,
				Value: set.Value,
			})
		}
		ruleKey = append(ruleKey, fmt.Sprintf("r%d()", index))
	}
	// expression
	expressionCondition := " && "
	if gobyPoc.ScanSteps[0] == "OR" {
		expressionCondition = " || "
	}
	xrayPoc.Expression = strings.Join(ruleKey, expressionCondition)

	return
}

func (p *PocRunner) analyseScanStep(ScanStepOpt map[string]interface{}) (rule xraypocv2.Rule, sets []xraypocv2.VariableMapItem, err error) {
	//在Rule中的SetVariable
	if _, ok := ScanStepOpt["SetVariable"]; ok {
		//通过Marshal/Unmarshal，将map[]interface{}转换为structure
		var setVbytes []byte
		if setVbytes, err = json.Marshal(ScanStepOpt["SetVariable"]); err != nil {
			return
		}
		var setVst []string
		if err = json.Unmarshal(setVbytes, &setVst); err != nil {
			return
		}
		if len(setVst) > 0 {
			var outputs []xraypocv2.VariableMapItem
			if outputs, err = p.analyseSetVariable(setVst); err != nil {
				return
			}
			for _, output := range outputs {
				rule.Output = append(rule.Output, yaml.MapItem{
					Key:   output.Key,
					Value: output.Value,
				})
			}
		}
	}
	var scanStep gobypoc.ScanStep
	// Request部份
	if _, ok := ScanStepOpt["Request"]; !ok {
		err = errors.New("no Request")
		return
	}
	//通过Marshal/Unmarshal，将map[]interface{}转换为structure
	var request, response []byte
	if request, err = json.Marshal(ScanStepOpt["Request"]); err != nil {
		return
	}
	if err = json.Unmarshal(request, &scanStep.Requests); err != nil {
		return
	}
	// 在Request中的set_variable，需设置为全局变量
	if scanStep.Requests.SetVariable != nil && len(scanStep.Requests.SetVariable) > 0 {
		if sets, err = p.analyseRequestSetVariable(scanStep.Requests.SetVariable); err != nil {
			return
		}
	}
	rule.Request.Path = scanStep.Requests.Uri
	rule.Request.FollowRedirects = scanStep.Requests.FollowRedirect
	rule.Request.Body = scanStep.Requests.Data
	rule.Request.Headers = scanStep.Requests.Header
	rule.Request.Method = scanStep.Requests.Method
	// ResponseTest部份
	if _, ok := ScanStepOpt["ResponseTest"]; !ok {
		err = errors.New("no Response")
		return
	}
	//通过Marshal/Unmarshal，将map[]interface{}转换为structure
	if response, err = json.Marshal(ScanStepOpt["ResponseTest"]); err != nil {
		return
	}
	if err = json.Unmarshal(response, &scanStep.ResponseTest); err != nil {
		return
	}
	if len(scanStep.ResponseTest.Checks) > 0 {
		if scanStep.ResponseTest.Type == "group" {
			checks := p.analyseCheck(scanStep.ResponseTest.Checks)
			if scanStep.ResponseTest.Operation == "AND" {
				rule.Expression = strings.Join(checks, " && ")
			} else {
				rule.Expression = strings.Join(checks, " || ")
			}
		}
	} else {
		//err = errors.New("no response check")
		// 有些poc的步骤没有check,默认增加一个true
		// expression: true
		// rule.Expression = strings.Join([]string{"true"}, "")
		// 由于expression要返回bool，string的true为加""，直接加expression: true不行
		// 所以还是用response.status >= 200代表
		rule.Expression = strings.Join([]string{"response.status >= 200"}, "")
	}
	return
}

func (p *PocRunner) analyseCheck(checks []gobypoc.Check) (checkResult []string) {
	for _, check := range checks {
		if check.Type == "item" {
			switch check.Variable {
			case "$code":
				if check.Operation == "==" || check.Operation == ">" || check.Operation == ">=" || check.Operation == "<" || check.Operation == "<=" || check.Operation == "!=" {
					r := fmt.Sprintf("response.status %s %s", check.Operation, check.Value)
					checkResult = append(checkResult, r)
				} else {
					r := fmt.Sprintf("response.status %s %s", "==", check.Value)
					checkResult = append(checkResult, r)
				}
			case "$body":
				switch check.Operation {
				case "contains":
					r := fmt.Sprintf("response.body.bcontains(b\"%s\")", convertStr(check.Value))
					checkResult = append(checkResult, r)
				case "regex":
					r := fmt.Sprintf("\"%s\".bmatches(response.body)", convertStr(check.Value))
					checkResult = append(checkResult, r)
				case "not contains":
					r := fmt.Sprintf("!response.body.bcontains(b\"%s\")", convertStr(check.Value))
					checkResult = append(checkResult, r)
				}
			case "$head":
				switch check.Operation {
				case "contains":
					r := fmt.Sprintf("response.raw_header.bcontains(b\"%s\")", convertStr(check.Value))
					checkResult = append(checkResult, r)
				case "regex":
					r := fmt.Sprintf("\"%s\".bmatches(response.raw_header)", convertStr(check.Value))
					checkResult = append(checkResult, r)
				case "not contains":
					r := fmt.Sprintf("!response.raw_header.bcontains(b\"%s\")", convertStr(check.Value))
					checkResult = append(checkResult, r)
				}
			}
		}
	}
	return
}

func (p *PocRunner) analyseSetVariable(varLine []string) (outputItem []xraypocv2.VariableMapItem, err error) {
	//由于xray poc的正则规则匹配语法与goby poc差异较大，只能简单解析不带正则，或正则为空的SetVariable
	//output|lastbody
	//output|lastbody|
	//output|lastbody||
	//output|lastbody|regex|
	//output|lastbody|regex|(.*)
	for _, v := range varLine {
		if len(v) == 0 {
			continue
		}
		vars := strings.Split(v, "|")
		if len(vars) < 2 {
			continue
		}
		//变量名为output一般无什么用，直接忽略
		if vars[0] == "output" {
			continue
		}
		if len(vars) == 2 ||
			(len(vars) == 3 && len(vars[2]) == 0) ||
			(len(vars) == 4 && len(vars[2]) == 0 && len(vars[3]) == 0) ||
			(len(vars) == 4 && vars[3] == "regex" && (len(vars[3]) == 0 || vars[3] == "(.*)")) {
			item := xraypocv2.VariableMapItem{Key: vars[0]}
			switch vars[1] {
			case "lastbody":
				item.Value = "response.body"
			case "lastheader":
				item.Value = "response.headers"
			case "statusline":
				err = errors.New("not support SetVariable statusline")
				return
			default:
				err = errors.New("not support SetVariable type")
				return
			}
			outputItem = append(outputItem, item)
		} else {
			//不支持正则：
			err = errors.New("not support SetVariable regex")
			return
		}
	}
	return
}

func (p *PocRunner) analyseRequestSetVariable(varLine []string) (outputItem []xraypocv2.VariableMapItem, err error) {
	//请求时的参数设置，目前只实现随机数生成
	//r1|rand|int|8
	//r2|rand|str|32
	for _, v := range varLine {
		if len(v) == 0 {
			continue
		}
		vars := strings.Split(v, "|")
		if len(vars) < 4 {
			continue
		}
		if vars[1] == "rand" {
			item := xraypocv2.VariableMapItem{Key: vars[0]}
			if vars[2] == "int" {
				item.Value = "randomInt(400000, 448000)"
			} else if vars[2] == "str" {
				item.Value = fmt.Sprintf("randomLowercase(%s)", vars[3])
			} else {
				err = errors.New("not support SetVariable type")
				return
			}
			outputItem = append(outputItem, item)
		}
	}
	return
}

func convertStr(input string) (output string) {
	output = strings.ReplaceAll(input, "\"", "\\\"")
	return
}
