# pocGoby2Xray
将Goby的json格式Poc转为xray的yaml格式Poc。

[Goby](https://gobysec.net/)和[Xray](https://github.com/chaitin/xray)是深受网络安全爱好者（包括本人）使用的社区/商业化的渗透测试工具，在[Nemo](https://github.com/hanc00l/nemo_go)项目中也集成了调用Xray进行Poc扫描。pocGoby2Xray的初衷是通过“翻译”两种工具的Poc规则和语法后进行“转换”，方便统一使用Xray的调用Poc。



## 使用

```bash
Usage of main:
  -f string
    	transform one goby poc file
  -o string
    	the xray poc output path,default is same as src file
  -p string
    	transform poc file for path
    	
➜  pocGoby2Xray git:(main) ✗ go run main.go -p gobypoc -o xraypoc
2022/11/17 21:40:51 Discuz!ML_3.x_RCE_CNVD-2019-22239.json:not support SetVariable regex
2022/11/17 21:40:51 Discuz_ML_3.x_RCE__CNVD_2019_22239.json:not support SetVariable regex
2022/11/17 21:40:51 GitLab_RCE_CVE-2021-22205.json:not support SetVariable regex
2022/11/17 21:40:51 Gitlab_RCE_CVE_2021_22205.json:not support SetVariable regex
2022/11/17 21:40:51 Tianwen_ERP_system_FileUpload_CNVD_2020_28119.json:not support SetVariable regex
2022/11/17 21:40:51 Tianwen_ERP_system__uploadfile.aspx_Arbitraryvfilevupload.json:not support SetVariable regex
2022/11/17 21:40:51 VMware_vCenter_Log4shell_CVE-2021-44228_(1).json:not support SetVariable regex
2022/11/17 21:40:51 VMware_vCenter_Log4shell_CVE_2021_44228_1.json:not support SetVariable regex
2022/11/17 21:40:51 [+]Finish: 360/368
```



## 说明

1. 由于Goby和Xray的Poc规则和语法的差异性，并不能保证Goby的Poc转换后在Xray中完全准确性，最大的差异是Goby的SetVariable与Xray的Output中，关于正则的定义和使用，如果Goby中有正则提取将取消当前Poc的转换；只保留了ScanSteps（放弃了ExploitSteps）。
2. 转换后建议使用xray poclint --script <转换后的yml文件>进行poc的lint；xray对poc的文件名和name字段有规则要求，需要转换后进行手工修改；部份rule的expression如果包含正则式需要手工修改。
3. 测试的Poc：[github.com/GREENHAT7/pxplan](https://github.com/GREENHAT7/pxplan/tree/main/goby_pocs)



## 参考

- [XrayPoc](https://docs.xray.cool/#/guide/poc/v2)
- [GobyExp](https://gobysec.net/exp)

