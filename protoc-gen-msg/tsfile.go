package main

import (
	"bytes"
	"text/template"

	"github.com/davyxu/pbmeta"
)

const codeTsTemplate = `
import { pb } from "../protobuf/pb";
export module pbCMD {
    export function getProtoClassMap() {
        return {
			{{range .Protos}}
			{{range .Messages}}{{if .CheckGen}}
			{{.CMDID}}:pb.{{.Name}},
			{{end}}{{end}}{{end}}
        }
    }
}
window["pbCMD"] = pbCMD;`

func (self *msgModel) CMDID() int32 {
	return self.parent.EnumByName("CMD").ValueByName(self.LeadingComment()).Value()
}

func printTsFile(pool *pbmeta.DescriptorPool) (string, bool) {

	tpl, err := template.New("msgid").Parse(codeTsTemplate)
	if err != nil {
		log.Errorln(err)
		return "", false
	}

	if pool.FileCount() == 0 {
		return "", false
	}

	var model fileModel
	model.PackageName = pool.File(0).PackageName()

	for f := 0; f < pool.FileCount(); f++ {

		file := pool.File(f)

		pm := &protoModel{
			FileDescriptor: file,
		}

		for m := 0; m < file.MessageCount(); m++ {

			d := file.Message(m)

			pm.Messages = append(pm.Messages, &msgModel{
				Descriptor: d,
				parent:     file,
			})

		}

		model.TotalMessages += file.MessageCount()

		model.Protos = append(model.Protos, pm)

	}

	var bf bytes.Buffer

	err = tpl.Execute(&bf, &model)
	if err != nil {
		log.Errorln(err)
		return "", false
	}

	return bf.String(), true
}
