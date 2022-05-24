package tag

import (
	"strings"
)

func addTag(beforeTag string, tag string) string {
	if beforeTag == "" {
		return tag
	}
	return beforeTag + ", " + tag
}

func removeTag(beforeTag string, tag string) string {
	if beforeTag == tag {
		return ""
	} else {
		beforeTag = strings.TrimSpace(strings.ReplaceAll(beforeTag, ", "+tag, ""))
		beforeTag = strings.TrimSpace(strings.ReplaceAll(beforeTag, ","+tag, ""))
		beforeTag = strings.TrimSpace(strings.ReplaceAll(beforeTag, tag, ""))
		return beforeTag

	}
}

//// 为poc添加标签
//func AddTags(tags []string, xrayPocMap map[string]xray_structs.Poc, nucleiPocMap map[string]nuclei_structs.Poc) {
//	for _, tag := range tags {
//		for pocPath, poc := range xrayPocMap {
//			pocTags := poc.Detail.Tags
//			if !strings.Contains(pocTags, tag) {
//				poc.Detail.Tags = addTag(pocTags, tag)
//				out, err := yaml.Marshal(poc)
//				if err != nil {
//					utils.CliError("Can't Marshal poc: "+poc.Name, 2)
//				}
//				err = utils.WriteFile(pocPath, out)
//				if err != nil {
//					utils.CliError("Can't write file: "+pocPath, 3)
//				}
//				utils.SuccessF("%s: Add [%s] Tag", poc.Name, tag)
//			}
//		}
//		for pocPath, poc := range nucleiPocMap {
//			pocTags := poc.Info.Tags.String()
//			if !strings.Contains(pocTags, tag) {
//				poc.Info.Tags = stringslice.StringSlice{
//					Value: addTag(pocTags, tag),
//				}
//				out, err := yaml.Marshal(poc)
//				if err != nil {
//					utils.CliError("Can't Marshal poc: "+poc.ID, 2)
//				}
//				err = utils.WriteFile(pocPath, out)
//				if err != nil {
//					utils.CliError("Can't write file: "+pocPath, 3)
//				}
//				utils.SuccessF("%s: Add [%s] Tag", poc.ID, tag)
//			}
//		}
//	}
//
//}
//
//// 为poc移除标签
//func RemoveTags(tags []string, xrayPocMap map[string]xray_structs.Poc, nucleiPocMap map[string]nuclei_structs.Poc) {
//	for _, tag := range tags {
//		for pocPath, poc := range xrayPocMap {
//			pocTags := poc.Detail.Tags
//			if strings.Contains(pocTags, tag) {
//				poc.Detail.Tags = removeTag(pocTags, tag)
//				out, err := yaml.Marshal(poc)
//				if err != nil {
//					utils.CliError("Can't Marshal poc: "+poc.Name, 2)
//				}
//				err = utils.WriteFile(pocPath, out)
//				if err != nil {
//					utils.CliError("Can't write file: "+pocPath, 3)
//				}
//				utils.SuccessF("%s: Remove [%s] Tag", poc.Name, tag)
//			}
//		}
//		// nuclei tag 不区分大小写
//		for pocPath, poc := range nucleiPocMap {
//			pocTags := poc.Info.Tags.String()
//			tag = strings.ToLower(tag)
//			if strings.Contains(pocTags, tag) {
//				poc.Info.Tags = stringslice.StringSlice{
//					Value: removeTag(pocTags, tag),
//				}
//				out, err := yaml.Marshal(poc)
//				if err != nil {
//					utils.CliError("Can't Marshal poc: "+poc.ID, 2)
//				}
//				err = utils.WriteFile(pocPath, out)
//				if err != nil {
//					utils.CliError("Can't write file: "+pocPath, 3)
//				}
//				utils.SuccessF("%s: Remove [%s] Tag", poc.ID, tag)
//			}
//		}
//	}
//
//}
