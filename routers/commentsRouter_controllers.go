package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"] = append(beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"] = append(beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"],
        beego.ControllerComments{
            Method: "File",
            Router: `/file`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"] = append(beego.GlobalControllerRouter["github.com/sinksmell/files-cmp/controllers:CheckController"],
        beego.ControllerComments{
            Method: "Hash",
            Router: `/hash`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
