package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"text/template"

	"github.com/davyxu/pbmeta"
)

const codeTsTemplate = `
import { pb } from "./pb";
export default class ProtoUtil {
	public static readonly CmdClassMap: { [k:number]: any } = {
				{{range .Protos}}
				{{range .Messages}}{{if .CheckGen}}
				{{.CMDID}}:pb.{{.Name}}, //pb.CMD.{{.MsgID}}
				{{end}}{{end}}{{end}}
	}
	public static readonly Cmd2Name: { [k:number]: string } = {
		{{range .Protos}}
		{{range .Messages}}{{if .CheckGen}}
		{{.CMDID}}:"{{.MsgID}}", //pb.CMD.{{.Name}}
		{{end}}{{end}}{{end}}
	}

    public static encode(cmd: number, msgObj: any): Uint8Array {
	   let clazz: any = this.CmdClassMap[cmd];
	   if (!clazz) {
		   return null;
	   }
	   let msg: typeof clazz = new clazz(msgObj);
	   let array: Uint8Array = clazz.encode(msg).finish();
	   return array;
    }

    public static decode(cmd: number, buffer: ArrayBuffer, offset: number): any {
	   let clazz: any = this.CmdClassMap[cmd];
	   if (!clazz) {
		   return null;
	   }

	   let array: Uint8Array = new Uint8Array(buffer, offset);
	   let msg: typeof clazz = clazz.decode(array);
	   return msg;
    }
}`

type msgModel struct {
	*pbmeta.Descriptor

	parent *pbmeta.FileDescriptor
}

func (self *msgModel) CMDID() int32 {
	return self.parent.EnumByName("CMD").ValueByName(self.LeadingComment()).Value()
}

func (self *msgModel) MsgID() string {
	return self.LeadingComment()
}

func (self *msgModel) FullName() string {
	return fmt.Sprintf("%s.%s", self.parent.PackageName(), self.Name())
}

func (self *msgModel) CheckGen() bool {
	return self.parent.EnumByName("CMD").ValueByName(self.LeadingComment()) != nil
}

type protoModel struct {
	*pbmeta.FileDescriptor

	Messages []*msgModel
}

func (self *protoModel) Name() string {
	return self.FileDescriptor.FileName()
}

type fileModel struct {
	TotalMessages int
	Protos        []*protoModel
	PackageName   string
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

// func printFile(pool *pbmeta.DescriptorPool) (string, bool) {

// 	tpl, err := template.New("msgid").Parse(codeTemplate)
// 	if err != nil {
// 		log.Errorln(err)
// 		return "", false
// 	}

// 	if pool.FileCount() == 0 {
// 		return "", false
// 	}

// 	var model fileModel
// 	model.PackageName = pool.File(0).PackageName()

// 	for f := 0; f < pool.FileCount(); f++ {

// 		file := pool.File(f)

// 		pm := &protoModel{
// 			FileDescriptor: file,
// 		}

// 		for m := 0; m < file.MessageCount(); m++ {

// 			d := file.Message(m)

// 			pm.Messages = append(pm.Messages, &msgModel{
// 				Descriptor: d,
// 				parent:     file,
// 			})

// 		}

// 		model.TotalMessages += file.MessageCount()

// 		model.Protos = append(model.Protos, pm)

// 	}

// 	var bf bytes.Buffer

// 	err = tpl.Execute(&bf, &model)
// 	if err != nil {
// 		log.Errorln(err)
// 		return "", false
// 	}

// 	err = formatCode(&bf)
// 	if err != nil {
// 		log.Errorln(err)
// 		return "", false
// 	}

// 	return bf.String(), true
// }

func formatCode(bf *bytes.Buffer) error {
	// Reformat generated code.
	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", bf, parser.ParseComments)
	if err != nil {
		return err
	}

	bf.Reset()

	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(bf, fset, ast)
	if err != nil {
		return err
	}

	return nil
}
