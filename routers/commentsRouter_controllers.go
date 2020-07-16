package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "Conn",
            Router: "/Conn",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "CreateConversation",
            Router: "/CreateConversation",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "DelConversation",
            Router: "/DelConversation",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "GetConversation",
            Router: "/GetConversation",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "GetMsgList",
            Router: "/GetMsgList",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chatservice/controllers:MsgController"] = append(beego.GlobalControllerRouter["chatservice/controllers:MsgController"],
        beego.ControllerComments{
            Method: "SendMsg",
            Router: "/SendMsg",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
